package virtual

import (
	"fmt"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
	"strings"
	"time"
)

type UsageCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewUsageCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *UsageCommand) {
	return &UsageCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *UsageCommand) Run(c *cli.Context) error {
	var periodic int
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	vsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	if !c.IsSet("summary-period") {
		periodic = 3600
	} else {
		periodic = c.Int("summary-period")
	}

	var start, end string
	var startDate, endDate time.Time

	if c.IsSet("start") {
		start = c.String("start")
		startDate, err = time.Parse(GetDateFormat(start), start)
		if err != nil {
			return errors.NewInvalidUsageError("Invalid start date: " + err.Error())
		}
	} else {
		startDate = time.Now()
	}
	if c.IsSet("end") {
		end = c.String("end")
		endDate, err = time.Parse(GetDateFormat(end), end)
		if err != nil {
			return errors.NewInvalidUsageError("Invalid end date: " + err.Error())
		}

	} else {
		endDate = startDate.AddDate(0, -1, 0)
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	vsUsage, err := cmd.VirtualServerManager.GetSummaryUsage(vsID, startDate, endDate, strings.ToUpper(c.String("valid-data")), periodic)
	if err != nil {
		return cli.NewExitError(T("Failed to upgrade virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": vsID})+err.Error(), 2)
	}
	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, vsUsage)
	}

	tableAverage := cmd.UI.Table([]string{T("Average")})
	tableUsage := cmd.UI.Table([]string{T("Counter"), T("Date"), T("Type")})
	count := 0
	counter := 0.0
	for _, data := range vsUsage {
		usageCounter := 0.0
		if strings.ToUpper(c.String("valid-data")) == "MEMORY-USAGE" {
			usageCounter = float64(*data.Counter) / 2.00 * 30.00
		} else {
			usageCounter = float64(*data.Counter)
		}
		tableUsage.Add(fmt.Sprintf("%.2f", usageCounter), utils.FormatSLTimePointer(data.DateTime), utils.FormatStringPointer(data.Type))
		count = count + 1
		counter = counter + usageCounter

	}
	average := counter / float64(count)
	tableAverage.Add(fmt.Sprintf("%.2f", average))

	tableAverage.Print()
	tableUsage.Print()
	return nil
}

func VSUsageMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "usage",
		Description: T("usage data over date range."),
		Usage: T(`${COMMAND_NAME} sl {{.Command}} usage IDENTIFIER [OPTIONS]
Usage information of a virtual server.
Example:
   ${COMMAND_NAME} sl {{.Command}} usage 1234 --start 2006-01-02 --end 2006-01-02 --valid-data cpu0`, map[string]interface{}{"Command": "vs"}),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "s,start",
				Usage:    T("Start Date e.g. 2019-3-4 (yyyy-MM-dd)  [required]"),
				Required: true,
			},
			cli.StringFlag{
				Name:     "e,end",
				Usage:    T("End Date e.g. 2019-4-2 (yyyy-MM-dd)  [required]"),
				Required: true,
			},
			cli.StringFlag{
				Name:     "t,valid-data",
				Usage:    T("Metric_Data_Type keyName e.g. CPU0, CPU1, MEMORY_USAGE, etc.  [required]"),
				Required: true,
			},
			cli.IntFlag{
				Name:  "p,summary-period",
				Usage: T("300, 600, 1800, 3600, 43200 or 86400 seconds."),
			},
			metadata.OutputFlag(),
		},
	}
}

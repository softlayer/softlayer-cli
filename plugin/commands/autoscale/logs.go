package autoscale

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type LogsCommand struct {
	UI               terminal.UI
	AutoScaleManager managers.AutoScaleManager
	SecurityManager  managers.SecurityManager
}

func NewLogsCommand(ui terminal.UI, autoScaleManager managers.AutoScaleManager, securityManager managers.SecurityManager) (cmd *LogsCommand) {
	return &LogsCommand{
		UI:               ui,
		AutoScaleManager: autoScaleManager,
		SecurityManager:  securityManager,
	}
}

func (cmd *LogsCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one identifier."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	autoScaleGroupId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return errors.NewInvalidUsageError(T("Autoscale group ID should be a number."))
	}

	datefilter := ""
	if c.IsSet("date-min") {
		date := c.String("date-min")
		_, err := time.Parse(time.RFC3339, date+"T00:00:00Z")
		if err != nil {
			return errors.NewInvalidUsageError(T("Invalid format date."))
		}
		splitDate := strings.Split(date, "-")
		datefilter = fmt.Sprintf("%s/%s/%s", splitDate[1], splitDate[2], splitDate[0])
	}

	mask := "mask[createDate,description]"
	autoScaleGroupLogs, err := cmd.AutoScaleManager.GetLogsScaleGroup(autoScaleGroupId, mask, datefilter)
	if err != nil {
		return cli.NewExitError(T("Failed to get AutoScale group logs.\n")+err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Date"), T("Entry")})
	for _, log := range autoScaleGroupLogs {
		table.Add(utils.FormatSLTimePointer(log.CreateDate), utils.FormatStringPointer(log.Description))
	}

	if outputFormat == "JSON" {
		table.PrintJson()
	} else {
		table.Print()
	}

	return nil
}

func AutoScaleLogsMetaData() cli.Command {
	return cli.Command{
		Category:    "autoscale",
		Name:        "logs",
		Description: T("Retreive logs for an Autoscale group."),
		Usage: T(`${COMMAND_NAME} sl autoscale logs IDENTIFIER [OPTIONS]

EXAMPLE: 
   ${COMMAND_NAME} sl autoscale logs 123456 [OPTIONS]`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "d,date-min",
				Usage: T("Earliest date to retrieve logs for [YYYY-MM-DD]."),
			},
			metadata.OutputFlag(),
		},
	}
}

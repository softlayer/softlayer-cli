package virtual

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type UsageCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Start                string
	End                  string
	ValidData            string
	SummaryPeriod        int
}

func NewUsageCommand(sl *metadata.SoftlayerCommand) (cmd *UsageCommand) {
	thisCmd := &UsageCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	subs := map[string]interface{}{"Command": "vs"}
	cobraCmd := &cobra.Command{
		Use:   "usage " + T("IDENTIFIER"),
		Short: T("usage data over date range."),
		Long: T(`${COMMAND_NAME} sl {{.Command}} usage IDENTIFIER [OPTIONS]
Usage information of a virtual server.
Example:
   ${COMMAND_NAME} sl {{.Command}} usage 1234 --start 2006-01-02 --end 2006-01-02 --valid-data cpu0`, subs),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().StringVarP(&thisCmd.Start, "start", "s", "", T("Start Date e.g. 2019-3-4 (yyyy-MM-dd)  [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.End, "end", "e", "", T("End Date e.g. 2019-4-2 (yyyy-MM-dd)  [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.ValidData, "valid-data", "t", "", T("Metric_Data_Type keyName e.g. CPU0, CPU1, MEMORY_USAGE, etc.  [required]"))
	cobraCmd.Flags().IntVarP(&thisCmd.SummaryPeriod, "summary-period", "p", 3600, T("300, 600, 1800, 3600, 43200 or 86400 seconds."))

	// the docs say these are required, but the code gives them a default value... so going to leave unrequired.
	// cobraCmd.MarkFlagRequired("start")
	// cobraCmd.MarkFlagRequired("end")
	cobraCmd.MarkFlagRequired("valid-data")
	return thisCmd
}

func (cmd *UsageCommand) Run(args []string) error {
	var periodic int

	vsID, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	periodic = cmd.SummaryPeriod

	var start, end string
	var startDate, endDate time.Time

	if cmd.Start != "" {
		start = cmd.Start
		startDate, err = time.Parse(GetDateFormat(start), start)
		if err != nil {
			return slErrors.NewInvalidUsageError("Invalid start date: " + err.Error())
		}
	} else {
		startDate = time.Now()
	}

	if cmd.End != "" {
		end = cmd.End
		endDate, err = time.Parse(GetDateFormat(end), end)
		if err != nil {
			return slErrors.NewInvalidUsageError("Invalid end date: " + err.Error())
		}
	} else {
		endDate = startDate.AddDate(0, -1, 0)
	}

	outputFormat := cmd.GetOutputFlag()

	vsUsage, err := cmd.VirtualServerManager.GetSummaryUsage(vsID, startDate, endDate, strings.ToUpper(cmd.ValidData), periodic)
	subs := map[string]interface{}{"VsID": vsID}
	if err != nil {
		return slErrors.NewAPIError(T("Failed to upgrade virtual server instance: {{.VsID}}.\n", subs), err.Error(), 2)
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
		if strings.ToUpper(cmd.ValidData) == "MEMORY-USAGE" {
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

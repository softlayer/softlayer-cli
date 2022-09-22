package cdn

import (
	"fmt"
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	*metadata.SoftlayerCommand
	CdnManager managers.CdnManager
	Command    *cobra.Command
	History    int
}

func NewDetailCommand(sl *metadata.SoftlayerCommand) *DetailCommand {
	thisCmd := &DetailCommand{
		SoftlayerCommand: sl,
		CdnManager:       managers.NewCdnManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "detail",
		Short: T("Detail a CDN Account."),
		Long:  T("${COMMAND_NAME} sl cdn detail"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.History, "history", 30, T("Bandwidth, Hits, Ratio counted over history number of days ago. 89 is the maximum."))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DetailCommand) Run(args []string) error {
	cdnId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("cdn ID")
	}

	history := cmd.History
	if history <= 0 || history > 89 {
		return slErr.NewInvalidUsageError("history")
	}

	outputFormat := cmd.GetOutputFlag()

	mask := ""
	cdnDetail, err := cmd.CdnManager.GetDetailCDN(cdnId, mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get CDN Detail. "), err.Error(), 2)
	}

	cdnMetrics, err := cmd.CdnManager.GetUsageMetrics(cdnId, history, mask)
	if err != nil {
		return cli.NewExitError(T("Failed to get CDN Metrics. ")+err.Error(), 2)
	}

	PrintDetailCDN(cdnDetail, cdnMetrics, cmd.UI, outputFormat)
	return nil
}

func PrintDetailCDN(cdnDetail datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, cdnMetrics datatypes.Container_Network_CdnMarketplace_Metrics, ui terminal.UI, outputFormat string) {

	totalBandwidth := fmt.Sprintf("%s GB", utils.FormatStringPointer(&cdnMetrics.Totals[0]))
	totalHits := utils.FormatStringPointer(&cdnMetrics.Totals[1])
	hitRatio := fmt.Sprintf("%s %%", utils.FormatStringPointer(&cdnMetrics.Totals[2]))

	table := ui.Table([]string{
		T("Name"),
		T("Value"),
	})
	table.Add(T("Unique id"), utils.FormatStringPointer(cdnDetail.UniqueId))
	table.Add(T("Hostname"), utils.FormatStringPointer(cdnDetail.Domain))
	table.Add(T("Protocol"), utils.FormatStringPointer(cdnDetail.Protocol))
	table.Add(T("Origin"), utils.FormatStringPointer(cdnDetail.OriginHost))
	table.Add(T("Origin type"), utils.FormatStringPointer(cdnDetail.OriginType))
	table.Add(T("Path"), utils.FormatStringPointer(cdnDetail.Path))
	table.Add(T("Provider"), utils.FormatStringPointer(cdnDetail.VendorName))
	table.Add(T("Status"), utils.FormatStringPointer(cdnDetail.Status))
	table.Add(T("Total bandwidth"), utils.FormatStringPointer(&totalBandwidth))
	table.Add(T("Total hits"), utils.FormatStringPointer(&totalHits))
	table.Add(T("Hit Radio"), utils.FormatStringPointer(&hitRatio))

	utils.PrintTable(ui, table, outputFormat)
}

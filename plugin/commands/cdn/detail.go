package cdn

import (
	"fmt"
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	UI         terminal.UI
	CdnManager managers.CdnManager
}

func NewDetailCommand(ui terminal.UI, cdnManager managers.CdnManager) (cmd *DetailCommand) {
	return &DetailCommand{
		UI:         ui,
		CdnManager: cdnManager,
	}
}

func DetailMetaData() cli.Command {
	return cli.Command{
		Category:    "cdn",
		Name:        "detail",
		Description: T("Detail a CDN Account."),
		Usage:       T(`${COMMAND_NAME} sl cdn detail`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "history",
				Usage: T("Bandwidth, Hits, Ratio counted over history number of days ago. 89 is the maximum."),
				Value: 30,
			},
			metadata.OutputFlag(),
		},
	}
}

func (cmd *DetailCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return slErr.NewInvalidUsageError(T("This command requires one argument."))
	}

	cdnId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("cdn ID")
	}

	history := c.Int("history")
	if history <= 0 || history > 89 {
		return slErr.NewInvalidUsageError("history")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	mask := ""
	cdnDetail, err := cmd.CdnManager.GetDetailCDN(cdnId, mask)
	if err != nil {
		return cli.NewExitError(T("Failed to get CDN Detail. ")+err.Error(), 2)
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

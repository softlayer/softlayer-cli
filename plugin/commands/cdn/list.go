package cdn

import (
	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	UI         terminal.UI
	CdnManager managers.CdnManager
}

func NewListCommand(ui terminal.UI, cdnManager managers.CdnManager) (cmd *ListCommand) {
	return &ListCommand{
		UI:         ui,
		CdnManager: cdnManager,
	}
}

func ListMetaData() cli.Command {
	return cli.Command{
		Category:    "cdn",
		Name:        "list",
		Description: T("List all CDN accounts."),
		Usage:       T(`${COMMAND_NAME} sl cdn list`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

func (cmd *ListCommand) Run(c *cli.Context) error {
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	mask := ""
	cdnList, err := cmd.CdnManager.GetNetworkCdnMarketplaceConfigurationMapping(mask)
	if err != nil {
		return cli.NewExitError(T("Failed to get CDN List. ")+err.Error(), 2)
	}

	PrintNetworkCdnMarketplaceConfigurationMapping(cdnList, cmd.UI, outputFormat)
	return nil
}

func PrintNetworkCdnMarketplaceConfigurationMapping(cdnList []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, ui terminal.UI, outputFormat string) {
	table := ui.Table([]string{
		T("Unique Id"),
		T("Domain"),
		T("Origin"),
		T("Vendor"),
		T("Cname"),
		T("Status"),
	})
	for _, cdn := range cdnList {
		table.Add(
			utils.FormatStringPointer(cdn.UniqueId),
			utils.FormatStringPointer(cdn.Domain),
			utils.FormatStringPointer(cdn.OriginHost),
			utils.FormatStringPointer(cdn.VendorName),
			utils.FormatStringPointer(cdn.Cname),
			utils.FormatStringPointer(cdn.Status),
		)
	}
	utils.PrintTable(ui, table, outputFormat)
}

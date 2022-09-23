package cdn

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	*metadata.SoftlayerCommand
	CdnManager managers.CdnManager
	Command    *cobra.Command
}

func NewListCommand(sl *metadata.SoftlayerCommand) *ListCommand {
	thisCmd := &ListCommand{
		SoftlayerCommand: sl,
		CdnManager:       managers.NewCdnManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List all CDN accounts."),
		Long:  T("${COMMAND_NAME} sl cdn list"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ListCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	cdnList, err := cmd.CdnManager.GetNetworkCdnMarketplaceConfigurationMapping()
	if err != nil {
		return errors.NewAPIError(T("Failed to get CDN List. "), err.Error(), 2)
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

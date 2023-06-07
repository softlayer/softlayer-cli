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

type OriginListCommand struct {
	*metadata.SoftlayerCommand
	CdnManager managers.CdnManager
	Command    *cobra.Command
}

func NewOriginListCommand(sl *metadata.SoftlayerCommand) *OriginListCommand {
	thisCmd := &OriginListCommand{
		SoftlayerCommand: sl,
		CdnManager:       managers.NewCdnManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "origin-list " + T("IDENTIFIER"),
		Short: T("List origin path for an existing CDN mapping."),
		Long:  T("${COMMAND_NAME} sl cdn origin-list"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *OriginListCommand) Run(args []string) error {
	cdnId := args[0]

	outputFormat := cmd.GetOutputFlag()

	cdnOriginList, err := cmd.CdnManager.GetOrigins(cdnId)
	if err != nil {
		return errors.NewAPIError(T("Failed to get origins list for CDN: {{.cdnId}}.", map[string]interface{}{"cdnId": cdnId}), err.Error(), 2)
	}

	PrintOriginsList(cdnOriginList, cmd.UI, outputFormat)
	return nil
}

func PrintOriginsList(cdnOriginList []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path, ui terminal.UI, outputFormat string) {
	table := ui.Table([]string{
		T("Path"),
		T("Origin"),
		T("Http Port"),
		T("Https Port"),
		T("Status"),
	})
	for _, origin := range cdnOriginList {
		table.Add(
			utils.FormatStringPointer(origin.Path),
			utils.FormatStringPointer(origin.Origin),
			utils.FormatIntPointer(origin.HttpPort),
			utils.FormatIntPointer(origin.HttpsPort),
			utils.FormatStringPointer(origin.Status),
		)
	}
	utils.PrintTable(ui, table, outputFormat)
}

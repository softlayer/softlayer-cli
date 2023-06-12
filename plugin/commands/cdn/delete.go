package cdn

import (
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type DeleteCommand struct {
	*metadata.SoftlayerCommand
	CdnManager managers.CdnManager
	Command    *cobra.Command
}

func NewDeleteCommand(sl *metadata.SoftlayerCommand) *DeleteCommand {
	thisCmd := &DeleteCommand{
		SoftlayerCommand: sl,
		CdnManager:       managers.NewCdnManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "delete " + T("IDENTIFIER"),
		Short: T("Delete a CDN domain mapping."),
		Long:  T(`${COMMAND_NAME} sl cdn delete`),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DeleteCommand) Run(args []string) error {
	cdnID := args[0]

	_, err := cmd.CdnManager.DeleteCDN(cdnID)
	if err != nil {
		return errors.NewAPIError(T("Failed to deleted a CDN."), err.Error(), 2)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Cdn with uniqueId: {{.ID}} was deleted.", map[string]interface{}{"ID": cdnID}))
	return nil
}

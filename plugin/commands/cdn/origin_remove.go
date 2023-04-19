package cdn

import (
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type OriginRemoveCommand struct {
	*metadata.SoftlayerCommand
	CdnManager managers.CdnManager
	Command    *cobra.Command
}

func NewOriginRemoveCommand(sl *metadata.SoftlayerCommand) *OriginRemoveCommand {
	thisCmd := &OriginRemoveCommand{
		SoftlayerCommand: sl,
		CdnManager:       managers.NewCdnManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "origin-remove",
		Short: T("Removes an origin path for an existing CDN mapping."),
		Long:  T("${COMMAND_NAME} sl cdn origin-remove"),
		Args:  metadata.TwoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *OriginRemoveCommand) Run(args []string) error {
	cdnId := args[0]
	path := args[1]

	_, err := cmd.CdnManager.RemoveOrigin(cdnId, path)
	if err != nil {
		return errors.NewAPIError(T("Failed to delete origin."), err.Error(), 2)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("The origin {{.Id}} with path {{.Path}} was deleted.", map[string]interface{}{"Path": path, "Id": cdnId}))
	return nil
}

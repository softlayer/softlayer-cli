package tags

import (
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DeleteCommand struct {
	*metadata.SoftlayerCommand
	TagsManager managers.TagsManager
	Command     *cobra.Command
}

func NewDeleteCommand(sl *metadata.SoftlayerCommand) (cmd *DeleteCommand) {
	thisCmd := &DeleteCommand{
		SoftlayerCommand: sl,
		TagsManager:      managers.NewTagsManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "delete " + T("[TAG NAME]"),
		Short: T("Removes an empty tag from your account."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DeleteCommand) Run(args []string) error {

	tagName := args[0]
	outputFormat := cmd.GetOutputFlag()
	success, err := cmd.TagsManager.DeleteTag(tagName)

	if err != nil {
		return errors.NewAPIError(T("Failed to delete Tag"), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, success)
	}

	cmd.UI.Print("%v", success)
	return nil

}

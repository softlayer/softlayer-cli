package tags

import (
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CleanupCommand struct {
	*metadata.SoftlayerCommand
	TagsManager managers.TagsManager
	Command     *cobra.Command
	DryRun      bool
}

func NewCleanupCommand(sl *metadata.SoftlayerCommand) (cmd *CleanupCommand) {
	thisCmd := &CleanupCommand{
		SoftlayerCommand: sl,
		TagsManager:      managers.NewTagsManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "cleanup",
		Short: T("Removes all empty tags."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVar(&thisCmd.DryRun, "dry-run", false, T("Don't delete, just show what will be deleted."))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CleanupCommand) Run(args []string) error {

	unattachedTags, err := cmd.TagsManager.GetUnattachedTags("")
	if err != nil {
		return errors.NewAPIError(T("Failed to get Unattached Tags."), err.Error(), 2)
	}
	for _, tag := range unattachedTags {
		tag_replace := map[string]interface{}{"tag": *tag.Name}
		if cmd.DryRun {
			cmd.UI.Print(T("(Dry Run) Removing Tag: {{.tag}}.", tag_replace))
		} else {
			success, err := cmd.TagsManager.DeleteTag(*tag.Name)
			if err != nil {
				cmd.UI.Print(T("Failed to delete Tag: {{.tag}}.", tag_replace) + "\n" + err.Error() + "\n")
			}
			if success {
				cmd.UI.Print(T("Removing Tag: {{.tag}}.", tag_replace))
			}
		}
	}
	return nil

}

package tags

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

const DRY_RUN_FLAG = "dry-run"

type CleanupCommand struct {
	UI          terminal.UI
	TagsManager managers.TagsManager
}

func NewCleanupCommand(ui terminal.UI, tagsManager managers.TagsManager) (cmd *CleanupCommand) {
	return &CleanupCommand{
		UI:          ui,
		TagsManager: tagsManager,
	}
}

func (cmd *CleanupCommand) Run(c *cli.Context) error {

	unattachedTags, err := cmd.TagsManager.GetUnattachedTags("")
	if err != nil {
		return cli.NewExitError(T("Failed to get Unattached Tags.")+"\n"+err.Error(), 2)
	}
	for _, tag := range unattachedTags {
		tag_replace := map[string]interface{}{"tag": *tag.Name}
		if c.IsSet(DRY_RUN_FLAG) && c.Bool(DRY_RUN_FLAG) {
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

func TagsCleanupMetaData() cli.Command {
	return cli.Command{
		Category:    "tags",
		Name:        "cleanup",
		Description: T("Removes all empty tags."),
		Usage: T(`${COMMAND_NAME} sl tags cleanup [OPTIONS]

EXAMPLE:
	${COMMAND_NAME} sl tags cleanup`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "dry-run",
				Usage: T("Don't delete, just show what will be deleted."),
			},
		},
	}
}

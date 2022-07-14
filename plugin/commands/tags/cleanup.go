package tags

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

const dryRunFlag = "dry-run"

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
		return cli.NewExitError(T("Failed to get Unattached Tags.\n")+err.Error(), 2)
	}
	for _, tag := range unattachedTags {
		if c.IsSet(dryRunFlag) && c.Bool(dryRunFlag) {
			cmd.UI.Print(T("(Dry Run) Removing Tag: {{.tag}}.", map[string]interface{}{"tag": *tag.Name}))
		} else {
			success, err := cmd.TagsManager.DeleteTag(*tag.Name)
			if err != nil {
				cmd.UI.Print(T("Failed to delete Tag: {{.tag}}.\n", map[string]interface{}{"tag": *tag.Name}) + err.Error() + "\n")
			}
			if success {
				cmd.UI.Print(T("Removing Tag: {{.tag}}.", map[string]interface{}{"tag": *tag.Name}))
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
			metadata.OutputFlag(),
		},
	}
}

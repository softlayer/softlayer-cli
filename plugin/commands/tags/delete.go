package tags

import (
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)


type DeleteCommand struct {
	UI 			terminal.UI
	TagsManager managers.TagsManager
}


func NewDeleteCommand(ui terminal.UI, tagsManager managers.TagsManager) (cmd *DeleteCommand) {
	return &DeleteCommand{
		UI:             ui,
		TagsManager: 	tagsManager,
	}	
}

func (cmd *DeleteCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	tagName := c.Args()[0]
	outputFormat, _ := metadata.CheckOutputFormat(c, cmd.UI)
	success, err := cmd.TagsManager.DeleteTag(tagName)

	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, success)
	}

	cmd.UI.Print("%v", success)
	return nil

}

func TagsDeleteMetaData() cli.Command {
	return cli.Command{
		Category:    "tags",
		Name:        "delete",
		Description: T("Removes an empty tag from your account."),
		Usage: T(`${COMMAND_NAME} sl tags delete [TAG NAME]

EXAMPLE:
	${COMMAND_NAME} sl tags delete tag1
	Removes "tag" from your account.
`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

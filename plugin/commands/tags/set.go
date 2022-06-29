package tags

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type SetCommand struct {
	UI          terminal.UI
	TagsManager managers.TagsManager
}

func NewSetCommand(ui terminal.UI, tagsManager managers.TagsManager) (cmd *SetCommand) {
	return &SetCommand{
		UI:          ui,
		TagsManager: tagsManager,
	}
}

func (cmd *SetCommand) Run(c *cli.Context) error {
	tags := c.String("tags")
	keyName := c.String("key-name")
	resourceId := c.Int("resource-id")

	response, err := cmd.TagsManager.SetTags(tags, keyName, resourceId)
	if err != nil {
		return cli.NewExitError(T("Failed to set tags.\n")+err.Error(), 2)
	}
	if response {
		cmd.UI.Ok()
		cmd.UI.Print("Set tags successfully")
	}

	return nil
}

func TagsSetMetaData() cli.Command {
	return cli.Command{
		Category:    "tags",
		Name:        "set",
		Description: T("Set Tags."),
		Usage: T(`${COMMAND_NAME} sl tags set [OPTIONS]

EXAMPLE:
	${COMMAND_NAME} sl tags set --tags 'tag1,tag2' --key-name HARDWARE --resource-id 123456
`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Required: true,
				Name:     "tags",
				Usage:    T("Comma seperated list of tags, enclosed in quotes. 'tag1,tag2'  [required]"),
			},
			cli.StringFlag{
				Required: true,
				Name:     "key-name",
				Usage:    T("Key name of a tag type e.g. GUEST, HARDWARE. See slcli tags taggable output.  [required]"),
			},
			cli.IntFlag{
				Required: true,
				Name:     "resource-id",
				Usage:    T("ID of the object being tagged  [required]"),
			},
		},
	}
}

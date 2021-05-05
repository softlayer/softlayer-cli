package tags

import (
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)


func GetCommandAcionBindings( ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	tagsManager := managers.NewTagsManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		//tags
		NS_TAGS_NAME + "-" + CMD_TAGS_LIST_NAME: func(c *cli.Context) error {
			return NewListCommand(ui, tagsManager).Run(c)
		},
		NS_TAGS_NAME + "-" + CMD_TAGS_DETAIL_NAME: func(c *cli.Context) error {
			return NewDetailCommand(ui, tagsManager).Run(c)
		},
		NS_TAGS_NAME + "-" + CMD_TAGS_DELETE_NAME: func(c *cli.Context) error {
			return NewDeleteCommand(ui, tagsManager).Run(c)
		},
	}

	return CommandActionBindings
}
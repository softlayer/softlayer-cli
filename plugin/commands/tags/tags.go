package tags

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	tagsManager := managers.NewTagsManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"tags-delete": func(c *cli.Context) error {
			return NewDeleteCommand(ui, tagsManager).Run(c)
		},
		"tags-detail": func(c *cli.Context) error {
			return NewDetailCommand(ui, tagsManager).Run(c)
		},
		"tags-list": func(c *cli.Context) error {
			return NewListCommand(ui, tagsManager).Run(c)
		},
		"tags-set": func(c *cli.Context) error {
			return NewSetCommand(ui, tagsManager).Run(c)
		},
		"tags-cleanup": func(c *cli.Context) error {
			return NewCleanupCommand(ui, tagsManager).Run(c)
		},
	}
	return CommandActionBindings
}

func TagsNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "tags",
		Description: T("Classic infrastructure Tag management"),
	}
}

func TagsMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "tags",
		Description: T("Classic infrastructure Tag management"),
		Usage:       "${COMMAND_NAME} sl tags",
		Subcommands: []cli.Command{
			TagsListMetaData(),
			TagsDetailsMetaData(),
			TagsDeleteMetaData(),
			TagsSetMetaData(),
			TagsCleanupMetaData(),
		},
	}
}

package email

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	emailManager := managers.NewEmailManager(session)
	CommandActionBindings := map[string]func(c *cli.Context) error{
		"email-list": func(c *cli.Context) error {
			return NewListCommand(ui, emailManager).Run(c)
		},
	}

	return CommandActionBindings
}

func EmailNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "email",
		Description: T("Classic infrastructure Email commands"),
	}
}

func EmailMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "email",
		Description: T("Classic infrastructure Email commands"),
		Usage:       "${COMMAND_NAME} sl email",
		Subcommands: []cli.Command{
			ListMetaData(),
		},
	}
}

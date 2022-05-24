package licenses

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	licensesManager := managers.NewLicensesManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"licenses-create-options": func(c *cli.Context) error {
			return NewLicensesOptionsCommand(ui, licensesManager).Run(c)
		},
		"licenses-create": func(c *cli.Context) error {
			return NewCreateCommand(ui, licensesManager).Run(c)
		},
		"licenses-cancel": func(c *cli.Context) error {
			return NewCancelItemCommand(ui, licensesManager).Run(c)
		},
	}

	return CommandActionBindings
}

func LicensesNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "licenses",
		Description: T("Classic infrastructure Licenses"),
	}
}

func LicensesMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "licenses",
		Description: T("Classic infrastructure Licenses"),
		Usage:       "${COMMAND_NAME} sl licenses",
		Subcommands: []cli.Command{
			LicensesCreateOptionsMetaData(),
			CreateMetaData(),
			CancelItemMetaData(),
		},
	}
}

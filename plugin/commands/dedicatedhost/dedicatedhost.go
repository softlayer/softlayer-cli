package dedicatedhost

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	dedicatedhostManager := managers.NewDedicatedhostManager(session)
	networkManager := managers.NewNetworkManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"dedicatedhost-list-guests": func(c *cli.Context) error {
			return NewListGuestsCommand(ui, dedicatedhostManager).Run(c)
		},
		"dedicatedhost-create": func(c *cli.Context) error {
			return NewCreateCommand(ui, dedicatedhostManager, networkManager, context).Run(c)
		},
		"dedicatedhost-detail": func(c *cli.Context) error {
			return NewDetailCommand(ui, dedicatedhostManager).Run(c)
		},
		"dedicatedhost-cancel-guests": func(c *cli.Context) error {
			return NewCancelCommand(ui, dedicatedhostManager).Run(c)
		},
	}

	return CommandActionBindings
}

func DedicatedhostNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "dedicatedhost",
		Description: T("Classic infrastructure Dedicatedhost"),
	}
}

func DedicatedhostMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "dedicatedhost",
		Description: T("Classic infrastructure Dedicatedhost"),
		Usage:       "${COMMAND_NAME} sl dedicatedhost",
		Subcommands: []cli.Command{
			// DedicatedhostListMetaData(),
			DedicatedhostListGuestsMetaData(),
			DedicatedhostCreateMetaData(),
			DedicatedhostDetailMetaData(),
			DedicatedhostCancelGuestsMetaData(),
		},
	}
}

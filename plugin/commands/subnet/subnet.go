package subnet

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	networkManager := managers.NewNetworkManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"subnet-cancel": func(c *cli.Context) error {
			return NewCancelCommand(ui, networkManager).Run(c)
		},
		"subnet-create": func(c *cli.Context) error {
			return NewCreateCommand(ui, networkManager).Run(c)
		},
		"subnet-detail": func(c *cli.Context) error {
			return NewDetailCommand(ui, networkManager).Run(c)
		},
		"subnet-list": func(c *cli.Context) error {
			return NewListCommand(ui, networkManager).Run(c)
		},
		"subnet-lookup": func(c *cli.Context) error {
			return NewLookupCommand(ui, networkManager).Run(c)
		},
		"subnet-route": func(c *cli.Context) error {
			return NewRouteCommand(ui, networkManager).Run(c)
		},
		"subnet-clear-route": func(c *cli.Context) error {
			return NewClearRouteCommand(ui, networkManager).Run(c)
		},
		"subnet-edit": func(c *cli.Context) error {
			return NewEditCommand(ui, networkManager).Run(c)
		},
		"subnet-edit-ip": func(c *cli.Context) error {
			return NewEditIpCommand(ui, networkManager).Run(c)
		},
	}

	return CommandActionBindings
}

func SubnetNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "subnet",
		Description: T("Classic infrastructure Network subnets"),
	}
}

func SubnetMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "subnet",
		Description: T("Classic infrastructure Network subnets"),
		Usage:       "${COMMAND_NAME} sl subnet",
		Subcommands: []cli.Command{
			SubnetCancelMetaData(),
			SubnetCreateMetaData(),
			SubnetDetailMetaData(),
			SubnetListMetaData(),
			SubnetLookupMetaData(),
			SubnetRouteMetaData(),
			SubnetClearRouteMetaData(),
			SubnetEditMetaData(),
			SubnetEditIpMetaData(),
		},
	}
}

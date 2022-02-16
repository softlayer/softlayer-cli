package globalip

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
		"globalip-assign": func(c *cli.Context) error {
			return NewAssignCommand(ui, networkManager).Run(c)
		},
		"globalip-cancel": func(c *cli.Context) error {
			return NewCancelCommand(ui, networkManager).Run(c)
		},
		"globalip-create": func(c *cli.Context) error {
			return NewCreateCommand(ui, networkManager).Run(c)
		},
		"globalip-list": func(c *cli.Context) error {
			return NewListCommand(ui, networkManager).Run(c)
		},
		"globalip-unassign": func(c *cli.Context) error {
			return NewUnassignCommand(ui, networkManager).Run(c)
		},
	}

	return CommandActionBindings
}

func GlobalIpNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "globalip",
		Description: T("Classic infrastructure Global IP addresses"),
	}
}

func GlobalIpMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "globalip",
		Description: T("Classic infrastructure Global IP addresses"),
		Usage:       "${COMMAND_NAME} sl globalip",
		Subcommands: []cli.Command{
			GlobalIpCreateMetaData(),
			GlobalIpAssignMetaData(),
			GlobalIpCancelMetaData(),
			GlobalIpListMetaData(),
			GlobalIpUnassignMetaData(),
		},
	}
}

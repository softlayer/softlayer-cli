package vlan

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
	context = plugin.InitPluginContext("softlayer")

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"vlan-cancel": func(c *cli.Context) error {
			return NewCancelCommand(ui, networkManager).Run(c)
		},
		"vlan-create": func(c *cli.Context) error {
			return NewCreateCommand(ui, networkManager, context).Run(c)
		},
		"vlan-detail": func(c *cli.Context) error {
			return NewDetailCommand(ui, networkManager).Run(c)
		},
		"vlan-edit": func(c *cli.Context) error {
			return NewEditCommand(ui, networkManager).Run(c)
		},
		"vlan-list": func(c *cli.Context) error {
			return NewListCommand(ui, networkManager).Run(c)
		},
		"vlan-options": func(c *cli.Context) error {
			return NewOptionsCommand(ui, networkManager).Run(c)
		},
	}

	return CommandActionBindings
}

func VlanNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "vlan",
		Description: T("Classic infrastructure Network VLANs"),
	}
}

func VlanMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "vlan",
		Description: T("Classic infrastructure Network VLANs"),
		Usage:       "${COMMAND_NAME} sl vlan",
		Subcommands: []cli.Command{
			VlanCreateMetaData(),
			VlanCancelMetaData(),
			VlanDetailMetaData(),
			VlanEditMetaData(),
			VlanListMetaData(),
			VlanOptionsMetaData(),
		},
	}
}

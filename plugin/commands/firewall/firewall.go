package firewall

// These commands are not actually part of the ibmcloud sl command at the moment because they are legacy
// Need to update these at some point.

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func FirewallNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "file",
		Description: T("Classic infrastructure Firewalls"),
	}
}

func FirewallMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "firewall",
		Description: T("Classic infrastructure Firewalls"),
		Usage:       "${COMMAND_NAME} sl firewall",
		Subcommands: []cli.Command{
			FirewallAddMetaData(),
			FirewallCancelMetaData(),
			FirewallDetailMetaData(),
			FirewallEditMetaData(),
			FirewallListMetaData(),
		},
	}
}


func GetCommandAcionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	firewallManager := managers.NewFirewallManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"firewall-add": func(c *cli.Context) error {
			return NewAddCommand(ui, firewallManager).Run(c)
		},
		"firewall-cancel": func(c *cli.Context) error {
			return NewCancelCommand(ui, firewallManager).Run(c)
		},
		"firewall-detail": func(c *cli.Context) error {
			return NewDetailCommand(ui, firewallManager).Run(c)
		},
		"firewall-edit": func(c *cli.Context) error {
			return NewEditCommand(ui, firewallManager).Run(c)
		},
		"firewall-list": func(c *cli.Context) error {
			return NewListCommand(ui, firewallManager).Run(c)
		},
	}

	return CommandActionBindings
}
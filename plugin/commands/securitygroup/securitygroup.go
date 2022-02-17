package securitygroup

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
	vsManager := managers.NewVirtualServerManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"securitygroup-create": func(c *cli.Context) error {
			return NewCreateCommand(ui, networkManager).Run(c)
		},
		"securitygroup-delete": func(c *cli.Context) error {
			return NewDeleteCommand(ui, networkManager).Run(c)
		},
		"securitygroup-detail": func(c *cli.Context) error {
			return NewDetailCommand(ui, networkManager).Run(c)
		},
		"securitygroup-edit": func(c *cli.Context) error {
			return NewEditCommand(ui, networkManager).Run(c)
		},
		"securitygroup-interface-add": func(c *cli.Context) error {
			return NewInterfaceAddCommand(ui, networkManager, vsManager).Run(c)
		},
		"securitygroup-interface-list": func(c *cli.Context) error {
			return NewInterfaceListCommand(ui, networkManager).Run(c)
		},
		"securitygroup-interface-remove": func(c *cli.Context) error {
			return NewInterfaceRemoveCommand(ui, networkManager, vsManager).Run(c)
		},
		"securitygroup-list": func(c *cli.Context) error {
			return NewListCommand(ui, networkManager).Run(c)
		},
		"securitygroup-rule-add": func(c *cli.Context) error {
			return NewRuleAddCommand(ui, networkManager).Run(c)
		},
		"securitygroup-rule-edit": func(c *cli.Context) error {
			return NewRuleEditCommand(ui, networkManager).Run(c)
		},
		"securitygroup-rule-list": func(c *cli.Context) error {
			return NewRuleListCommand(ui, networkManager).Run(c)
		},
		"securitygroup-rule-remove": func(c *cli.Context) error {
			return NewRuleRemoveCommand(ui, networkManager).Run(c)
		},
	}
	return CommandActionBindings
}

func SecurityGroupNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "securitygroup",
		Description: T("Classic infrastructure network security groups"),
	}
}

func SecurityGroupMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "securitygroup",
		Description: T("Classic infrastructure network security groups"),
		Usage:       "${COMMAND_NAME} sl securitygroup",
		Subcommands: []cli.Command{
			SecurityGroupCreateMetaData(),
			SecurityGroupDeleteMetaData(),
			SecurityGroupDetailMetaData(),
			SecurityGroupEditMetaData(),
			SecurityGroupInterfaceAddMetaData(),
			SecurityGroupInterfaceListMetaData(),
			SecurityGroupInterfaceRemoveMetaData(),
			SecurityGroupListMetaData(),
			SecurityGroupRuleAddMetaData(),
			SecurityGroupRuleEditMetaData(),
			SecurityGroupRuleListMetaData(),
			SecurityGroupRuleRemoveMetaData(),
		},
	}
}

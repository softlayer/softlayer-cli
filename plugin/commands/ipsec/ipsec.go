package ipsec

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	ipsecManager := managers.NewIPSECManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"ipsec-config": func(c *cli.Context) error {
			return NewConfigCommand(ui, ipsecManager).Run(c)
		},
		"ipsec-cancel": func(c *cli.Context) error {
			return NewCancelCommand(ui, ipsecManager).Run(c)
		},
		"ipsec-order": func(c *cli.Context) error {
			return NewOrderCommand(ui, ipsecManager, context).Run(c)
		},
		"ipsec-detail": func(c *cli.Context) error {
			return NewDetailCommand(ui, ipsecManager).Run(c)
		},
		"ipsec-list": func(c *cli.Context) error {
			return NewListCommand(ui, ipsecManager).Run(c)
		},
		"ipsec-subnet-add": func(c *cli.Context) error {
			return NewAddSubnetCommand(ui, ipsecManager).Run(c)
		},
		"ipsec-subnet-remove": func(c *cli.Context) error {
			return NewRemoveSubnetCommand(ui, ipsecManager).Run(c)
		},
		"ipsec-translation-add": func(c *cli.Context) error {
			return NewAddTranslationCommand(ui, ipsecManager).Run(c)
		},
		"ipsec-translation-remove": func(c *cli.Context) error {
			return NewRemoveTranslationCommand(ui, ipsecManager).Run(c)
		},
		"ipsec-translation-update": func(c *cli.Context) error {
			return NewUpdateTranslationCommand(ui, ipsecManager).Run(c)
		},
		"ipsec-update": func(c *cli.Context) error {
			return NewUpdateCommand(ui, ipsecManager).Run(c)
		},
	}

	return CommandActionBindings
}

func IpsecNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "ipsec",
		Description: T("Classic infrastructure IPSEC VPN"),
	}
}

func IpsecMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "ipsec",
		Description: T("Classic infrastructure IPSEC VPN"),
		Usage:       "${COMMAND_NAME} sl ipsec",
		Subcommands: []cli.Command{
			IpsecCancelMetaData(),
			IpsecConfigMetaData(),
			IpsecOrderMetaData(),
			IpsecDetailMetaData(),
			IpsecListMetaData(),
			IpsecSubnetAddMetaData(),
			IpsecSubnetRemoveMetaData(),
			IpsecTransAddMetaData(),
			IpsecTransRemoveMetaData(),
			IpsecTransUpdataMetaData(),
			IpsecUpdateMetaData(),
		},
	}
}

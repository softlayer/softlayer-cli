package dns

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	dnsManager := managers.NewDNSManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"dns-import": func(c *cli.Context) error {
			return NewImportCommand(ui, dnsManager).Run(c)
		},
		"dns-record-add": func(c *cli.Context) error {
			return NewRecordAddCommand(ui, dnsManager).Run(c)
		},
		"dns-record-edit": func(c *cli.Context) error {
			return NewRecordEditCommand(ui, dnsManager).Run(c)
		},
		"dns-record-list": func(c *cli.Context) error {
			return NewRecordListCommand(ui, dnsManager).Run(c)
		},
		"dns-record-remove": func(c *cli.Context) error {
			return NewRecordRemoveCommand(ui, dnsManager).Run(c)
		},
		"dns-zone-create": func(c *cli.Context) error {
			return NewZoneCreateCommand(ui, dnsManager).Run(c)
		},
		"dns-zone-delete": func(c *cli.Context) error {
			return NewZoneDeleteCommand(ui, dnsManager).Run(c)
		},
		"dns-zone-list": func(c *cli.Context) error {
			return NewZoneListCommand(ui, dnsManager).Run(c)
		},
		"dns-zone-print": func(c *cli.Context) error {
			return NewZonePrintCommand(ui, dnsManager).Run(c)
		},
	}

	return CommandActionBindings
}

func DnsNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "dns",
		Description: T("Classic infrastructure Domain Name System"),
	}
}

func DnsMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "dns",
		Description: T("Classic infrastructure Domain Name System"),
		Usage:       "${COMMAND_NAME} sl dns",
		Subcommands: []cli.Command{
			DnsImportMetaData(),
			DnsRecordAddMetaData(),
			DnsRecordEditMetaData(),
			DnsRecordListMetaData(),
			DnsRecordRemoveMetaData(),
			DnsZoneCreateMetaData(),
			DnsZoneDeleteMetaData(),
			DnsZoneListMetaData(),
			DnsZonePrintMetaData(),
		},
	}
}

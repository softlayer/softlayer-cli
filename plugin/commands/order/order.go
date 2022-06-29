package order

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	orderManager := managers.NewOrderManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"order-category-list": func(c *cli.Context) error {
			return NewCategoryListCommand(ui, orderManager).Run(c)
		},
		"order-item-list": func(c *cli.Context) error {
			return NewItemListCommand(ui, orderManager).Run(c)
		},
		"order-package-list": func(c *cli.Context) error {
			return NewPackageListCommand(ui, orderManager).Run(c)
		},
		"order-package-locations": func(c *cli.Context) error {
			return NewPackageLocationCommand(ui, orderManager).Run(c)
		},
		"order-place": func(c *cli.Context) error {
			return NewPlaceCommand(ui, orderManager, context).Run(c)
		},
		"order-place-quote": func(c *cli.Context) error {
			return NewPlaceQuoteCommand(ui, orderManager, context).Run(c)
		},
		"order-preset-list": func(c *cli.Context) error {
			return NewPresetListCommand(ui, orderManager).Run(c)
		},
		"order-quote-list": func(c *cli.Context) error {
			return NewQuoteListCommand(ui, orderManager).Run(c)
		},
	}

	return CommandActionBindings
}

func OrderNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "order",
		Description: T("Classic infrastructure Orders"),
	}
}

func OrderMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "order",
		Usage:       "${COMMAND_NAME} sl order",
		Description: T("Classic infrastructure Orders"),
		Subcommands: []cli.Command{
			OrderCategoryListMetaData(),
			OrderItemListMetaData(),
			OrderPackageListMetaData(),
			OrderPackageLocaionMetaData(),
			OrderPlaceMetaData(),
			OrderPlaceQuoteMetaData(),
			OrderPresetListMetaData(),
			OrderQuoteListMetaData(),
		},
	}
}

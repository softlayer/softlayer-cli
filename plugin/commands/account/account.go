package account

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

/*
This account package follows a slightly different pattern than the other CLI commands
because I'd like to eventually adpot this pattern throughout to get away from having metadata files
for every command.
*/

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use: "account",
		Short: T("Classic infrastructure Account commands"),
		Long: "${COMMAND_NAME} sl account",
		RunE: nil,
	}
	// cobraCmd.AddCommand(account.New<COMMAND>Command(ui, session))
	cobraCmd.AddCommand(NewBandwidthPoolsCommand(sl).Command)
	cobraCmd.AddCommand(NewBandwidthPoolsDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewBillingItemsCommand(sl).Command)
	return cobraCmd	
}

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	accountManager := managers.NewAccountManager(session)
	CommandActionBindings := map[string]func(c *cli.Context) error{
		// "account-bandwidth-pools": func(c *cli.Context) error {
		// 	return NewBandwidthPoolsCommand(ui, session).Run(c)
		// },
		"account-cancel-item": func(c *cli.Context) error {
			return NewCancelItemCommand(ui, accountManager).Run(c)
		},
		// "account-billing-items": func(c *cli.Context) error {
		// 	return NewBillingItemsCommand(ui, accountManager).Run(c)
		// },
		"account-invoice-detail": func(c *cli.Context) error {
			return NewInvoiceDetailCommand(ui, accountManager).Run(c)
		},
		"account-events": func(c *cli.Context) error {
			return NewEventsCommand(ui, accountManager).Run(c)
		},
		"account-event-detail": func(c *cli.Context) error {
			return NewEventDetailCommand(ui, accountManager).Run(c)
		},
		"account-invoices": func(c *cli.Context) error {
			return NewInvoicesCommand(ui, accountManager).Run(c)
		},
		"account-item-detail": func(c *cli.Context) error {
			return NewItemDetailCommand(ui, accountManager).Run(c)
		},
		"account-licenses": func(c *cli.Context) error {
			return NewLicensesCommand(ui, accountManager).Run(c)
		},
		"account-orders": func(c *cli.Context) error {
			return NewOrdersCommand(ui, accountManager).Run(c)
		},
		"account-summary": func(c *cli.Context) error {
			return NewSummaryCommand(ui, accountManager).Run(c)
		},
		// "account-bandwidth-pools-detail": func(c *cli.Context) error {
		// 	return NewBandwidthPoolsDetailCommand(ui, accountManager).Run(c)
		// },
	}

	return CommandActionBindings
}

func AccountNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "account",
		Description: T("Classic infrastructure Account commands"),
	}
}

/*
func AccountMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "account",
		Description: T("Classic infrastructure Account commands"),
		Usage:       "${COMMAND_NAME} sl account",
		Subcommands: []cli.Command{
			BandwidthPoolsMetaData(),
			CancelItemMetaData(),
			BillingItemsMetaData(),
			InvoiceDetailMetaData(),
			EventsMetaData(),
			EventDetailMetaData(),
			InvoicesMetaData(),
			ItemDetailMetaData(),
			LicensesMetaData(),
			OrdersMetaData(),
			SummaryMetaData(),
			BandwidthPoolsDetailMetaData(),
		},
	}
}
*/

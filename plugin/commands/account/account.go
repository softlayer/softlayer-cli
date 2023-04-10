package account

import (
	"github.com/spf13/cobra"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "account",
		Short: T("Classic infrastructure Account commands"),
		RunE:  nil,
	}
	cobraCmd.AddCommand(NewBandwidthPoolsCommand(sl).Command)
	cobraCmd.AddCommand(NewBandwidthPoolsDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewBillingItemsCommand(sl).Command)
	cobraCmd.AddCommand(NewCancelItemCommand(sl).Command)
	cobraCmd.AddCommand(NewInvoiceDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewEventsCommand(sl).Command)
	cobraCmd.AddCommand(NewEventDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewInvoicesCommand(sl).Command)
	cobraCmd.AddCommand(NewItemDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewLicensesCommand(sl).Command)
	cobraCmd.AddCommand(NewOrdersCommand(sl).Command)
	cobraCmd.AddCommand(NewSummaryCommand(sl).Command)
	cobraCmd.AddCommand(NewHooksCommand(sl).Command)
	cobraCmd.AddCommand(NewHookCreateCommand(sl).Command)
	cobraCmd.AddCommand(NewHookDeleteCommand(sl).Command)
	return cobraCmd
}

func AccountNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "account",
		Description: T("Classic infrastructure Account commands"),
	}
}

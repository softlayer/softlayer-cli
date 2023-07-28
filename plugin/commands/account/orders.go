package account

import (
	"bytes"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type OrdersCommand struct {
	*metadata.SoftlayerCommand
	AccountManager managers.AccountManager
	Command        *cobra.Command
	Limit          int
	Upgrades       bool
}

func NewOrdersCommand(sl *metadata.SoftlayerCommand) *OrdersCommand {
	thisCmd := &OrdersCommand{
		SoftlayerCommand: sl,
		AccountManager:   managers.NewAccountManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "orders",
		Short: T("Lists account orders. Use `slcli order lookup <ID>` to find more details about a specific order."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.Limit, "limit", 100, T("How many results to get in one api call."))
	cobraCmd.Flags().BoolVar(&thisCmd.Upgrades, "upgrades", false, T("Show upgrades orders."))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *OrdersCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	mask := "mask[orderTotalAmount,userRecord,initialInvoice[id,amount,invoiceTotalAmount],items[description]]"
	orders, err := cmd.AccountManager.GetAccountAllBillingOrders(mask, cmd.Limit)
	if err != nil {
		return errors.NewAPIError(T("Failed to get orders."), err.Error(), 2)
	}
	PrintOrders(orders, cmd.UI, outputFormat)

	if cmd.Upgrades {
		mask = "mask[id,maintenanceStartTimeUtc,statusId,createDate,ticketId]"
		upgrades, err := cmd.AccountManager.GetUpgradeRequests(mask, cmd.Limit)
		if err != nil {
			return errors.NewAPIError(T("Failed to get Upgrade Requests."), err.Error(), 2)
		}
		PrintUpgrades(upgrades, cmd.UI, outputFormat)
	}

	return nil
}

func PrintOrders(orders []datatypes.Billing_Order, ui terminal.UI, outputFormat string) {
	bufEvent := new(bytes.Buffer)
	table := terminal.NewTable(bufEvent, []string{
		T("Id"),
		T("State"),
		T("User"),
		T("Date"),
		T("Amount"),
		T("Item"),
	})

	for _, order := range orders {

		items := []string{}
		for _, item := range order.Items {
			items = append(items, utils.FormatStringPointer(item.Description))
		}

		table.Add(
			utils.FormatIntPointer(order.Id),
			utils.FormatStringPointer(order.Status),
			utils.FormatStringPointer(order.UserRecord.Username),
			utils.FormatSLTimePointer(order.CreateDate),
			utils.FormatSLFloatPointerToFloat(order.OrderTotalAmount),
			utils.ShortenString(strings.Join(items[:], ",")),
		)
	}

	utils.PrintTableWithTitle(ui, table, bufEvent, "Orders", outputFormat)
}

func PrintUpgrades(upgrades []datatypes.Product_Upgrade_Request, ui terminal.UI, outputFormat string) {
	bufEvent := new(bytes.Buffer)
	table := terminal.NewTable(bufEvent, []string{
		T("Id"),
		T("Maintance Window"),
		T("Status"),
		T("Created Date"),
		T("Case"),
	})

	for _, upgrade := range upgrades {
		table.Add(
			utils.FormatIntPointer(upgrade.Id),
			utils.FormatSLTimePointer(upgrade.MaintenanceStartTimeUtc),
			utils.FormatIntPointer(upgrade.StatusId),
			utils.FormatSLTimePointer(upgrade.CreateDate),
			utils.FormatIntPointer(upgrade.TicketId),
		)
	}

	utils.PrintTableWithTitle(ui, table, bufEvent, "Upgrade orders", outputFormat)
}

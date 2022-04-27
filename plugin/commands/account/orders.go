package account

import (
	"bytes"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type OrdersCommand struct {
	UI             terminal.UI
	AccountManager managers.AccountManager
}

func NewOrdersCommand(ui terminal.UI, accountManager managers.AccountManager) (cmd *OrdersCommand) {
	return &OrdersCommand{
		UI:             ui,
		AccountManager: accountManager,
	}
}

func OrdersMetaData() cli.Command {
	return cli.Command{
		Category:    "account",
		Name:        "orders",
		Description: T("Lists account orders."),
		Usage:       T(`${COMMAND_NAME} sl account orders [OPTIONS]`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "limit",
				Usage: T("How many results to get in one api call. [default: 50]"),
			},
			metadata.OutputFlag(),
		},
	}
}

func (cmd *OrdersCommand) Run(c *cli.Context) error {
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	limit := 50
	if c.IsSet("limit") {
		limit = c.Int("limit")
	}

	mask := "mask[orderTotalAmount,userRecord,initialInvoice[id,amount,invoiceTotalAmount],items[description]]"
	orders, err := cmd.AccountManager.GetAccountAllBillingOrders(mask, limit)
	if err != nil {
		return cli.NewExitError(T("Failed to get orders.")+err.Error(), 2)
	}
	PrintOrders(orders, cmd.UI, outputFormat)

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

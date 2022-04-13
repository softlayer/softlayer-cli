package account

import (
	"bytes"
	"fmt"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type BillingItemsCommand struct {
	UI             terminal.UI
	AccountManager managers.AccountManager
}

func NewBillingItemsCommand(ui terminal.UI, accountManager managers.AccountManager) (cmd *BillingItemsCommand) {
	return &BillingItemsCommand{
		UI:             ui,
		AccountManager: accountManager,
	}
}

func BillingItemsMetaData() cli.Command {
	return cli.Command{
		Category:    "account",
		Name:        "billing-items",
		Description: T("Lists billing items with some other useful information."),
		Usage:       T(`${COMMAND_NAME} slcli account billing-items [OPTIONS]`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

func (cmd *BillingItemsCommand) Run(c *cli.Context) error {

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	mask := "mask[orderItem[id,order[id,userRecord[id,email,displayName,userStatus]]],nextInvoiceTotalRecurringAmount,location, hourlyFlag]"
	billingItems, err := cmd.AccountManager.GetBillingItems(mask)
	if err != nil {
		return cli.NewExitError(T("Failed to get billing items.")+err.Error(), 2)
	}
	PrintBillingItems(billingItems, cmd.UI, outputFormat)
	return nil
}

func PrintBillingItems(billingItems []datatypes.Billing_Item, ui terminal.UI, outputFormat string) {
	bufEvent := new(bytes.Buffer)
	table := terminal.NewTable(bufEvent, []string{
		T("Id"),
		T("Create Date"),
		T("Cost"),
		T("Category Code"),
		T("Ordered By"),
		T("Description"),
		T("Notes"),
	})
	for _, billingItems := range billingItems {
		Description := billingItems.Description
		fqdn := fmt.Sprintf("%s.%s", utils.FormatStringPointerName(billingItems.HostName), utils.FormatStringPointerName(billingItems.DomainName))
		if fqdn != "." {
			Description = &fqdn
		}
		OrderedBy := "IBM"

		if billingItems.OrderItem != nil {
			OrderedBy = utils.FormatStringPointer(billingItems.OrderItem.Order.UserRecord.DisplayName)
		}
		table.Add(
			utils.FormatIntPointer(billingItems.Id),
			utils.FormatSLTimePointer(billingItems.CreateDate),
			fmt.Sprintf("%.2f", *billingItems.NextInvoiceTotalRecurringAmount),
			utils.FormatStringPointer(billingItems.CategoryCode),
			utils.FormatStringPointer(&OrderedBy),
			utils.ShortenStringWithLimit(utils.FormatStringPointer(Description), 50),
			utils.ShortenStringWithLimit(utils.FormatStringPointer(billingItems.Notes), 50),
		)
	}
	utils.PrintTableWithTitle(ui, table, bufEvent, "Billing Items", outputFormat)
}

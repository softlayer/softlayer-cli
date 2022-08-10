package account

import (
	"bytes"
	"fmt"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type BillingItemsCommand struct {
	*metadata.SoftlayerCommand
	AccountManager managers.AccountManager
	Command *cobra.Command
}

func NewBillingItemsCommand(sl *metadata.SoftlayerCommand) *BillingItemsCommand {
	thisCmd := &BillingItemsCommand{
		SoftlayerCommand: sl,
		AccountManager: managers.NewAccountManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use: "billing-items",
		Short: T("Lists billing items with some other useful information."),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}


func (cmd *BillingItemsCommand)  Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	mask := "mask[orderItem[id,order[id,userRecord[id,email,displayName,userStatus]]],nextInvoiceTotalRecurringAmount,location, hourlyFlag]"
	billingItems, err := cmd.AccountManager.GetBillingItems(mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get billing items."), err.Error(), 2)
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

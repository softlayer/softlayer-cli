package account

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ItemDetailCommand struct {
	*metadata.SoftlayerCommand
	AccountManager managers.AccountManager
	Command        *cobra.Command
}

func NewItemDetailCommand(sl *metadata.SoftlayerCommand) *ItemDetailCommand {
	thisCmd := &ItemDetailCommand{
		SoftlayerCommand: sl,
		AccountManager:   managers.NewAccountManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "item-detail " + T("IDENTIFIER"),
		Short: T("Gets detailed information about a billing item."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ItemDetailCommand) Run(args []string) error {
	itemID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Item ID")
	}

	outputFormat := cmd.GetOutputFlag()

	mask := `mask[
orderItem[id,order[id,userRecord[id,email,displayName,userStatus]]],
nextInvoiceTotalRecurringAmount,location, hourlyFlag, children
]`
	// The ID entered is from 'account billing-items'
	item, err := cmd.AccountManager.GetItemDetail(itemID, mask)
	if err != nil {
		// ID entered might be from 'account invoice-detail <ID>'
		item, err = cmd.AccountManager.GetItemDetailFromInvoiceItem(itemID, mask)
		if err != nil {
			subs := map[string]interface{}{"itemID": itemID}
			return slErr.NewAPIError(T("Failed to get the item {{.itemID}}. ", subs), err.Error(), 2)
		}

	}
	PrintItemDetail(itemID, item, cmd.UI, outputFormat)
	return nil
}

func PrintItemDetail(itemID int, item datatypes.Billing_Item, ui terminal.UI, outputFormat string) {
	bufEvent := new(bytes.Buffer)
	table := terminal.NewTable(bufEvent, []string{
		T("Key"),
		T("Value"),
	})

	table.Add("createDate", utils.FormatSLTimePointer(item.CreateDate))
	table.Add("cycleStartDate", utils.FormatSLTimePointer(item.CycleStartDate))
	table.Add("cancellationDate", utils.FormatSLTimePointer(item.CancellationDate))
	table.Add("description", utils.FormatStringPointer(item.Description))

	fqdn := fmt.Sprintf("%s.%s", utils.FormatStringPointerName(item.HostName), utils.FormatStringPointerName(item.DomainName))
	if fqdn != "." {
		table.Add("FQDN", fqdn)
	}
	if utils.FormatBoolPointer(item.HourlyFlag) == "true" {
		table.Add("hourlyRecurringFee", utils.FormatSLFloatPointerToFloat(item.HourlyRecurringFee))
		table.Add("hoursUsed", utils.FormatStringPointer(item.HoursUsed))
		table.Add("currentHourlyCharge", utils.FormatStringPointer(item.CurrentHourlyCharge))
	} else {
		table.Add("recurringFee", utils.FormatSLFloatPointerToFloat(item.RecurringFee))
	}
	OrderedBy := "IBM"
	if item.OrderItem != nil {
		OrderedBy = fmt.Sprintf("%s (%s)", utils.FormatStringPointer(item.OrderItem.Order.UserRecord.DisplayName), utils.FormatStringPointer(item.OrderItem.Order.UserRecord.UserStatus.Name))
	}
	table.Add("Ordered By", OrderedBy)

	table.Add("Notes", utils.FormatStringPointer(item.Notes))
	if item.Location != nil {
		table.Add("Location", utils.FormatStringPointer(item.Location.Name))
	}
	if item.Children != nil {
		for _, child := range item.Children {
			table.Add(utils.FormatStringPointer(child.CategoryCode), utils.FormatStringPointer(child.Description))
		}
	}

	utils.PrintTableWithTitle(ui, table, bufEvent, utils.FormatStringPointer(item.Description), outputFormat)
}

package order

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type LookupCommand struct {
	*metadata.SoftlayerCommand
	OrderManager managers.OrderManager
	Command      *cobra.Command
	Details      bool
}

func NewLookupCommand(sl *metadata.SoftlayerCommand) (cmd *LookupCommand) {
	thisCmd := &LookupCommand{
		SoftlayerCommand: sl,
		OrderManager:     managers.NewOrderManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "lookup " + T("IDENTIFIER"),
		Short: T("Provides some details related to order owner, date order, cost information, initial invoice."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVar(&thisCmd.Details, "details", false, T("Shows a very detailed list of charges"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *LookupCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	orderId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Order ID")
	}

	order, err := cmd.OrderManager.GetOrderDetail(orderId, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get Order."), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(order.Id))

	orderedBy := "IBM"
	if order.UserRecord != nil {
		orderedBy = fmt.Sprintf("%s (%s)", *order.UserRecord.DisplayName, *order.UserRecord.UserStatus.Name)
	}
	table.Add(T("Ordered By"), orderedBy)

	table.Add(T("Create Date"), utils.FormatSLTimePointer(order.CreateDate))
	table.Add(T("Modify Date"), utils.FormatSLTimePointer(order.ModifyDate))
	table.Add(T("Order Approval Date"), utils.FormatSLTimePointer(order.OrderApprovalDate))
	table.Add(T("Status"), utils.FormatStringPointer(order.Status))
	table.Add(T("Order Total Amount"), fmt.Sprintf("%.2f", float64(*order.OrderTotalAmount)))
	table.Add(T("Invoice Total Amount"), fmt.Sprintf("%.2f", float64(*order.InitialInvoice.InvoiceTotalAmount)))

	buf := new(bytes.Buffer)
	itemsTable := terminal.NewTable(buf, []string{T("Item Description")})
	for _, item := range order.Items {
		itemsTable.Add(utils.FormatStringPointer(item.Description))
	}
	itemsTable.Print()
	table.Add(T("Items"), buf.String())

	buf = new(bytes.Buffer)
	invoiceTable := terminal.NewTable(buf, []string{
		T("Item Id"),
		T("Category"),
		T("Description"),
		T("Single"),
		T("Monthly"),
		T("Create Date"),
		T("Location"),
	})
	for _, invoiceDetail := range order.InitialInvoice.InvoiceTopLevelItems {
		Category := utils.FormatStringPointerName(invoiceDetail.Category.Name)
		if Category == "" {
			Category = utils.FormatStringPointer(invoiceDetail.CategoryCode)
		}
		fqdn := fmt.Sprintf("%s.%s", utils.FormatStringPointerName(invoiceDetail.HostName), utils.FormatStringPointerName(invoiceDetail.DomainName))
		Description := utils.FormatStringPointer(invoiceDetail.Description)
		if fqdn != "." {
			Description = fmt.Sprintf("%s (%s)", Description, fqdn)
		}
		invoiceTable.Add(
			utils.FormatIntPointer(invoiceDetail.Id),
			Category,
			utils.ShortenString(Description),
			fmt.Sprintf("%.2f", *invoiceDetail.OneTimeAfterTaxAmount),
			fmt.Sprintf("%.2f", *invoiceDetail.RecurringAfterTaxAmount),
			utils.FormatSLTimePointer(invoiceDetail.CreateDate),
			utils.FormatStringPointer(invoiceDetail.Location.Name),
		)
		if cmd.Details {
			for _, child := range invoiceDetail.Children {
				invoiceTable.Add(
					">>>",
					utils.FormatStringPointer(child.Category.Name),
					utils.ShortenString(utils.FormatStringPointer(child.Description)),
					fmt.Sprintf("%.2f", *invoiceDetail.OneTimeAfterTaxAmount),
					fmt.Sprintf("%.2f", *invoiceDetail.RecurringAfterTaxAmount),
					"---",
					"---",
				)
			}
		}
	}
	invoiceTable.Print()
	table.Add(T("Initial Invoice"), buf.String())

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

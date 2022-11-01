package account

import (
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

type InvoiceDetailCommand struct {
	*metadata.SoftlayerCommand
	AccountManager managers.AccountManager
	Command        *cobra.Command
	Details        bool
}

func NewInvoiceDetailCommand(sl *metadata.SoftlayerCommand) *InvoiceDetailCommand {
	thisCmd := &InvoiceDetailCommand{
		SoftlayerCommand: sl,
		AccountManager:   managers.NewAccountManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "invoice-detail " + T("IDENTIFIER"),
		Short: T("Invoice details."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVar(&thisCmd.Details, "details", false, T("Shows a very detailed list of charges"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *InvoiceDetailCommand) Run(args []string) error {

	invoiceID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Invoice ID")
	}

	outputFormat := cmd.GetOutputFlag()

	mask := `mask[id, description, hostName, domainName, oneTimeAfterTaxAmount, recurringAfterTaxAmount, ` +
		`createDate,categoryCode,category[name],location[name],children[id, category[name], description, ` +
		`oneTimeAfterTaxAmount, recurringAfterTaxAmount]]`
	invoice, err := cmd.AccountManager.GetInvoiceDetail(invoiceID, mask)
	if err != nil {
		subs := map[string]interface{}{"invoiceID": invoiceID}
		return slErr.NewAPIError(T("Failed to get the invoice {{.invoiceID}}. ", subs), err.Error(), 2)
	}
	PrintInvoiceDetail(invoiceID, invoice, cmd.UI, outputFormat, cmd.Details)
	return nil
}

func PrintInvoiceDetail(invoiceID int, invoice []datatypes.Billing_Invoice_Item, ui terminal.UI, outputFormat string, details bool) {
	var rows [][]string
	rows = append(rows, []string{
		T("Item Id"),
		T("Category"),
		T("Description"),
		T("Single"),
		T("Monthly"),
		T("Create Date"),
		T("Location"),
	})

	for _, invoiceDetail := range invoice {
		Category := utils.FormatStringPointerName(invoiceDetail.Category.Name)
		if Category == "" {
			Category = utils.FormatStringPointer(invoiceDetail.CategoryCode)
		}
		fqdn := fmt.Sprintf("%s.%s", utils.FormatStringPointerName(invoiceDetail.HostName), utils.FormatStringPointerName(invoiceDetail.DomainName))
		Description := utils.FormatStringPointer(invoiceDetail.Description)
		if fqdn != "." {
			Description = fmt.Sprintf("%s (%s)", Description, fqdn)
		}
		rows = append(rows, []string{
			utils.FormatIntPointer(invoiceDetail.Id),
			Category,
			utils.ShortenString(Description),
			fmt.Sprintf("%.2f", *invoiceDetail.OneTimeAfterTaxAmount),
			fmt.Sprintf("%.2f", *invoiceDetail.RecurringAfterTaxAmount),
			utils.FormatSLTimePointer(invoiceDetail.CreateDate),
			utils.FormatStringPointer(invoiceDetail.Location.Name)},
		)
		if details {
			for _, child := range invoiceDetail.Children {
				rows = append(rows, []string{
					">>>",
					utils.FormatStringPointer(child.Category.Name),
					utils.ShortenString(utils.FormatStringPointer(child.Description)),
					fmt.Sprintf("%.2f", *invoiceDetail.OneTimeAfterTaxAmount),
					fmt.Sprintf("%.2f", *invoiceDetail.RecurringAfterTaxAmount),
					"---",
					"---"},
				)
			}
		}
	}
	utils.PrintTableWithTitleAndCSV(ui, rows, "Invoice: "+strconv.Itoa(invoiceID), outputFormat)
}

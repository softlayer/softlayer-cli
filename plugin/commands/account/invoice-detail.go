package account

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type InvoiceDetailCommand struct {
	UI             terminal.UI
	AccountManager managers.AccountManager
}

func NewInvoiceDetailCommand(ui terminal.UI, accountManager managers.AccountManager) (cmd *InvoiceDetailCommand) {
	return &InvoiceDetailCommand{
		UI:             ui,
		AccountManager: accountManager,
	}
}

func InvoiceDetailMetaData() cli.Command {
	return cli.Command{
		Category:    "account",
		Name:        "invoice-detail",
		Description: T("Invoice details."),
		Usage:       T(`${COMMAND_NAME} sl account invoice-detail IDENTIFIER [OPTIONS]`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "details",
				Usage: T("Shows a very detailed list of charges"),
			},
			metadata.OutputFlag(),
		},
	}
}

func (cmd *InvoiceDetailCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError("Invoice ID is required.")
	}

	invoiceID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return errors.NewInvalidUsageError(T("The invoice ID has to be a positive integer."))
	}

	details := false
	if c.IsSet("details") {
		details = true
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	invoice, err := cmd.AccountManager.GetInvoiceDetail(invoiceID)
	if err != nil {
		return cli.NewExitError(T("Failed to get the invoice {{.invoiceID}}. ", map[string]interface{}{"invoiceID": invoiceID})+err.Error(), 2)
	}
	PrintInvoiceDetail(invoiceID, invoice, cmd.UI, outputFormat, details)
	return nil
}

func PrintInvoiceDetail(invoiceID int, invoice []datatypes.Billing_Invoice_Item, ui terminal.UI, outputFormat string, details bool) {
	tableTitle := ui.Table([]string{
		T("Invoice: " + strconv.Itoa(invoiceID)),
	})

	bufEvent := new(bytes.Buffer)
	table := terminal.NewTable(bufEvent, []string{
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
		table.Add(
			utils.FormatIntPointer(invoiceDetail.Id),
			Category,
			NiceString(Description),
			fmt.Sprintf("%.2f", *invoiceDetail.OneTimeAfterTaxAmount),
			fmt.Sprintf("%.2f", *invoiceDetail.RecurringAfterTaxAmount),
			utils.FormatSLTimePointer(invoiceDetail.CreateDate),
			utils.FormatStringPointer(invoiceDetail.Location.Name),
		)
		if details {
			for _, child := range invoiceDetail.Children {
				table.Add(
					">>>",
					utils.FormatStringPointer(child.Category.Name),
					NiceString(utils.FormatStringPointer(child.Description)),
					fmt.Sprintf("%.2f", *invoiceDetail.OneTimeAfterTaxAmount),
					fmt.Sprintf("%.2f", *invoiceDetail.RecurringAfterTaxAmount),
					"---",
					"---",
				)
			}
		}
	}
	if outputFormat == "JSON" {
		table.PrintJson()
		tableTitle.Add(bufEvent.String())
		tableTitle.PrintJson()
	} else {
		table.Print()
		tableTitle.Add(bufEvent.String())
		tableTitle.Print()
	}
}

func NiceString(ugly_string string) string {
	limit := 100
	if len(ugly_string) > limit {
		return ugly_string[:limit] + ".."
	}
	return ugly_string
}

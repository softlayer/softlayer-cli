package account

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type InvoicesCommand struct {
	UI             terminal.UI
	AccountManager managers.AccountManager
}

func NewInvoicesCommand(ui terminal.UI, accountManager managers.AccountManager) (cmd *InvoicesCommand) {
	return &InvoicesCommand{
		UI:             ui,
		AccountManager: accountManager,
	}
}

func InvoicesMetaData() cli.Command {
	return cli.Command{
		Category:    "account",
		Name:        "invoices",
		Description: T("List invoices"),
		Usage:       T(`${COMMAND_NAME} sl account invoices [OPTIONS]`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "limit",
				Usage: T("How many invoices to get back. [default: 50]"),
			},
			cli.BoolFlag{
				Name:  "closed",
				Usage: T("Include invoices with a CLOSED status. [default: False]"),
			},
			cli.BoolFlag{
				Name:  "all",
				Usage: T("Return ALL invoices. There may be a lot of these. [default: False]"),
			},
			metadata.OutputFlag(),
		},
	}
}

func (cmd *InvoicesCommand) Run(c *cli.Context) error {
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	limit := 50
	if c.IsSet("limit") {
		limit = c.Int("limit")
	}

	closed := false
	if c.IsSet("closed") {
		closed = true
	}

	all := false
	if c.IsSet("all") {
		all = true
	}

	invoices, err := cmd.AccountManager.GetInvoices(limit, closed, all)
	if err != nil {
		return cli.NewExitError(T("Failed to get invoices.")+err.Error(), 2)
	}
	table := cmd.UI.Table([]string{
		T("Id"),
		T("Created"),
		T("Type"),
		T("Status"),
		T("Starting Balance"),
		T("Ending Balance"),
		T("Invoice Amount"),
		T("Items"),
	})
	for _, invoice := range invoices {
		table.Add(
			utils.FormatIntPointer(invoice.Id),
			utils.FormatSLTimePointer(invoice.CreateDate),
			utils.FormatStringPointer(invoice.TypeCode),
			utils.FormatStringPointer(invoice.StatusCode),
			utils.FormatSLFloatPointerToFloat(invoice.StartingBalance),
			utils.FormatSLFloatPointerToFloat(invoice.EndingBalance),
			utils.FormatSLFloatPointerToFloat(invoice.InvoiceTotalAmount),
			utils.FormatUIntPointer(invoice.ItemCount),
		)
	}
	if outputFormat == "JSON" {
		table.PrintJson()
	} else {
		table.Print()
	}
	return nil
}

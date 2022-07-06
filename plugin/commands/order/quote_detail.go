package order

import (
	"bytes"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type QuoteDetailCommand struct {
	UI           terminal.UI
	OrderManager managers.OrderManager
}

func NewQuoteDetailCommand(ui terminal.UI, orderManager managers.OrderManager) (cmd *QuoteDetailCommand) {
	return &QuoteDetailCommand{
		UI:           ui,
		OrderManager: orderManager,
	}
}

func (cmd *QuoteDetailCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	quoteId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Quote ID")
	}

	quote, err := cmd.OrderManager.GetQuote(quoteId, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get Quote.\n"+err.Error()), 2)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("Id"), utils.FormatIntPointer(quote.Id))
	table.Add(T("Name"), utils.FormatStringPointer(quote.Name))
	table.Add(T("Package"), utils.FormatStringPointer(quote.Order.Items[0].Package.KeyName))

	buf := new(bytes.Buffer)
	itemsTable := terminal.NewTable(buf, []string{T("Category"), T("Description"), T("Quantity"), T("Recurring"), T("One Time")})
	for _, item := range quote.Order.Items {
		itemsTable.Add(
			utils.FormatStringPointer(item.CategoryCode),
			utils.FormatStringPointer(item.Description),
			utils.FormatIntPointer(item.Quantity),
			utils.FormatSLFloatPointerToFloat(item.RecurringFee),
			utils.FormatSLFloatPointerToFloat(item.OneTimeFee),
		)
	}
	itemsTable.Print()
	table.Add(T("Items"), buf.String())

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

func OrderQuoteDetailMetaData() cli.Command {
	return cli.Command{
		Category:    "order",
		Name:        "quote-detail",
		Description: T("View a quote"),
		Usage: T(`${COMMAND_NAME} sl order quote-detail IDENTIFIER [OPTIONS]

EXAMPLE: 
	${COMMAND_NAME} sl order quote-detail 123456`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

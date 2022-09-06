package order

import (
	"bytes"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type QuoteDetailCommand struct {
	*metadata.SoftlayerCommand
	OrderManager managers.OrderManager
	Command      *cobra.Command
}

func NewQuoteDetailCommand(sl *metadata.SoftlayerCommand) (cmd *QuoteDetailCommand) {
	thisCmd := &QuoteDetailCommand{
		SoftlayerCommand: sl,
		OrderManager:     managers.NewOrderManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "quote-detail " + T("IDENTIFIER"),
		Short: T("View a quote"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *QuoteDetailCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	quoteId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Quote ID")
	}

	quote, err := cmd.OrderManager.GetQuote(quoteId, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get Quote.\n"), err.Error(), 2)
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

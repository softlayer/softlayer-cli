package order

import (
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type QuoteListCommand struct {
	*metadata.SoftlayerCommand
	OrderManager managers.OrderManager
	Command      *cobra.Command
}

func NewQuoteListCommand(sl *metadata.SoftlayerCommand) (cmd *QuoteListCommand) {
	thisCmd := &QuoteListCommand{
		SoftlayerCommand: sl,
		OrderManager:     managers.NewOrderManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "quote-list",
		Short: T("List all active quotes on an account"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *QuoteListCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	quotes, err := cmd.OrderManager.GetActiveQuotes("")
	if err != nil {
		return errors.NewAPIError(T("Failed to get Quotes.\n"), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Id"), T("Name"), T("Created"), T("Expiration"), T("Status"), T("Package Name"), T("Package Id")})
	for _, quote := range quotes {
		table.Add(
			utils.FormatIntPointer(quote.Id),
			utils.FormatStringPointer(quote.Name),
			utils.FormatSLTimePointer(quote.CreateDate),
			utils.FormatSLTimePointer(quote.ModifyDate),
			utils.FormatStringPointer(quote.Status),
			utils.FormatStringPointer(quote.Order.Items[0].Package.KeyName),
			utils.FormatIntPointer(quote.Order.Items[0].Package.Id),
		)
	}
	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

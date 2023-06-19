package order

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type QuoteDeleteCommand struct {
	*metadata.SoftlayerCommand
	OrderManager managers.OrderManager
	Command      *cobra.Command
}

func NewQuoteDeleteCommand(sl *metadata.SoftlayerCommand) (cmd *QuoteDeleteCommand) {
	thisCmd := &QuoteDeleteCommand{
		SoftlayerCommand: sl,
		OrderManager:     managers.NewOrderManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "quote-delete " + T("IDENTIFIER"),
		Short: T("Delete the quote of an order."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *QuoteDeleteCommand) Run(args []string) error {

	quoteId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Quote ID")
	}

	deletedQuote, err := cmd.OrderManager.DeleteQuote(quoteId)
	if err != nil {
		return errors.NewAPIError(T("Failed to delete Quote"), err.Error(), 2)
	}

	i18nSub := map[string]interface{}{"quoteID": deletedQuote.Id}
	cmd.UI.Ok()
	cmd.UI.Print(T("Quote: {{.quoteID}} was deleted.", i18nSub))
	return nil
}

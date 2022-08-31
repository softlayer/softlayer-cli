package order

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type QuoteSaveCommand struct {
	*metadata.SoftlayerCommand
	OrderManager managers.OrderManager
	Command      *cobra.Command
}

func NewQuoteSaveCommand(sl *metadata.SoftlayerCommand) (cmd *QuoteSaveCommand) {
	thisCmd := &QuoteSaveCommand{
		SoftlayerCommand: sl,
		OrderManager:     managers.NewOrderManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "quote-save " + T("IDENTIFIER"),
		Short: T("Save a quote"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *QuoteSaveCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	quoteId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Quote ID")
	}

	quote, err := cmd.OrderManager.SaveQuote(quoteId)
	if err != nil {
		return errors.NewAPIError(T("Failed to save Quote.\n"), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Id"), T("Name"), T("Created"), T("Modified"), T("Status")})
	table.Add(
		utils.FormatIntPointer(quote.Id),
		utils.FormatStringPointer(quote.Name),
		utils.FormatSLTimePointer(quote.CreateDate),
		utils.FormatSLTimePointer(quote.ModifyDate),
		utils.FormatStringPointer(quote.Status),
	)

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

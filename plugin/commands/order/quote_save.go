package order

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type QuoteSaveCommand struct {
	UI           terminal.UI
	OrderManager managers.OrderManager
}

func NewQuoteSaveCommand(ui terminal.UI, orderManager managers.OrderManager) (cmd *QuoteSaveCommand) {
	return &QuoteSaveCommand{
		UI:           ui,
		OrderManager: orderManager,
	}
}

func (cmd *QuoteSaveCommand) Run(c *cli.Context) error {
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

	quote, err := cmd.OrderManager.SaveQuote(quoteId)
	if err != nil {
		return cli.NewExitError(T("Failed to save Quote.\n")+err.Error(), 2)
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

func OrderQuoteSaveMetaData() cli.Command {
	return cli.Command{
		Category:    "order",
		Name:        "quote-save",
		Description: T("Save a quote"),
		Usage: T(`${COMMAND_NAME} sl order quote-save IDENTIFIER [OPTIONS]

EXAMPLE: 
	${COMMAND_NAME} sl order quote-save 123456`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

package order

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type QuoteListCommand struct {
	UI           terminal.UI
	OrderManager managers.OrderManager
}

func NewQuoteListCommand(ui terminal.UI, orderManager managers.OrderManager) (cmd *QuoteListCommand) {
	return &QuoteListCommand{
		UI:           ui,
		OrderManager: orderManager,
	}
}

func (cmd *QuoteListCommand) Run(c *cli.Context) error {

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	quotes, err := cmd.OrderManager.GetActiveQuotes("")
	if err != nil {
		return cli.NewExitError(T("Failed to get Quotes.\n"+err.Error()), 2)
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

func OrderQuoteListMetaData() cli.Command {
	return cli.Command{
		Category:    "order",
		Name:        "quote-list",
		Description: T("List all active quotes on an account"),
		Usage: T(`${COMMAND_NAME} sl order quote-list [OPTIONS]

   EXAMPLE: 
	  ${COMMAND_NAME} sl order quote-list`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

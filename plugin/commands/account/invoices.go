package account

import (
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type InvoicesCommand struct {
	*metadata.SoftlayerCommand
	AccountManager managers.AccountManager
	Command        *cobra.Command
	Limit          int
	Closed         bool
	All            bool
}

func NewInvoicesCommand(sl *metadata.SoftlayerCommand) *InvoicesCommand {
	thisCmd := &InvoicesCommand{
		SoftlayerCommand: sl,
		AccountManager:   managers.NewAccountManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "invoices",
		Short: T("List invoices."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVar(&thisCmd.All, "all", false, T("Return ALL invoices. There may be a lot of these."))
	cobraCmd.Flags().BoolVar(&thisCmd.Closed, "closed", false, T("Include invoices with a CLOSED status."))
	cobraCmd.Flags().IntVar(&thisCmd.Limit, "limit", 50, T("How many invoices to get back."))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *InvoicesCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	invoices, err := cmd.AccountManager.GetInvoices(cmd.Limit, cmd.Closed, cmd.All)
	if err != nil {
		return errors.NewAPIError(T("Failed to get invoices."), err.Error(), 2)
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
	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

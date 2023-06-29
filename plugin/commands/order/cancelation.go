package order

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CancelationCommand struct {
	*metadata.SoftlayerCommand
	OrderManager managers.OrderManager
	Command      *cobra.Command
}

func NewCancelationCommand(sl *metadata.SoftlayerCommand) (cmd *CancelationCommand) {
	thisCmd := &CancelationCommand{
		SoftlayerCommand: sl,
		OrderManager:     managers.NewOrderManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "cancelation",
		Short: T("List cancelations."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CancelationCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	mask := ""
	items, err := cmd.OrderManager.GetAllCancelation(mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to list all item cancelations."), err.Error(), 2)
	}

	PrintItemsCancelation(items, cmd.UI, outputFormat)
	return nil
}

func PrintItemsCancelation(items []datatypes.Billing_Item_Cancellation_Request, ui terminal.UI, outputFormat string) {
	table := ui.Table([]string{
		T("Case Number"),
		T("Number Of Items Cancelled"),
		T("Created"),
		T("Status"),
		T("Requested by"),
	})

	for _, item := range items {
		requestedBy := utils.FormatStringPointer(item.User.FirstName) + " " + utils.FormatStringPointer(item.User.LastName)
		table.Add(
			utils.FormatIntPointer(item.TicketId),
			utils.FormatUIntPointer(item.ItemCount),
			utils.FormatSLTimePointer(item.CreateDate),
			utils.FormatStringPointer(item.Status.Name),
			requestedBy,
		)
	}
	utils.PrintTable(ui, table, outputFormat)
}

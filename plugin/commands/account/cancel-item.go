package account

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CancelItemCommand struct {
	*metadata.SoftlayerCommand
	Command        *cobra.Command
	AccountManager managers.AccountManager
}

func NewCancelItemCommand(sl *metadata.SoftlayerCommand) *CancelItemCommand {
	thisCmd := &CancelItemCommand{
		SoftlayerCommand: sl,
		AccountManager:   managers.NewAccountManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "cancel-item",
		Short: T("Cancels a billing item."),
		Long: T(`Cancel the resource or service for a billing Item. By default the billing item will be canceled
on the next bill date and reclaim of the resource will begin shortly after the cancellation`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CancelItemCommand) Run(args []string) error {

	itemID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Item ID")
	}

	err = cmd.AccountManager.CancelItem(itemID)
	itemIdSub := map[string]interface{}{"itemID": itemID}
	if err != nil {
		if strings.Contains(err.Error(), slErr.SL_EXP_OBJ_NOT_FOUND) {
			return slErr.NewAPIError(T("Unable to find item with ID: {{.itemID}}.\n", itemIdSub), err.Error(), 0)
		}
		return slErr.NewAPIError(T("Failed to cancel item: {{.itemID}}.\n", itemIdSub), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Item: {{.itemID}} was cancelled.", itemIdSub))
	return nil
}

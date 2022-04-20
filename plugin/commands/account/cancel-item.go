package account

import (
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type CancelItemCommand struct {
	UI             terminal.UI
	AccountManager managers.AccountManager
}

func NewCancelItemCommand(ui terminal.UI, accountManager managers.AccountManager) (cmd *CancelItemCommand) {
	return &CancelItemCommand{
		UI:             ui,
		AccountManager: accountManager,
	}
}

func CancelItemMetaData() cli.Command {
	return cli.Command{
		Category:    "account",
		Name:        "cancel-item",
		Description: T("Cancels a billing item."),
		Usage:       T(`${COMMAND_NAME} sl account cancel-item IDENTIFIER`),
		Flags:       []cli.Flag{},
	}
}

func (cmd *CancelItemCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return slErr.NewInvalidUsageError(T("This command requires one argument."))
	}

	itemID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Item ID")
	}

	err = cmd.AccountManager.CancelItem(itemID)
	if err != nil {
		if strings.Contains(err.Error(), slErr.SL_EXP_OBJ_NOT_FOUND) {
			return cli.NewExitError(T("Unable to find item with ID: {{.itemID}}.\n", map[string]interface{}{"itemID": itemID})+err.Error(), 0)
		}
		return cli.NewExitError(T("Failed to cancel item: {{.itemID}}.\n", map[string]interface{}{"itemID": itemID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Item: {{.itemID}} was cancelled.", map[string]interface{}{"itemID": itemID}))
	return nil
}

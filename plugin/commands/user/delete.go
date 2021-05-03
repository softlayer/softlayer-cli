package user

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
)

type DeleteCommand struct {
	UI          terminal.UI
	UserManager managers.UserManager
}

func NewDeleteCommand(ui terminal.UI, userManager managers.UserManager) (cmd *DeleteCommand) {
	return &DeleteCommand{
		UI:          ui,
		UserManager: userManager,
	}
}

func (cmd *DeleteCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	userId := c.Args()[0]

	id, err := strconv.Atoi(userId)
	if err != nil {
		return errors.NewInvalidUsageError(T("User ID should be a number."))
	}

	userStatusId := 1021
	templateObject := datatypes.User_Customer{
		UserStatusId: &userStatusId,
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This will delete the user: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": id}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	_, err = cmd.UserManager.EditUser(templateObject, id)
	if err != nil {
		cli.NewExitError(T("Failed to delete user.\n")+err.Error(), 2)
	}

	cmd.UI.Ok()
	return nil

}

package user

import (
	"strconv"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type DeleteCommand struct {
	*metadata.SoftlayerCommand
	UserManager managers.UserManager
	ForceFlag   bool
	Command     *cobra.Command
}

func NewDeleteCommand(sl *metadata.SoftlayerCommand) (cmd *DeleteCommand) {
	thisCmd := &DeleteCommand{
		SoftlayerCommand: sl,
		UserManager:      managers.NewUserManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "delete " + T("USER_ID"),
		Short: T("Sets a user's status to CANCEL_PENDING, which will immediately disable the account, and will eventually be fully removed from the account by an automated internal process"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVarP(&thisCmd.ForceFlag, "force", "f", false, T("Force operation without confirmation"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DeleteCommand) Run(args []string) error {
	userId := args[0]

	id, err := strconv.Atoi(userId)
	if err != nil {
		return errors.NewInvalidUsageError(T("User ID should be a number."))
	}

	userStatusId := 1021
	templateObject := datatypes.User_Customer{
		UserStatusId: &userStatusId,
	}

	if !cmd.ForceFlag {
		confirm, err := cmd.UI.Confirm(T("This will delete the user: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": id}))
		if err != nil {
			return errors.NewInvalidUsageError(err.Error())
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	_, err = cmd.UserManager.EditUser(templateObject, id)
	if err != nil {
		return errors.NewAPIError(T("Failed to delete user.\n"), err.Error(), 2)
	}

	cmd.UI.Ok()
	return nil

}

package user

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type EditPermissionCommand struct {
	*metadata.SoftlayerCommand
	UserManager managers.UserManager
	Command     *cobra.Command
	Enable      string
	Permission  []string
	FromUser    int
}

func NewEditPermissionCommand(sl *metadata.SoftlayerCommand) (cmd *EditPermissionCommand) {
	thisCmd := &EditPermissionCommand{
		SoftlayerCommand: sl,
		UserManager:      managers.NewUserManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "permission-edit " + T("IDENTIFIER"),
		Short: T("Enable or Disable specific permissions"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Enable, "enable", "true", T("Enable or Disable selected permissions. Accepted inputs are 'true' and 'false'. default is 'true'"))
	cobraCmd.Flags().StringSliceVar(&thisCmd.Permission, "permission", []string{}, T("Permission keyName to set. Use keyword ALL to select ALL permissions"))
	cobraCmd.Flags().IntVar(&thisCmd.FromUser, "from-user", 0, T("Set permissions to match this user's permissions. Adds and removes the appropriate permissions"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *EditPermissionCommand) Run(args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidUsageError(T("User ID should be a number."))
	}

	if (cmd.FromUser != 0 && len(cmd.Permission) != 0) || (cmd.FromUser == 0 && len(cmd.Permission) == 0) {
		return errors.NewInvalidUsageError(T("one of --permission and --from-user should be used to specify permissions"))
	}

	permissionKeynames := cmd.Permission
	permissions, err := cmd.UserManager.FormatPermissionObject(permissionKeynames)
	if err != nil {
		return err
	}

	if cmd.FromUser != 0 {
		fromUser := cmd.FromUser
		err = cmd.UserManager.PermissionFromUser(id, fromUser)
	} else if strings.ToLower(cmd.Enable) == "true" {
		_, err = cmd.UserManager.AddPermission(id, permissions)
	} else if strings.ToLower(cmd.Enable) == "false" {
		_, err = cmd.UserManager.RemovePermission(id, permissions)
	} else {
		return errors.NewInvalidUsageError(T("options for --enable are true, false"))
	}

	if err != nil {
		return errors.NewAPIError(T("Failed to update permissions"), err.Error(), 1)
	}
	cmd.UI.Print(T("Permissions updated successfully"))
	return nil
}

package user

import (
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type VpnEnableCommand struct {
	*metadata.SoftlayerCommand
	UserManager managers.UserManager
	Command     *cobra.Command
}

func NewVpnEnableCommand(sl *metadata.SoftlayerCommand) (cmd *VpnEnableCommand) {
	thisCmd := &VpnEnableCommand{
		SoftlayerCommand: sl,
		UserManager:      managers.NewUserManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "vpn-enable " + T("USER_ID"),
		Short: T("Enable vpn for a user."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VpnEnableCommand) Run(args []string) error {

	userID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("User ID")
	}

	userTemplate := datatypes.User_Customer{
		SslVpnAllowedFlag: sl.Bool(true),
	}

	subs := map[string]interface{}{"UserId": userID}
	confirm, err := cmd.UI.Confirm(T("This will Enable vpn for the user id: {{.UserId}}. Continue?", subs))
	if err != nil {
		return err
	}
	if !confirm {
		cmd.UI.Print(T("Aborted."))
		return nil
	}

	success, err := cmd.UserManager.EditUser(userTemplate, userID)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to enable vpn for the 'User ID' \n", subs), err.Error(), 2)
	}

	if success {
		cmd.UI.Ok()
		cmd.UI.Print(T("{{.UserId}} is successfully enabled.", subs))
	}
	return nil
}

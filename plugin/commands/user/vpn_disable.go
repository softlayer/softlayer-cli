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

type VpnDisableCommand struct {
	*metadata.SoftlayerCommand
	UserManager managers.UserManager
	Command     *cobra.Command
}

func NewVpnDisableCommand(sl *metadata.SoftlayerCommand) (cmd *VpnDisableCommand) {
	thisCmd := &VpnDisableCommand{
		SoftlayerCommand: sl,
		UserManager:      managers.NewUserManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "vpn-disable " + T("USER_ID"),
		Short: T("Disable vpn for a user."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VpnDisableCommand) Run(args []string) error {
	userID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("User ID")
	}
	userTemplate := datatypes.User_Customer{
		SslVpnAllowedFlag: sl.Bool(false),
	}
	success, err := cmd.UserManager.EditUser(userTemplate, userID)
	if err != nil {
		return slErrors.NewAPIError(T(""), err.Error(), 2)
	}
	if success {
		cmd.UI.Ok()
	}
	return nil
}

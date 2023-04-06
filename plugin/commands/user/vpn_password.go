package user

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type VpnPasswordCommand struct {
	*metadata.SoftlayerCommand
	UserManager managers.UserManager
	Command     *cobra.Command
	Password    string
}

func NewVpnPasswordCommand(sl *metadata.SoftlayerCommand) (cmd *VpnPasswordCommand) {
	thisCmd := &VpnPasswordCommand{
		SoftlayerCommand: sl,
		UserManager:      managers.NewUserManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "vpn-password " + T("IDENTIFIER"),
		Short: T("Set the user VPN password."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Password, "password", "", T("Your new VPN password [required]"))

	//#nosec G104 -- This is a false positive
	cobraCmd.MarkFlagRequired("password")

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VpnPasswordCommand) Run(args []string) error {
	userID, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("User ID")
	}

	success, err := cmd.UserManager.UpdateVpnPassword(userID, cmd.Password)
	if err != nil {
		return errors.NewAPIError(T("Failed to update user vpn password."), err.Error(), 2)
	}
	if success {
		cmd.UI.Ok()
		cmd.UI.Print(T("Successfully updated user vpn."))
	}
	return nil
}

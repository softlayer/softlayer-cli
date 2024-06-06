package user

import (
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type VpnManualCommand struct {
	*metadata.SoftlayerCommand
	UserManager managers.UserManager
	Command     *cobra.Command
	Enable      bool
	Disable     bool
}

func NewVpnManualCommand(sl *metadata.SoftlayerCommand) (cmd *VpnManualCommand) {
	thisCmd := &VpnManualCommand{
		SoftlayerCommand: sl,
		UserManager:      managers.NewUserManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "vpn-manual " + T("USER_ID"),
		Short: T("Enable or disable user vpn subnets manual config."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVar(&thisCmd.Enable, "enable", false, T("Enable vpn subnets manual config."))
	cobraCmd.Flags().BoolVar(&thisCmd.Disable, "disable", false, T("Disable vpn subnets manual config."))
	cobraCmd.MarkFlagsMutuallyExclusive("enable", "disable")

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VpnManualCommand) Run(args []string) error {

	if !cmd.Enable && !cmd.Disable {
		return errors.NewInvalidUsageError(T("This command requires --enable or --disable option."))
	}

	userID, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("User ID")
	}

	vpnManualConfig := false
	if cmd.Enable {
		vpnManualConfig = true
	}

	userTemplate := datatypes.User_Customer{
		VpnManualConfig: sl.Bool(vpnManualConfig),
	}


	success, err := cmd.UserManager.EditUser(userTemplate, userID)
	if err != nil {
		return errors.NewInvalidUsageError(T("Failed to update user vpn subnets manual config"))
	}
	if success {
		cmd.UI.Ok()
		cmd.UI.Print(T("Successfully updated user vpn subnets manual config"))
	}
	return nil
}

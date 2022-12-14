package user

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type VpnSubnetCommand struct {
	*metadata.SoftlayerCommand
	UserManager managers.UserManager
	Command     *cobra.Command
	Add         bool
	Remove      bool
}

func NewVpnSubnetCommand(sl *metadata.SoftlayerCommand) (cmd *VpnSubnetCommand) {
	thisCmd := &VpnSubnetCommand{
		SoftlayerCommand: sl,
		UserManager:      managers.NewUserManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "vpn-subnet " + T("USER_ID") + " " + T("SUBNET_ID"),
		Short: T("Add or remove subnet access for a user."),
		Args:  metadata.TwoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVar(&thisCmd.Add, "add", false, T("Add access to subnet."))
	cobraCmd.Flags().BoolVar(&thisCmd.Remove, "remove", false, T("Remove access to subnet."))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VpnSubnetCommand) Run(args []string) error {

	if !cmd.Add && !cmd.Remove {
		return errors.NewInvalidUsageError(T("This command requires --add or --remove option."))
	}

	if cmd.Add && cmd.Remove {
		return errors.NewExclusiveFlagsError("--add", "--remove")
	}

	userID, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("User ID")
	}

	subnetID, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Subnet ID")
	}

	if cmd.Add {
		overrideSuccess, overrideError := cmd.UserManager.CreateUserVpnOverride(userID, subnetID)
		if overrideError != nil {
			return errors.NewAPIError(T("Failed to create user vpn override."), overrideError.Error(), 2)
		}
		if overrideSuccess {
			success, err := cmd.UserManager.UpdateVpnUser(userID)
			if err != nil {
				return errors.NewAPIError(T("Override created, but unable to update VPN user."), err.Error(), 2)
			}
			if success {
				cmd.UI.Ok()
				cmd.UI.Print(T("Successfully added subnet access for user."))
			}
		}
	}

	if cmd.Remove {
		overrides, err := cmd.UserManager.GetOverrides(userID)
		if err != nil {
			return errors.NewAPIError(T("Failed to get user vpn overrides."), err.Error(), 2)
		}
		overrideId := 0
		for _, override := range overrides {
			if *override.SubnetId == subnetID {
				overrideId = *override.Id
			}
		}
		if overrideId == 0 {
			mapValues := map[string]interface{}{"subnetID": subnetID, "userID": userID}
			return errors.NewInvalidUsageError(T("Subnet {{.subnetID}} is not assigned to User {{.userID}}", mapValues))
		}
		overrideSuccess, overrideError := cmd.UserManager.DeleteUserVpnOverride(overrideId)
		if overrideError != nil {
			return errors.NewAPIError(T("Failed to delete user vpn override."), overrideError.Error(), 2)
		}
		if overrideSuccess {
			success, err := cmd.UserManager.UpdateVpnUser(userID)
			if err != nil {
				return errors.NewAPIError(T("Failed to update VPN user."), err.Error(), 2)
			}
			if success {
				cmd.UI.Ok()
				cmd.UI.Print(T("Successfully removed subnet access for user."))
			}
		}
	}
	return nil
}

package user

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type RemoveAccessCommand struct {
	*metadata.SoftlayerCommand
	UserManager managers.UserManager
	Command     *cobra.Command
	Hardware    string
	Virtual     string
	Dedicated   string
}

func NewRemoveAccessCommand(sl *metadata.SoftlayerCommand) (cmd *RemoveAccessCommand) {
	thisCmd := &RemoveAccessCommand{
		SoftlayerCommand: sl,
		UserManager:      managers.NewUserManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "remove-access " + T("IDENTIFIER"),
		Short: T("Remove access from a user to an specific device"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Hardware, "hardware", "", T("Hardware ID"))
	cobraCmd.Flags().StringVar(&thisCmd.Virtual, "virtual", "", T("Virtual Guest ID"))
	cobraCmd.Flags().StringVar(&thisCmd.Dedicated, "dedicated", "", T("Dedicated Host ID"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *RemoveAccessCommand) Run(args []string) error {
	if cmd.Hardware == "" && cmd.Virtual == "" && cmd.Dedicated == "" {
		return errors.NewInvalidUsageError(T("This command requires one option."))
	}

	userId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidUsageError(T("User ID should be a number."))
	}

	if cmd.Hardware != "" {
		hardwareId, err := strconv.Atoi(cmd.Hardware)
		if err != nil {
			return errors.NewInvalidUsageError(T("Hardware ID should be a number."))
		} else {
			i18nsubs := map[string]interface{}{"userId": userId, "objectId": hardwareId}
			response, err := cmd.UserManager.RemoveHardwareAccess(userId, hardwareId)
			if err != nil {
				return errors.NewAPIError(T("Failed to update access.\n"), err.Error(), 2)
			}
			if response {
				cmd.UI.Ok()
				cmd.UI.Print(T("Access removed to user {{.userId}} for {{.objectId}}", i18nsubs))
			}
		}
	}

	if cmd.Dedicated != "" {
		dedicatedHostId, err := strconv.Atoi(cmd.Dedicated)
		if err != nil {
			return errors.NewInvalidUsageError(T("Dedicated host ID should be a number."))
		} else {
			i18nsubs := map[string]interface{}{"userId": userId, "objectId": dedicatedHostId}
			response, err := cmd.UserManager.RemoveDedicatedHostAccess(userId, dedicatedHostId)
			if err != nil {
				return errors.NewAPIError(T("Failed to update access.\n"), err.Error(), 2)
			}
			if response {
				cmd.UI.Ok()
				cmd.UI.Print(T("Access removed to user {{.userId}} for {{.objectId}}", i18nsubs))
			}
		}
	}

	if cmd.Virtual != "" {
		virtualId, err := strconv.Atoi(cmd.Virtual)
		if err != nil {
			return errors.NewInvalidUsageError(T("Virtual server ID should be a number."))
		} else {
			i18nsubs := map[string]interface{}{"userId": userId, "objectId": virtualId}
			response, err := cmd.UserManager.RemoveVirtualGuestAccess(userId, virtualId)
			if err != nil {
				return errors.NewAPIError(T("Failed to update access.\n"), err.Error(), 2)
			}
			if response {
				cmd.UI.Ok()
				cmd.UI.Print(T("Access removed to user {{.userId}} for {{.objectId}}", i18nsubs))
			}
		}
	}

	return nil

}

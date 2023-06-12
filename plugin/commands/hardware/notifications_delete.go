package hardware

import (
	"strconv"

	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type NotificationsDeleteCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
}

func NewNotificationsDeleteCommand(sl *metadata.SoftlayerCommand) (cmd *NotificationsDeleteCommand) {
	thisCmd := &NotificationsDeleteCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "notifications-delete " + T("IDENTIFIER"),
		Short: T("Remove a user hardware notification entry."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *NotificationsDeleteCommand) Run(args []string) error {
	userCustomerNotificationId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	success, err := cmd.HardwareManager.DeleteUserCustomerNotification(userCustomerNotificationId)
	if err != nil {
		cmd.UI.Failed(T("Failed to delete User Customer Notification.") + "\n" + err.Error())
	}
	if success {
		cmd.UI.Ok()
		cmd.UI.Print(T("Successfully removed User Customer notification."))
	}

	return nil
}

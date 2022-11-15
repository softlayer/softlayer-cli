package hardware

import (
	"strconv"

	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type NotificationsCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
}

func NewNotificationsCommand(sl *metadata.SoftlayerCommand) (cmd *NotificationsCommand) {
	thisCmd := &NotificationsCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "notifications " + T("IDENTIFIER"),
		Short: T("Shows who gets notified when the server has a monitoring issues."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *NotificationsCommand) Run(args []string) error {
	hardwareId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat := cmd.GetOutputFlag()

	userCustomers, err := cmd.HardwareManager.GetUserCustomerNotificationsByHardwareId(hardwareId, "")
	if err != nil {
		return slErr.NewAPIError(T("Failed to get User Customer Notifications."), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("ID"), T("Last Name"), T("First Name"), T("Email"), T("User ID")})
	for _, userCustomer := range userCustomers {
		table.Add(
			utils.FormatIntPointer(userCustomer.Id),
			utils.FormatStringPointer(userCustomer.User.LastName),
			utils.FormatStringPointer(userCustomer.User.FirstName),
			utils.FormatStringPointer(userCustomer.User.Email),
			utils.FormatStringPointer(userCustomer.User.Username),
		)
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

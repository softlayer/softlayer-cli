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

type NotificationsAddCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
	Users           []int
}

func NewNotificationsAddCommand(sl *metadata.SoftlayerCommand) (cmd *NotificationsAddCommand) {
	thisCmd := &NotificationsAddCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "notifications-add " + T("IDENTIFIER"),
		Short: T("Create a user hardware notification entry."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().IntSliceVar(&thisCmd.Users, "users", []int{}, T("User ID to be notified on monitoring failure, multiple occurrence allowed"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *NotificationsAddCommand) Run(args []string) error {
	hardwareId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat := cmd.GetOutputFlag()

	userIds := cmd.Users
	printTable := false

	table := cmd.UI.Table([]string{T("Id"), T("Hostname"), T("Username"), T("Email"), T("First Name"), T("Last Name")})
	for _, userId := range userIds {

		UserCustomerNotification, err := cmd.HardwareManager.CreateUserCustomerNotification(hardwareId, userId)
		if err != nil {
			userIdMap := map[string]interface{}{"userID": userId}
			cmd.UI.Failed(T("Failed to create User Customer Notification with user ID: {{.userID}}."), userIdMap)
		} else {
			printTable = true
			table.Add(
				utils.FormatIntPointer(UserCustomerNotification.Id),
				utils.FormatStringPointer(UserCustomerNotification.Hardware.FullyQualifiedDomainName),
				utils.FormatStringPointer(UserCustomerNotification.User.Username),
				utils.FormatStringPointer(UserCustomerNotification.User.Email),
				utils.FormatStringPointer(UserCustomerNotification.User.FirstName),
				utils.FormatStringPointer(UserCustomerNotification.User.LastName),
			)
		}
	}

	if printTable {
		utils.PrintTable(cmd.UI, table, outputFormat)
	}
	return nil
}

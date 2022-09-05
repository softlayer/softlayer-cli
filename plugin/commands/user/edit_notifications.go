package user

import (
	"errors"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type EditNotificationsCommand struct {
	*metadata.SoftlayerCommand
	UserManager managers.UserManager
	Command     *cobra.Command
	Enable      []string
	Disable     []string
}

func NewEditNotificationsCommand(sl *metadata.SoftlayerCommand) (cmd *EditNotificationsCommand) {
	thisCmd := &EditNotificationsCommand{
		SoftlayerCommand: sl,
		UserManager:      managers.NewUserManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "edit-notifications",
		Short: T("Enable or Disable specific notifications for the active user."),
		Long: T(`
Notification names should be enclosed in quotation marks. Examples:
	${COMMAND_NAME} sl user edit-notifications --enable 'Order Approved'
	${COMMAND_NAME} sl user edit-notifications --enable 'Order Approved' --enable  'Reload Complete'`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringSliceVar(&thisCmd.Enable, "enable", []string{}, T("Enable (DEFAULT) selected notifications"))
	cobraCmd.Flags().StringSliceVar(&thisCmd.Disable, "disable", []string{}, T("Disable selected notifications"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *EditNotificationsCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	if len(cmd.Disable) != 0 && len(cmd.Enable) != 0 {
		return slErr.NewInvalidUsageError(T("Only set --enable or --disable options."))
	}

	if len(cmd.Disable) != 0 && len(args) > 0 {
		return slErr.NewInvalidUsageError(T("Only set --enable or --disable options."))
	}

	if len(cmd.Disable) == 0 && len(cmd.Enable) == 0 && (len(args) == 0) {
		return slErr.NewInvalidUsageError(T("This command requires notification names as arguments and options flags."))
	}

	notificationsInput := []string{}
	succesNotifications := []string{}
	failedNotifications := []string{}

	allNotifications, err := cmd.UserManager.GetAllNotifications("mask[id,name]")
	if err != nil {
		return slErr.NewAPIError(T("Failed to update notifications: "+printNotifications(notificationsInput)+"\n"), err.Error(), 2)
	}

	if len(cmd.Disable) != 0 {
		notificationsInput = cmd.Disable
		succesNotifications, failedNotifications = setNotifications(cmd, "disable", notificationsInput, allNotifications)
	}

	if len(cmd.Enable) != 0 {
		notificationsInput = cmd.Enable
		succesNotifications, failedNotifications = setNotifications(cmd, "enable", notificationsInput, allNotifications)
	}

	if len(cmd.Disable) == 0 && len(cmd.Disable) == 0 {
		notificationsInput = append(notificationsInput, args...)
		succesNotifications, failedNotifications = setNotifications(cmd, "enable", notificationsInput, allNotifications)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, len(failedNotifications) == 0)
	}

	if len(succesNotifications) > 0 {
		cmd.UI.Ok()
		cmd.UI.Print(T("Notifications updated successfully: " + printNotifications(succesNotifications)))
	}

	if len(failedNotifications) > 0 {
		cmd.UI.Print(T("Notifications updated unsuccessfully: " + printNotifications(failedNotifications) + ". Review if already set or if the name is correct."))
	}

	return nil
}

func printNotifications(notifications []string) string {
	res := ""
	for _, notification := range notifications {
		res = res + notification + ", "
	}
	return res[0 : len(res)-2]
}

func setNotifications(cmd *EditNotificationsCommand, action string, notificationsInput []string, allNotifications []datatypes.Email_Subscription) (succesNotifications []string, failedNotifications []string) {
	requestResponse := true
	errorRequest := errors.New("")

	for _, notificationName := range notificationsInput {
		for _, notification := range allNotifications {
			if *notification.Name == notificationName {
				if action == "enable" {
					requestResponse, errorRequest = cmd.UserManager.EnableEmailSubscriptionNotification(*notification.Id)
				} else if action == "disable" {
					requestResponse, errorRequest = cmd.UserManager.DisableEmailSubscriptionNotification(*notification.Id)
				}
				if errorRequest == nil && requestResponse {
					succesNotifications = append(succesNotifications, notificationName)
					break
				} else {
					requestResponse = false
				}
			} else {
				requestResponse = false
			}
		}
		if !requestResponse {
			failedNotifications = append(failedNotifications, notificationName)
		}
	}
	return succesNotifications, failedNotifications
}

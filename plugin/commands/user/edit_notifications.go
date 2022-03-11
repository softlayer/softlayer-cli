package user

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type EditNotificationsCommand struct {
	UI          terminal.UI
	UserManager managers.UserManager
}

func NewEditNotificationsCommand(ui terminal.UI, userManager managers.UserManager) (cmd *EditNotificationsCommand) {
	return &EditNotificationsCommand{
		UI:          ui,
		UserManager: userManager,
	}
}

func (cmd *EditNotificationsCommand) Run(c *cli.Context) error {

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	if c.IsSet("disable") && c.IsSet("enable") {
		return slErr.NewInvalidUsageError(T("Only set --enable or --disable options."))
	}

	if c.IsSet("disable") && c.NArg() > 0 {
		return slErr.NewInvalidUsageError(T("Only set --enable or --disable options."))
	}

	if !c.IsSet("disable") && !c.IsSet("enable") && (c.NArg() == 0) {
		return slErr.NewInvalidUsageError(T("This command requires notification names as arguments and options flags"))
	}

	notificationsInput := []string{}
	succesNotifications := []string{}
	failedNotifications := []string{}

	allNotifications, err := cmd.UserManager.GetAllNotifications("mask[id,name]")
	if err != nil {
		return cli.NewExitError(T("Failed to update notifications: "+printNotifications(notificationsInput)+"\n")+err.Error(), 2)
	}

	if c.IsSet("disable") {
		notificationsInput = c.StringSlice("disable")
		succesNotifications, failedNotifications = setNotifications(cmd, "disable", notificationsInput, allNotifications)
	}

	if c.IsSet("enable") {
		notificationsInput = c.StringSlice("enable")
		succesNotifications, failedNotifications = setNotifications(cmd, "enable", notificationsInput, allNotifications)
	}

	if !c.IsSet("disable") && !c.IsSet("enable") {
		notificationsInput = append(notificationsInput, c.Args()...)
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
		cmd.UI.Print(T("Notifications updated unsuccessfully: " + printNotifications(failedNotifications) + ". Review if already set or if the name is correct"))
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

func UserEditNotificationsMetaData() cli.Command {
	return cli.Command{
		Category:    "user",
		Name:        "edit-notifications",
		Description: T("Enable or Disable specific notifications for the active user."),
		Usage: T(`${COMMAND_NAME} sl user edit-notifications [OPTIONS] NOTIFICATIONS

		Notification names should be enclosed in quotation marks. Examples:
			slcli user edit-notifications --enable 'Order Approved'
			slcli user edit-notifications --enable 'Order Approved' --enable  'Reload Complete'`),
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "enable",
				Usage: T("Enable (DEFAULT) selected notifications"),
			},
			cli.StringSliceFlag{
				Name:  "disable",
				Usage: T("Disable selected notifications"),
			},
			metadata.OutputFlag(),
		},
	}
}

package user

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type NotificationsCommand struct {
	UI          terminal.UI
	UserManager managers.UserManager
}

func NewNotificationsCommand(ui terminal.UI, userManager managers.UserManager) (cmd *NotificationsCommand) {
	return &NotificationsCommand{
		UI:          ui,
		UserManager: userManager,
	}
}

func (cmd *NotificationsCommand) Run(c *cli.Context) error {

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	mask := "mask[id, name,description,enabled]"
	notifications, err := cmd.UserManager.GetAllNotifications(mask)
	if err != nil {
		return cli.NewExitError(T("Failed to get notifications.\n")+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, notifications)
	}

	columns := []string{"Id", "Name", "Description", "Enabled"}
	table := cmd.UI.Table(columns)

	for _, notification := range notifications {
		id := utils.FormatIntPointer(notification.Id)
		name := utils.FormatStringPointer(notification.Name)
		description := utils.FormatStringPointer(notification.Description)
		enabled := utils.FormatBoolPointer(notification.Enabled)

		table.Add(id, name, description, enabled)
	}

	table.Print()
	return nil
}

func UserNotificationsMetaData() cli.Command {
	return cli.Command{
		Category:    "user",
		Name:        "notifications",
		Description: T("List email subscription notifications"),
		Usage:       "${COMMAND_NAME} sl user notifications [OPTIONS]",
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

package user

import (
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type NotificationsCommand struct {
	*metadata.SoftlayerCommand
	UserManager managers.UserManager
	Command     *cobra.Command
}

func NewNotificationsCommand(sl *metadata.SoftlayerCommand) (cmd *NotificationsCommand) {
	thisCmd := &NotificationsCommand{
		SoftlayerCommand: sl,
		UserManager:      managers.NewUserManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "notifications",
		Short: T("List email subscription notifications"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *NotificationsCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	mask := "mask[id, name,description,enabled]"
	notifications, err := cmd.UserManager.GetAllNotifications(mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get notifications.\n"), err.Error(), 2)
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

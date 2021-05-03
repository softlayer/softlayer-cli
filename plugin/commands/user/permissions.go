package user

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type PermissionsCommand struct {
	UI          terminal.UI
	UserManager managers.UserManager
}

func NewPermissionsCommand(ui terminal.UI, userManager managers.UserManager) (cmd *PermissionsCommand) {
	return &PermissionsCommand{
		UI:          ui,
		UserManager: userManager,
	}
}

func (cmd *PermissionsCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	userId := c.Args()[0]
	id, err := strconv.Atoi(userId)
	if err != nil {
		return errors.NewInvalidUsageError(T("User ID should be a number."))
	}

	object_mask := "mask[id, permissions, isMasterUserFlag, roles]"
	user, err := cmd.UserManager.GetUser(id, object_mask)
	if err != nil {
		cli.NewExitError(T("Failed to get user.\n")+err.Error(), 2)
	}

	allPermission, err := cmd.UserManager.GetAllPermission()
	if err != nil {
		cli.NewExitError(T("Failed to get permissions.\n")+err.Error(), 2)
	}

	if user.IsMasterUserFlag != nil && *user.IsMasterUserFlag {
		cmd.UI.Print(T("This account is the Master User and has all permissions enabled"))

	}

	table := cmd.UI.Table([]string{T("ID"), T("Role Name"), T("Description")})

	for _, role := range user.Roles {
		roleId := utils.FormatIntPointer(role.Id)
		roleName := utils.FormatStringPointer(role.Name)
		roleDescription := utils.FormatStringPointer(role.Description)

		table.Add(roleId, roleName, roleDescription)
	}
	table.Add("", "", "")
	table.Print()

	tablePermission := cmd.UI.Table([]string{T("Description"), T("KeyName"), T("Assigned")})
	for _, perm := range allPermission {
		var assigned bool
		for _, userPerm := range user.Permissions {
			if perm.KeyName != nil && userPerm.KeyName != nil && *perm.KeyName == *userPerm.KeyName {
				assigned = true
			}
		}
		tablePermission.Add(utils.FormatStringPointer(perm.Name), utils.FormatStringPointer(perm.KeyName), strconv.FormatBool(assigned))
	}
	tablePermission.Print()
	return nil
}

package user

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PermissionsCommand struct {
	*metadata.SoftlayerCommand
	UserManager managers.UserManager
	Command     *cobra.Command
}

func NewPermissionsCommand(sl *metadata.SoftlayerCommand) (cmd *PermissionsCommand) {
	thisCmd := &PermissionsCommand{
		SoftlayerCommand: sl,
		UserManager:      managers.NewUserManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "permissions " + T("USER_ID"),
		Short: T("View user permissions"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *PermissionsCommand) Run(args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidUsageError(T("User ID should be a number."))
	}

	object_mask := "mask[id, permissions, isMasterUserFlag, roles]"
	user, err := cmd.UserManager.GetUser(id, object_mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get user."), err.Error(), 2)
	}

	allPermission, err := cmd.UserManager.GetAllPermission()
	if err != nil {
		return errors.NewAPIError(T("Failed to get permissions."), err.Error(), 2)
	}

	isMasterUser := false
	if user.IsMasterUserFlag != nil && *user.IsMasterUserFlag {
		cmd.UI.Print(T("This account is the Master User and has all permissions enabled"))
		isMasterUser = true
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
		// Display master user as having all permissions, even though they have none, technically.
		if isMasterUser {
			assigned = true
		}
		for _, userPerm := range user.Permissions {
			if perm.KeyName != nil && userPerm.KeyName != nil && *perm.KeyName == *userPerm.KeyName {
				assigned = true
			}

		}
		flag := true
		arr := []string{"ACCOUNT_SUMMARY_VIEW", "REQUEST_COMPLIANCE_REPORT", "COMPANY_EDIT", "ONE_TIME_PAYMENTS", "UPDATE_PAYMENT_DETAILS",
			"EU_LIMITED_PROCESSING_MANAGE", "TICKET_ADD", "TICKET_EDIT", "TICKET_SEARCH", "TICKET_VIEW", "TICKET_VIEW_ALL"}
		for i := 0; i < len(arr); i++ {
			if *perm.KeyName == arr[i] {
				flag = false
			}
		}
		if flag == true {
			tablePermission.Add(utils.FormatStringPointer(perm.Name), utils.FormatStringPointer(perm.KeyName), strconv.FormatBool(assigned))
		}
	}
	tablePermission.Print()
	return nil
}

package user

import (
	"bytes"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
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
	subs := map[string]interface{}{
		"HelpUrl": "https://cloud.ibm.com/docs/account?topic=account-migrated_permissions",
	}
	cobraCmd := &cobra.Command{
		Use:   "permissions " + T("USER_ID"),
		Short: T("View user permissions"),
		Long: T(`Some permissions here may also be managed by the IBM IAM service.
See {{.HelpUrl}} for more details.`, subs),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

type PermissionCollection struct {
	Department string
	Permissions []Permission
}

type Permission struct {
	KeyName string
	Assigned string
	Description string
}

func (cmd *PermissionsCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidUsageError(T("User ID should be a number."))
	}

	object_mask := "mask[id, permissions, isMasterUserFlag, roles]"
	user, err := cmd.UserManager.GetUser(id, object_mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get user."), err.Error(), 2)
	}

	allPermission, err := cmd.UserManager.GetAllPermissionDepartments()
	if err != nil {
		return errors.NewAPIError(T("Failed to get permissions."), err.Error(), 2)
	}

	userPermissions := []PermissionCollection{}

	isMasterUser := false
	if user.IsMasterUserFlag != nil && *user.IsMasterUserFlag {
		if outputFormat != "JSON" {
			cmd.UI.Print(T("This account is the Master User and has all permissions enabled"))
		}
		isMasterUser = true
	}


	for _, department := range allPermission {
		depPerm := PermissionCollection{Department: *department.KeyName}
		for _, perm := range department.Permissions {
			assignedPerm := UserHasPermission(user.Permissions, *perm.KeyName) || isMasterUser
			thisPerm := Permission{
				KeyName:  *perm.KeyName,
				Description: *perm.Description,
				Assigned: strconv.FormatBool(assignedPerm),
			}
			depPerm.Permissions = append(depPerm.Permissions, thisPerm)
		}
		userPermissions = append(userPermissions, depPerm)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, userPermissions)
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

	tablePermission := cmd.UI.Table([]string{T("Department"), T("Permissions")})
	for _, department := range userPermissions {
		buf := new(bytes.Buffer)
		headers := []string{T("KeyName"), T("Assigned"), T("Description")}
		subTable := terminal.NewTable(buf, headers)
		for _, perm := range department.Permissions {
			subTable.Add(perm.KeyName, perm.Assigned, perm.Description)
		}
		subTable.Print()
		tablePermission.Add(department.Department, buf.String())
	}
	tablePermission.Print()
	return nil
}


func UserHasPermission(userPermissions []datatypes.User_Customer_CustomerPermission_Permission, keyName string) bool {
	assigned := false
	for _, userPerm := range userPermissions {
		if *userPerm.KeyName == keyName {
			assigned = true
		}
	}
	return assigned
}
package user

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

var (
	NS_USER_NAME  = "user"
	CMD_USER_NAME = "user"

	CMD_USER_CREATE_NAME           = "create"
	CMD_USER_DELETE_NAME           = "delete"
	CMD_USER_DETAIL_NAME           = "detail"
	CMD_USER_EDIT_DETAILS_NAME     = "detail-edit"
	CMD_USER_EDIT_PERMISSIONS_NAME = "permission-edit"
	CMD_USER_LIST_NAME             = "list"
	CMD_USER_PERMISSIONS_NAME      = "permissions"
)

func GetCommandActionBindings(ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	userManager := managers.NewUserManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{

		CMD_USER_NAME + "-" + CMD_USER_CREATE_NAME: func(c *cli.Context) error {
			return NewCreateCommand(ui, userManager).Run(c)
		},
		CMD_USER_NAME + "-" + CMD_USER_LIST_NAME: func(c *cli.Context) error {
			return NewListCommand(ui, userManager).Run(c)
		},
		CMD_USER_NAME + "-" + CMD_USER_DELETE_NAME: func(c *cli.Context) error {
			return NewDeleteCommand(ui, userManager).Run(c)
		},
		CMD_USER_NAME + "-" + CMD_USER_DETAIL_NAME: func(c *cli.Context) error {
			return NewDetailsCommand(ui, userManager).Run(c)
		},
		CMD_USER_NAME + "-" + CMD_USER_PERMISSIONS_NAME: func(c *cli.Context) error {
			return NewPermissionsCommand(ui, userManager).Run(c)
		},
		CMD_USER_NAME + "-" + CMD_USER_EDIT_DETAILS_NAME: func(c *cli.Context) error {
			return NewEditCommand(ui, userManager).Run(c)
		},
		CMD_USER_NAME + "-" + CMD_USER_EDIT_PERMISSIONS_NAME: func(c *cli.Context) error {
			return NewEditPermissionCommand(ui, userManager).Run(c)
		},
	}
	return CommandActionBindings

}

func UserNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        NS_USER_NAME,
		Description: T("Classic infrastructure Manage Users"),
	}
}

func UserMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        CMD_USER_NAME,
		Usage:       "${COMMAND_NAME} sl user",
		Description: T("Classic infrastructure Manage Users"),
		Subcommands: []cli.Command{
			UserCreateMetaData(),
			UserDeleteMataData(),
			UserDetailMetaData(),
			UserEditMetaData(),
			UserEditPermissionMetaData(),
			UserListMetaData(),
			UserPermissionsMetaData(),
		},
	}
}

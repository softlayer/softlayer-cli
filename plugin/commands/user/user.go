package user

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)


func GetCommandActionBindings(ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	userManager := managers.NewUserManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{

		"user-create": func(c *cli.Context) error {
			return NewCreateCommand(ui, userManager).Run(c)
		},
		"user-list": func(c *cli.Context) error {
			return NewListCommand(ui, userManager).Run(c)
		},
		"user-delete": func(c *cli.Context) error {
			return NewDeleteCommand(ui, userManager).Run(c)
		},
		"user-detail": func(c *cli.Context) error {
			return NewDetailsCommand(ui, userManager).Run(c)
		},
		"user-permissions": func(c *cli.Context) error {
			return NewPermissionsCommand(ui, userManager).Run(c)
		},
		"user-detail-edit": func(c *cli.Context) error {
			return NewEditCommand(ui, userManager).Run(c)
		},
		"user-permission-edit": func(c *cli.Context) error {
			return NewEditPermissionCommand(ui, userManager).Run(c)
		},
	}
	return CommandActionBindings

}

func UserNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "user",
		Description: T("Classic infrastructure Manage Users"),
	}
}

func UserMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "user",
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

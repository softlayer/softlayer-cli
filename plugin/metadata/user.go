package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
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

func UserNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_USER_NAME,
		Description: T("Classic infrastructure Manage Users"),
	}
}

func UserMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
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

func UserCreateMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_USER_NAME,
		Name:        CMD_USER_CREATE_NAME,
		Description: T("Creates a user"),
		Usage: T(`${COMMAND_NAME} sl user create USERNAME [OPTIONS] 

EXAMPLE: 	
    ${COMMAND_NAME} sl user create my@email.com --email my@email.com --password generate --api-key --template '{"firstName": "Test", "lastName": "Testerson"}'
    Remember to set the permissions and access for this new user.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "email",
				Usage: T("Email address for this user. Required for creation"),
			},
			cli.StringFlag{
				Name:  "password",
				Usage: T("Password to set for this user. If no password is provided, the user is sent an email to generate one, which expires in 24 hours. Specify the '-p generate' option to generate a password for you. Passwords require 8+ characters, uppercase and lowercase, a number and a symbol"),
			},
			cli.IntFlag{
				Name:  "from-user",
				Usage: T("Base user to use as a template for creating this user. The default is to use the user that is running this command. Information provided in --template supersedes this template"),
			},
			cli.StringFlag{
				Name:  "template",
				Usage: T("A json string describing https://softlayer.github.io/reference/datatypes/SoftLayer_User_Customer/"),
			},
			cli.BoolFlag{
				Name:  "api-key",
				Usage: T("Create an API key for this user"),
			},
			cli.StringFlag{
				Name:  "vpn-password",
				Usage: T("VPN password to set for this user."),
			},
			ForceFlag(),
		},
	}
}

func UserDeleteMataData() cli.Command {
	return cli.Command{
		Category:    CMD_USER_NAME,
		Name:        CMD_USER_DELETE_NAME,
		Description: T("Sets a user's status to CANCEL_PENDING, which will immediately disable the account, and will eventually be fully removed from the account by an automated internal process"),
		Usage: T(`${COMMAND_NAME} sl user delete IDENTIFIER [OPTIONS]
	
EXAMPLE: 
   ${COMMAND_NAME} sl user delete userId
   This command delete user with userId.`),
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}

func UserDetailMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_USER_NAME,
		Name:        CMD_USER_DETAIL_NAME,
		Description: T("User details"),
		Usage:       "${COMMAND_NAME} sl user detail IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "keys",
				Usage: T("Show the users API key"),
			},
			cli.BoolFlag{
				Name:  "permissions",
				Usage: T("Display permissions assigned to this user. Master users do not show permissions"),
			},
			cli.BoolFlag{
				Name:  "hardware",
				Usage: T("Display hardware this user has access to"),
			},
			cli.BoolFlag{
				Name:  "virtual",
				Usage: T("Display virtual guests this user has access to"),
			},
			cli.BoolFlag{
				Name:  "logins",
				Usage: T("Show login history of this user for the last 24 hours"),
			},
			cli.BoolFlag{
				Name:  "events",
				Usage: T("Show audit log for this user"),
			},
		},
	}
}

func UserEditMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_USER_NAME,
		Name:        CMD_USER_EDIT_DETAILS_NAME,
		Description: T("Edit a user's details"),
		Usage: T(`${COMMAND_NAME} sl user detail-edit IDENTIFIER [OPTIONS]

EXAMPLE: 
    ${COMMAND_NAME} sl user detail-edit USER_ID --template '{"firstName": "Test", "lastName": "Testerson"}'
    This command edit a users details.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "template",
				Usage: T("A json string describing https://softlayer.github.io/reference/datatypes/SoftLayer_User_Customer/"),
			},
			OutputFlag(),
		},
	}
}

func UserEditPermissionMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_USER_NAME,
		Name:        CMD_USER_EDIT_PERMISSIONS_NAME,
		Description: T("Enable or Disable specific permissions"),
		Usage:       "${COMMAND_NAME} sl user permission-edit IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "enable",
				Usage: T("Enable or Disable selected permissions. Accepted inputs are 'true' and 'false'. default is 'true'"),
			},
			cli.StringSliceFlag{
				Name:  "permission",
				Usage: T("Permission keyName to set. Use keyword ALL to select ALL permissions"),
			},
			cli.IntFlag{
				Name:  "from-user",
				Usage: T("Set permissions to match this user's permissions. Adds and removes the appropriate permissions"),
			},
		},
	}
}

func UserListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_USER_NAME,
		Name:        CMD_USER_LIST_NAME,
		Description: T("List Users"),
		Usage:       "${COMMAND_NAME} sl user list [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "column",
				Usage: T("Column to display. options are: id,username,email,displayName,status,hardwareCount,virtualGuestCount. This option can be specified multiple times"),
			},
			cli.StringSliceFlag{
				Name:   "columns",
				Hidden: true,
			},
			OutputFlag(),
		},
	}
}

func UserPermissionsMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_USER_NAME,
		Name:        CMD_USER_PERMISSIONS_NAME,
		Description: T("View user permissions"),
		Usage:       "${COMMAND_NAME} sl user permissions IDENTIFIER",
		Flags:       []cli.Flag{},
	}
}

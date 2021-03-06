package user

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type EditPermissionCommand struct {
	UI          terminal.UI
	UserManager managers.UserManager
}

func NewEditPermissionCommand(ui terminal.UI, userManager managers.UserManager) (cmd *EditPermissionCommand) {
	return &EditPermissionCommand{
		UI:          ui,
		UserManager: userManager,
	}
}

func (cmd *EditPermissionCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	userId := c.Args()[0]
	id, err := strconv.Atoi(userId)
	if err != nil {
		return errors.NewInvalidUsageError(T("User ID should be a number."))
	}

	if (c.IsSet("from-user") && c.IsSet("permission")) || (!c.IsSet("from-user") && !c.IsSet("permission")) {
		return errors.NewInvalidUsageError(T("one of --permission and --from-user should be used to specify permissions"))
	}

	permissionKeynames := c.StringSlice("permission")
	permissions, err := cmd.UserManager.FormatPermissionObject(permissionKeynames)
	if err != nil {
		return err
	}

	enableFlag := true
	enable := c.String("enable")
	if enable != "" {
		enable = strings.ToLower(enable)
		if enable != "true" && enable != "false" {
			return cli.NewExitError(fmt.Sprintf(T("options for %s are true, false"), "enable"), 1)
		}
		enableFlag = (enable == "true")
	}

	if c.IsSet("from-user") {
		fromUser := c.Int("from-user")
		err = cmd.UserManager.PermissionFromUser(id, fromUser)
	} else if enableFlag {
		_, err = cmd.UserManager.AddPermission(id, permissions)
	} else {
		_, err = cmd.UserManager.RemovePermission(id, permissions)
	}

	if err != nil {
		return cli.NewExitError(fmt.Sprintf(T("Failed to update permissions: %s"), err.Error()), 1)
	}
	cmd.UI.Print(fmt.Sprintf(T("Permissions updated successfully: %s"), strings.Join(permissionKeynames, ",")))
	return nil
}

func UserEditPermissionMetaData() cli.Command {
	return cli.Command{
		Category:    "user",
		Name:        "permission-edit",
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

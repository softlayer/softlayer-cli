package user

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

const (
	NO_ZERO_VALUE = "yes"
)

type ListCommand struct {
	UI          terminal.UI
	UserManager managers.UserManager
}

func NewListCommand(ui terminal.UI, userManager managers.UserManager) (cmd *ListCommand) {
	return &ListCommand{
		UI:          ui,
		UserManager: userManager,
	}
}

var maskMap = map[string]string{
	"id":                "id",
	"username":          "username",
	"email":             "email",
	"displayName":       "displayName",
	"status":            "userStatus.name",
	"hardwareCount":     "hardwareCount",
	"virtualGuestCount": "virtualGuestCount",
	"2FA":               "externalBindingCount",
	"classicAPIKey":     "apiAuthenticationKeyCount",
}

func (cmd *ListCommand) Run(c *cli.Context) error {

	var columns []string
	if c.IsSet("column") {
		columns = c.StringSlice("column")
	} else if c.IsSet("columns") {
		columns = c.StringSlice("columns")
	}

	defaultColumns := []string{"id", "username", "email", "displayName", "2FA", "classicAPIKey"}
	optionalColumns := []string{"status", "hardwareCount", "virtualGuestCount"}

	showColumns, err := utils.ValidateColumns("", columns, defaultColumns, optionalColumns, []string{}, c)
	if err != nil {
		return err
	}

	mask := utils.GetMask(maskMap, showColumns, "")

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	users, err := cmd.UserManager.ListUsers(mask)
	if err != nil {
		return cli.NewExitError(T("Failed to list users.\n")+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, users)
	}

	table := cmd.UI.Table(showColumns)

	for _, user := range users {
		values := make(map[string]string)
		values["id"] = utils.FormatIntPointer(user.Id)
		values["username"] = utils.FormatStringPointer(user.Username)
		values["email"] = utils.FormatStringPointer(user.Email)
		values["displayName"] = utils.FormatStringPointer(user.DisplayName)

		values["2FA"] = utils.ReplaceUIntPointerValue(user.ExternalBindingCount, NO_ZERO_VALUE)
		values["classicAPIKey"] = utils.ReplaceUIntPointerValue(user.ApiAuthenticationKeyCount, NO_ZERO_VALUE)

		if user.UserStatus != nil {
			values["status"] = utils.FormatStringPointer(user.UserStatus.Name)
		} else {
			values["status"] = utils.EMPTY_VALUE
		}

		values["hardwareCount"] = utils.FormatUIntPointer(user.HardwareCount)
		values["virtualGuestCount"] = utils.FormatUIntPointer(user.VirtualGuestCount)

		row := make([]string, len(showColumns))
		for i, col := range showColumns {
			row[i] = values[col]
		}
		table.Add(row...)
	}
	table.Print()
	return nil
}

func UserListMetaData() cli.Command {
	return cli.Command{
		Category:    "user",
		Name:        "list",
		Description: T("List Users"),
		Usage:       "${COMMAND_NAME} sl user list [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "column",
				Usage: T("Column to display. options are: id,username,email,displayName,2FA,classicAPIKey,status,hardwareCount,virtualGuestCount. This option can be specified multiple times"),
			},
			cli.StringSliceFlag{
				Name:   "columns",
				Hidden: true,
			},
			metadata.OutputFlag(),
		},
	}
}

package user

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
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
}

func (cmd *ListCommand) Run(c *cli.Context) error {

	var columns []string
	if c.IsSet("column") {
		columns = c.StringSlice("column")
	} else if c.IsSet("columns") {
		columns = c.StringSlice("columns")
	}

	defaultColumns := []string{"id", "username", "email", "displayName"}
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
		cli.NewExitError(T("Failed to list users.\n")+err.Error(), 2)
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

		if user.UserStatus != nil {
			values["status"] = utils.FormatStringPointer(user.UserStatus.Name)
		} else {
			values["status"] = "-"
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

package user

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type EditCommand struct {
	UI          terminal.UI
	UserManager managers.UserManager
}

func NewEditCommand(ui terminal.UI, userManager managers.UserManager) (cmd *EditCommand) {
	return &EditCommand{
		UI:          ui,
		UserManager: userManager,
	}
}

func (cmd *EditCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	userId := c.Args()[0]
	id, err := strconv.Atoi(userId)
	if err != nil {
		return errors.NewInvalidUsageError(T("User ID should be a number."))
	}

	if !c.IsSet("template") {
		return errors.NewMissingInputError("--template")
	}

	var templateStruct datatypes.User_Customer
	template := c.String("template")
	err = json.Unmarshal([]byte(template), &templateStruct)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf(T("Unable to unmarshal template json: %s\n"), err.Error()), 1)
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	resp, err := cmd.UserManager.EditUser(templateStruct, id)
	if err != nil {
		return cli.NewExitError(T("Failed to update user {{.UserID}}.", map[string]interface{}{"UserID": id}), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Print(T("User {{.UserID}} updated successfully.", map[string]interface{}{"UserID": id}))
	return nil
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
			metadata.OutputFlag(),
		},
	}
}
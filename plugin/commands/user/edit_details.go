package user

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type EditCommand struct {
	*metadata.SoftlayerCommand
	UserManager managers.UserManager
	Command     *cobra.Command
	Template    string
}

func NewEditCommand(sl *metadata.SoftlayerCommand) (cmd *EditCommand) {
	thisCmd := &EditCommand{
		SoftlayerCommand: sl,
		UserManager:      managers.NewUserManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "detail-edit " + T("USER_ID"),
		Short: T("Edit a user's details"),
		Long: T(`
EXAMPLE: 
	${COMMAND_NAME} sl user detail-edit USER_ID --template '{"firstName": "Test", "lastName": "Testerson"}'
	This command edit a users details.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Template, "template", " ", T("A json string describing https://softlayer.github.io/reference/datatypes/SoftLayer_User_Customer/"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *EditCommand) Run(args []string) error {
	userId := args[0]
	id, err := strconv.Atoi(userId)
	if err != nil {
		return errors.NewInvalidUsageError(T("User ID should be a number."))
	}

	if cmd.Template == " " {
		return errors.NewMissingInputError("--template")
	}

	var templateStruct datatypes.User_Customer
	template := cmd.Template
	err = json.Unmarshal([]byte(template), &templateStruct)
	if err != nil {
		return errors.NewInvalidUsageError(fmt.Sprintf(T("Unable to unmarshal template json: %s\n"), err.Error()))
	}

	outputFormat := cmd.GetOutputFlag()

	resp, err := cmd.UserManager.EditUser(templateStruct, id)
	if err != nil {
		return errors.NewAPIError(T("Failed to update user {{.UserID}}.", map[string]interface{}{"UserID": id}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Print(T("User {{.UserID}} updated successfully.", map[string]interface{}{"UserID": id}))
	return nil
}

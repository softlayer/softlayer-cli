package email

import (
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type EditCommand struct {
	*metadata.SoftlayerCommand
	EmailManager managers.EmailManager
	Command      *cobra.Command
	Username     string
	Email        string
	Password     string
}

func NewEditCommand(sl *metadata.SoftlayerCommand) (cmd *EditCommand) {
	thisCmd := &EditCommand{
		SoftlayerCommand: sl,
		EmailManager:     managers.NewEmailManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "edit " + T("IDENTIFIER"),
		Short: T("Edit details of an email delivery account."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Username, "username", "", T("Sets username for this account."))
	cobraCmd.Flags().StringVar(&thisCmd.Email, "email", "", T("Sets the contact email for this account."))
	cobraCmd.Flags().StringVar(&thisCmd.Password, "password", "", T("Password must be between 8 and 20 characters and must contain one letter and one number."))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *EditCommand) Run(args []string) error {

	emailID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError(T("email ID"))
	}

	username := cmd.Username
	email := cmd.Email
	password := cmd.Password

	if username == "" && email == "" && password == "" {
		return slErr.NewInvalidUsageError(T("Please pass at least one of the flags."))
	}

	if email != "" {
		err = cmd.EmailManager.UpdateEmail(emailID, email)
		if err != nil {
			return slErr.NewAPIError(T("Failed to Edit emailAddress account: {{.emailID}}.\n", map[string]interface{}{"emailID": emailID}), err.Error(), 2)
		}
		cmd.UI.Ok()
		cmd.UI.Print(T("Email address {{.emailID}} was updated.", map[string]interface{}{"emailID": emailID}))
	}

	if username != "" || password != "" {
		templateObject := datatypes.Network_Message_Delivery{}
		if username != "" {
			templateObject.Username = sl.String(username)
		}
		if password != "" {
			templateObject.Password = sl.String(password)
		}
		err = cmd.EmailManager.EditObject(emailID, templateObject)
		if err != nil {
			return slErr.NewAPIError(T("Failed to Edit email account: {{.emailID}}.\n", map[string]interface{}{"emailID": emailID}), err.Error(), 2)
		}
		cmd.UI.Ok()
		cmd.UI.Print(T("Email account {{.emailID}} was updated.", map[string]interface{}{"emailID": emailID}))
	}

	return nil
}

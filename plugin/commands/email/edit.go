package email

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type EditCommand struct {
	UI           terminal.UI
	EmailManager managers.EmailManager
}

func NewEditCommand(ui terminal.UI, emailManager managers.EmailManager) (cmd *EditCommand) {
	return &EditCommand{
		UI:           ui,
		EmailManager: emailManager,
	}
}

func EditMetaData() cli.Command {
	return cli.Command{
		Category:    "email",
		Name:        "edit",
		Description: T("Edit details of an email delivery account."),
		Usage:       T(`${COMMAND_NAME} sl email edit`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "username",
				Usage: T("Sets username for this account."),
			},
			cli.StringFlag{
				Name:  "email",
				Usage: T("Sets the contact email for this account."),
			},
			cli.StringFlag{
				Name:  "password",
				Usage: T("Password must be between 8 and 20 characters and must contain one letter and one number."),
			},
		},
	}
}

func (cmd *EditCommand) Run(c *cli.Context) error {

	if c.NArg() != 1 {
		return slErr.NewInvalidUsageError(T("This command requires one argument."))
	}

	emailID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError(T("email ID"))
	}

	if !c.IsSet("username") && !c.IsSet("email") && !c.IsSet("password") {
		return slErr.NewInvalidUsageError(T("Please pass at least one of the flags."))
	}

	if c.IsSet("email") {
		err = cmd.EmailManager.UpdateEmail(emailID, c.String("email"))
		if err != nil {
			return cli.NewExitError(T("Failed to Edit emailAddress account: {{.emailID}}.\n", map[string]interface{}{"emailID": emailID})+err.Error(), 2)
		}
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Email address {{.emailID}} was updated.", map[string]interface{}{"emailID": emailID}))

	if c.IsSet("username") || c.IsSet("password") {
		err = cmd.EmailManager.EditObject(emailID, c.String("username"), c.String("password"))
		if err != nil {
			return cli.NewExitError(T("Failed to Edit email account: {{.emailID}}.\n", map[string]interface{}{"emailID": emailID})+err.Error(), 2)
		}
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Email account {{.emailID}} was updated.", map[string]interface{}{"emailID": emailID}))
	return nil
}

package ticket

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type UpdateTicketCommand struct {
	UI            terminal.UI
	TicketManager managers.TicketManager
}

func NewUpdateTicketCommand(ui terminal.UI, ticketManager managers.TicketManager) (cmd *UpdateTicketCommand) {
	return &UpdateTicketCommand{
		UI:            ui,
		TicketManager: ticketManager,
	}
}

func (cmd *UpdateTicketCommand) Run(c *cli.Context) error {
	nargs := c.NArg()
	if nargs < 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	args := c.Args()

	ticketid, err := strconv.Atoi(args[0])
	if err != nil || ticketid <= 0 {
		return errors.NewInvalidUsageError(T("The ticket id must be a positive non-zero number."))
	}

	content := ""

	if nargs == 1 {
		content, err = cmd.TicketManager.GetText()
		if err != nil {
			return err
		}
	} else {
		content = args[1]
	}

	err = cmd.TicketManager.AddUpdate(ticketid, content)
	if err != nil {
		return cli.NewExitError(T("Update could not be added: {{.Error}}\n", map[string]interface{}{"Error": err.Error()}), 2)
	}
	cmd.UI.Ok()
	return nil
}

func TicketUpdataMetaData() cli.Command {
	return cli.Command{
		Category:    "ticket",
		Name:        "update",
		Description: T("Adds an update to an existing ticket"),
		Usage: T(`${COMMAND_NAME} sl ticket update TICKETID ["CONTENTS"] 
  
    If the second argument is not specified on a non-Windows machine, it will attempt to use either the value stored in the EDITOR environmental variable, or find either nano, vim, or emacs in that order.
  
EXAMPLE:
  ${COMMAND_NAME} sl ticket update 767676 "A problem has been detected."
  ${COMMAND_NAME} sl ticket update 767667`),
	}
}

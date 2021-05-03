package ticket

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
)

type UploadFileTicketCommand struct {
	UI            terminal.UI
	TicketManager managers.TicketManager
}

func NewUploadFileTicketCommand(ui terminal.UI, ticketManager managers.TicketManager) (cmd *UploadFileTicketCommand) {
	return &UploadFileTicketCommand{
		UI:            ui,
		TicketManager: ticketManager,
	}
}

func (cmd *UploadFileTicketCommand) Run(c *cli.Context) error {
	if c.NArg() != 2 {
		return errors.NewInvalidUsageError(T("This command requires two arguments."))
	}

	args := c.Args()

	ticketId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidUsageError(T("The ticket id must be a number."))
	}

	file_path := args[1]
	name := c.String("name")

	err = cmd.TicketManager.AttachFileToTicket(ticketId, name, file_path)

	if err != nil {
		return cli.NewExitError(T("Error: {{.Error}}", map[string]interface{}{"Error": err.Error()}), 2)
	} else {
		cmd.UI.Ok()
		return nil
	}

}

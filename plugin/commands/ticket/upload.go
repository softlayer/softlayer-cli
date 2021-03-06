package ticket

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
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

func TicketUploadMetaData() cli.Command {
	return cli.Command{
		Category:    "ticket",
		Name:        "upload",
		Description: T("Adds an attachment to an existing ticket"),
		Usage: T(`${COMMAND_NAME} sl ticket upload TICKETID FILEPATH
  
EXAMPLE:
	${COMMAND_NAME} sl ticket upload 767676 "/home/user/screenshot.png"`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name",
				Usage: T("The name of the attachment shown in the ticket"),
			},
		},
	}
}

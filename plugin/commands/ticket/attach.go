package ticket

import (
	"strconv"

	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
)

type AttachDeviceTicketCommand struct {
	UI            terminal.UI
	TicketManager managers.TicketManager
}

func NewAttachDeviceTicketCommand(ui terminal.UI, ticketManager managers.TicketManager) (cmd *AttachDeviceTicketCommand) {
	return &AttachDeviceTicketCommand{
		UI:            ui,
		TicketManager: ticketManager,
	}
}

func (cmd *AttachDeviceTicketCommand) Run(c *cli.Context) error {

	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	args := c.Args()

	ticketid, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidUsageError(T("The ticket id must be a number."))
	}

	ishardware := c.IsSet("hardware")
	isvirtual := c.IsSet("virtual")
	if ishardware && isvirtual {
		return errors.NewInvalidUsageError(T("hardware and virtual flags cannot be set at the same time."))
	} else if !ishardware && !isvirtual {
		return errors.NewInvalidUsageError(T("either the hardware or virtual flag must be set."))
	}

	if ishardware {
		deviceid := c.Int("hardware")
		err = cmd.TicketManager.AttachDeviceToTicket(ticketid, deviceid, true)
	} else if isvirtual {
		deviceid := c.Int("virtual")
		err = cmd.TicketManager.AttachDeviceToTicket(ticketid, deviceid, false)
	}

	if err != nil {
		return cli.NewExitError(T("Error: {{.Error}}", map[string]interface{}{"Error": err.Error()}), 2)
	} else {
		cmd.UI.Ok()
		return nil
	}

}

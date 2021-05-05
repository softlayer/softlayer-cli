package ticket

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type DetachDeviceTicketCommand struct {
	UI            terminal.UI
	TicketManager managers.TicketManager
}

func NewDetachDeviceTicketCommand(ui terminal.UI, ticketManager managers.TicketManager) (cmd *DetachDeviceTicketCommand) {
	return &DetachDeviceTicketCommand{
		UI:            ui,
		TicketManager: ticketManager,
	}
}

func (cmd *DetachDeviceTicketCommand) Run(c *cli.Context) error {
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
		err = cmd.TicketManager.RemoveDeviceFromTicket(ticketid, deviceid, true)
	} else if isvirtual {
		deviceid := c.Int("virtual")
		err = cmd.TicketManager.RemoveDeviceFromTicket(ticketid, deviceid, false)
	}

	if err != nil {
		return cli.NewExitError(T("Error: {{.Error}}", map[string]interface{}{"Error": err.Error()}), 2)
	} else {
		cmd.UI.Ok()
		return nil
	}

}

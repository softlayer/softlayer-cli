package ticket

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type ListTicketCommand struct {
	UI            terminal.UI
	TicketManager managers.TicketManager
}

func NewListTicketCommand(ui terminal.UI, ticketManager managers.TicketManager) (cmd *ListTicketCommand) {
	return &ListTicketCommand{
		UI:            ui,
		TicketManager: ticketManager,
	}
}

func (cmd *ListTicketCommand) Run(c *cli.Context) error {
	if c.NArg() != 0 {
		return errors.NewInvalidUsageError(T("This command requires zero arguments."))
	}

	var err error
	var tickets, ticketsOpen, ticketsClose []datatypes.Ticket

	if c.Bool("open") && c.Bool("closed") {
		ticketsOpen, err = cmd.TicketManager.ListOpenTickets()
		ticketsClose, err = cmd.TicketManager.ListCloseTickets()
		tickets = append(ticketsOpen, ticketsClose...)
	} else if !c.Bool("open") && c.Bool("closed") {
		tickets, err = cmd.TicketManager.ListCloseTickets()
	} else {
		tickets, err = cmd.TicketManager.ListOpenTickets()
	}

	if err != nil {
		return cli.NewExitError(T("Error: {{.Error}}", map[string]interface{}{"Error": err.Error()}), 2)
	}

	columns := []string{T("Id"), T("Assigned User"), T("Title"), T("Last Edited"), T("Status"), T("Updates"), T("Priority")}

	table := cmd.UI.Table(utils.GetColumnHeader(columns))

	for _, ticket := range tickets {
		row := make([]string, len(columns))
		values := make(map[string]string, len(columns))

		user := "-"
		if ticket.AssignedUser != nil {
			user = utils.FormatStringPointer(ticket.AssignedUser.FirstName) + " " + utils.FormatStringPointer(ticket.AssignedUser.LastName)
		}
		if ticket.Id != nil {
			values["Id"] = strconv.Itoa(*ticket.Id)
		}
		values["Assigned User"] = user
		values["Title"] = utils.FormatStringPointer(ticket.Title)
		values["Last Edited"] = utils.FormatSLTimePointer(ticket.LastEditDate)
		values["Status"] = utils.FormatStringPointer(ticket.Status.Name)
		values["Updates"] = utils.FormatUIntPointer(ticket.UpdateCount)
		values["Priority"] = utils.FormatIntPointer(ticket.Priority)
		for i, col := range columns {
			row[i] = values[col]
		}

		table.Add(row...)
	}

	table.Print()
	return nil

}

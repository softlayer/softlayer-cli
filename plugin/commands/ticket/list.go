package ticket

import (
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListTicketCommand struct {
	*metadata.SoftlayerCommand
	TicketManager managers.TicketManager
	Command       *cobra.Command
	Open          bool
	Closed        bool
}

func NewListTicketCommand(sl *metadata.SoftlayerCommand) *ListTicketCommand {
	thisCmd := &ListTicketCommand{
		SoftlayerCommand: sl,
		TicketManager:    managers.NewTicketManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List tickets"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVar(&thisCmd.Open, "open", false, T("Display only open tickets"))
	cobraCmd.Flags().BoolVar(&thisCmd.Closed, "closed", false, T("Display only closed tickets"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ListTicketCommand) Run(args []string) error {
	var err error
	var tickets, ticketsOpen, ticketsClose []datatypes.Ticket

	if cmd.Open && cmd.Closed {
		ticketsOpen, err = cmd.TicketManager.ListOpenTickets()
		ticketsClose, err = cmd.TicketManager.ListCloseTickets()
		tickets = append(ticketsOpen, ticketsClose...)
	} else if !cmd.Open && cmd.Closed {
		tickets, err = cmd.TicketManager.ListCloseTickets()
	} else {
		tickets, err = cmd.TicketManager.ListOpenTickets()
	}

	if err != nil {
		return errors.New(T("Error: {{.Error}}", map[string]interface{}{"Error": err.Error()}))
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

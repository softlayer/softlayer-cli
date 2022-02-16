package ticket

import (
	"bytes"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailTicketCommand struct {
	UI            terminal.UI
	TicketManager managers.TicketManager
	UserManager   managers.UserManager
}

func NewDetailTicketCommand(ui terminal.UI, ticketManager managers.TicketManager, userManager managers.UserManager) (cmd *DetailTicketCommand) {
	return &DetailTicketCommand{
		UI:            ui,
		TicketManager: ticketManager,
		UserManager:   userManager,
	}
}

func (cmd *DetailTicketCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	args := c.Args()

	ticketid, err := strconv.Atoi(args[0])
	if err != nil || ticketid <= 0 {
		return errors.NewInvalidUsageError(T("The ticket id must be a positive non-zero number."))
	}

	ticket, err := cmd.TicketManager.GetTicket(ticketid)
	if err != nil {
		return cli.NewExitError(T("Error: {{.Error}}", map[string]interface{}{"Error": err.Error()}), 2)
	}
	ticketUpdates, err := cmd.TicketManager.GetAllUpdates(ticketid)
	if err != nil {
		return cli.NewExitError(T("Error: {{.Error}}", map[string]interface{}{"Error": err.Error()}), 2)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(ticket.Id))
	table.Add(T("Title"), utils.FormatStringPointer(ticket.Title))
	table.Add(T("Priority"), utils.FormatIntPointer(ticket.Priority))

	user := ticket.AssignedUser
	if user != nil {
		table.Add(T("User"), utils.FormatStringPointer(user.FirstName)+" "+utils.FormatStringPointer(user.LastName))
	}

	if ticket.Status != nil {

		table.Add(T("Status"), utils.FormatStringPointer(ticket.Status.Name))
	}
	table.Add(T("Created"), utils.FormatSLTimePointer(ticket.CreateDate))
	table.Add(T("Last Edited"), utils.FormatSLTimePointer(ticket.LastEditDate))

	updateCount := 10
	if c.IsSet("count") {
		updateCount = c.Int("count")
	}

	count := Min(len(ticketUpdates), updateCount)
	updates := ticketUpdates[len(ticketUpdates)-count:]

	num := len(ticketUpdates) - count
	columnsList := []string{T("Editor"), T("Create Date"), T("Update ID")}
	for _, update := range updates {
		buf := new(bytes.Buffer)
		values := make(map[string]string, len(columnsList))
		editor_type := utils.FormatStringPointer(update.EditorType)
		var editor_name string
		if editor_type == "USER" && update.EditorId != nil {
			user, err := cmd.UserManager.GetUser(*update.EditorId, "mask[firstName,lastName]")
			if err != nil {
				return cli.NewExitError(T("Error: {{.Error}}.\n", map[string]interface{}{"Error": err.Error()}), 2)
			} else {
				editor_name = utils.FormatStringPointer(user.FirstName) + " " + utils.FormatStringPointer(user.LastName)
			}
		} else {
			editor_name = "Employee"
		}

		values["Editor"] = editor_name
		values["Create Date"] = utils.FormatSLTimePointer(update.CreateDate)
		values["Update ID"] = utils.FormatIntPointer(update.Id)

		row := make([]string, len(columnsList))
		for i, col := range columnsList {
			row[i] = values[col]
		}
		tableUpdate := terminal.NewTable(buf, columnsList)
		tableUpdate.Add(row...)
		tableUpdate.Print()
		table.Add("Update "+strconv.Itoa(num+1), buf.String())
		if update.Entry != nil {
			table.Add("", *update.Entry)
		}
		table.Add("", "")
		num = num + 1
	}

	table.Print()

	return nil
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func TicketDetailMetaData() cli.Command {
	return cli.Command{
		Category:    "ticket",
		Name:        "detail",
		Description: T("Get details for a ticket"),
		Usage: T(`${COMMAND_NAME} sl ticket detail TICKETID [OPTIONS]
  
EXAMPLE:
  ${COMMAND_NAME} sl ticket detail 767676
  ${COMMAND_NAME} sl ticket detail 767676 --count 10`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "count",
				Usage: T("Number of updates"),
			},
		},
	}
}

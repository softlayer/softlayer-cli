package ticket

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
)

type SummaryTicketCommand struct {
	UI            terminal.UI
	TicketManager managers.TicketManager
	UserManager   managers.UserManager
}

func NewSummaryTicketCommand(ui terminal.UI, ticketManager managers.TicketManager) (cmd *SummaryTicketCommand) {
	return &SummaryTicketCommand{
		UI:            ui,
		TicketManager: ticketManager,
	}
}

func (cmd *SummaryTicketCommand) Run(c *cli.Context) error {
	summary, err := cmd.TicketManager.Summary()

	if err != nil {
		return cli.NewExitError(T("Error: {{.Error}}.\n", map[string]interface{}{"Error": err.Error()}), 2)
	} else {
		table := cmd.UI.Table([]string{T("Status:"), T("Count")})

		table.Add(T("Open:"), "")
		table.Add(T("Accounting"), strconv.Itoa(int(summary.Accounting)))
		table.Add(T("Billing"), strconv.Itoa(int(summary.Billing)))
		table.Add(T("Sales"), strconv.Itoa(int(summary.Sales)))
		table.Add(T("Support"), strconv.Itoa(int(summary.Support)))
		table.Add(T("Other"), strconv.Itoa(int(summary.Other)))
		table.Add(T("Total"), strconv.Itoa(int(summary.Open)))
		table.Add("", "")
		table.Add(T("Closed:"), "")
		table.Add(T("Total"), strconv.Itoa(int(summary.Closed)))

		table.Print()

		return nil
	}

}

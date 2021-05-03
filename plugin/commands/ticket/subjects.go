package ticket

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type SubjectsTicketCommand struct {
	UI            terminal.UI
	TicketManager managers.TicketManager
}

func NewSubjectsTicketCommand(ui terminal.UI, ticketManager managers.TicketManager) (cmd *SubjectsTicketCommand) {
	return &SubjectsTicketCommand{
		UI:            ui,
		TicketManager: ticketManager,
	}
}

func (cmd *SubjectsTicketCommand) Run(c *cli.Context) error {
	if c.NArg() != 0 {
		return errors.NewInvalidUsageError(T("This command uses 0 arguments"))
	}

	subjects, err := cmd.TicketManager.GetSubjects()

	if err != nil {
		return cli.NewExitError(T("Error: {{.Error}}", map[string]interface{}{"Error": err.Error()}), 2)
	} else {
		columnsList := []string{T("ID"), T("Subject")}
		table := cmd.UI.Table(utils.GetColumnHeader(columnsList))

		for _, sub := range *subjects {
			row := make([]string, len(columnsList))
			values := make(map[string]string, len(columnsList))
			if sub.Id != nil {
				values["ID"] = strconv.Itoa(*sub.Id)
			}
			if sub.Name != nil {
				values["Subject"] = *sub.Name
			}
			for i, col := range columnsList {
				row[i] = values[col]
			}
			table.Add(row...)
		}
		table.Print()
		return nil
	}

}

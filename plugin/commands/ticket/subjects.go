package ticket

import (
	"strconv"

	"github.com/spf13/cobra"
	
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SubjectsTicketCommand struct {
	*metadata.SoftlayerCommand
	TicketManager managers.TicketManager
	Command       *cobra.Command
}

func NewSubjectsTicketCommand(sl *metadata.SoftlayerCommand) *SubjectsTicketCommand {
	thisCmd := &SubjectsTicketCommand{
		SoftlayerCommand: sl,
		TicketManager:    managers.NewTicketManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "subjects",
		Short: T("List Subject IDs for ticket creation"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SubjectsTicketCommand) Run(args []string) error {

	subjects, err := cmd.TicketManager.GetSubjects()

	if err != nil {
		return errors.New(T("Error: {{.Error}}", map[string]interface{}{"Error": err.Error()}))
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

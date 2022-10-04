package ticket

import (
	"strconv"

	"github.com/spf13/cobra"
	
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type SummaryTicketCommand struct {
	*metadata.SoftlayerCommand
	TicketManager managers.TicketManager
	UserManager   managers.UserManager
	Command       *cobra.Command
}

func NewSummaryTicketCommand(sl *metadata.SoftlayerCommand) *SummaryTicketCommand {
	thisCmd := &SummaryTicketCommand{
		SoftlayerCommand: sl,
		TicketManager:    managers.NewTicketManager(sl.Session),
		UserManager:      managers.NewUserManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "summary",
		Short: T("Summary info about tickets"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SummaryTicketCommand) Run(args []string) error {
	summary, err := cmd.TicketManager.Summary()

	if err != nil {
		return errors.New(T("Error: {{.Error}}.\n", map[string]interface{}{"Error": err.Error()}))
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

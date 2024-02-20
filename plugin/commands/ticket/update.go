package ticket

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type UpdateTicketCommand struct {
	*metadata.SoftlayerCommand
	TicketManager managers.TicketManager
	Command       *cobra.Command
}

func NewUpdateTicketCommand(sl *metadata.SoftlayerCommand) *UpdateTicketCommand {
	thisCmd := &UpdateTicketCommand{
		SoftlayerCommand: sl,
		TicketManager:    managers.NewTicketManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "update " + T("TICKETID"),
		Short: T("Adds an update to an existing ticket"),
		Long: T(`${COMMAND_NAME} sl ticket update TICKETID ["CONTENTS"] 
  
    If the second argument is not specified on a non-Windows machine, it will attempt to use either the value stored in the EDITOR environmental variable, or find either nano, vim, or emacs in that order.
  
EXAMPLE:
  ${COMMAND_NAME} sl ticket update 767676 "A problem has been detected."
  ${COMMAND_NAME} sl ticket update 767667`),
		Args: metadata.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *UpdateTicketCommand) Run(args []string) error {
	nargs := args
	ticketid, err := strconv.Atoi(args[0])
	if err != nil || ticketid <= 0 {
		return errors.NewInvalidUsageError(T("The ticket id must be a positive non-zero number."))
	}

	content := ""

	if len(nargs) == 1 {
		content, err = cmd.TicketManager.GetText()
		if err != nil {
			return err
		}
	} else {
		content = args[1]
	}

	err = cmd.TicketManager.AddUpdate(ticketid, content)
	if err != nil {
		return errors.New(T("Update could not be added: {{.Error}}\n", map[string]interface{}{"Error": err.Error()}))
	}
	cmd.UI.Ok()
	return nil
}

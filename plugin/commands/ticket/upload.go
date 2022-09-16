package ticket

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type UploadFileTicketCommand struct {
	*metadata.SoftlayerCommand
	TicketManager managers.TicketManager
	Command       *cobra.Command
	Name          string
}

func NewUploadFileTicketCommand(sl *metadata.SoftlayerCommand) *UploadFileTicketCommand {
	thisCmd := &UploadFileTicketCommand{
		SoftlayerCommand: sl,
		TicketManager:    managers.NewTicketManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "upload " + T("TICKETID") + " " + T("FILEPATH"),
		Short: T("Adds an attachment to an existing ticket"),
		Long: T(`${COMMAND_NAME} sl ticket upload TICKETID FILEPATH
  
EXAMPLE:
	${COMMAND_NAME} sl ticket upload 767676 "/home/user/screenshot.png"`),
		Args: metadata.TwoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Name, "name", "", T("The name of the attachment shown in the ticket"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *UploadFileTicketCommand) Run(args []string) error {
	ticketId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidUsageError(T("The ticket id must be a number."))
	}

	file_path := args[1]
	name := cmd.Name

	err = cmd.TicketManager.AttachFileToTicket(ticketId, name, file_path)

	if err != nil {
		return cli.NewExitError(T("Error: {{.Error}}", map[string]interface{}{"Error": err.Error()}), 2)
	} else {
		cmd.UI.Ok()
		return nil
	}
}

package ticket

import (
	"strconv"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	"github.com/spf13/cobra"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type AttachDeviceTicketCommand struct {
	*metadata.SoftlayerCommand
	TicketManager managers.TicketManager
	Command       *cobra.Command
	Hardware      int
	Virtual       int
}

func NewAttachDeviceTicketCommand(sl *metadata.SoftlayerCommand) *AttachDeviceTicketCommand {
	thisCmd := &AttachDeviceTicketCommand{
		SoftlayerCommand: sl,
		TicketManager:    managers.NewTicketManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "attach",
		Short: T("Attach devices to ticket"),
		Long: T(`${COMMAND_NAME} sl ticket attach TICKETID [OPTIONS]
  
EXAMPLE:
  ${COMMAND_NAME} sl ticket attach 7676767 --hardware 8675654 
  ${COMMAND_NAME} sl ticket attach 7676767 --virtual 1234567 `),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.Hardware, "hardware", 0, T("The identifier for hardware to attach"))
	cobraCmd.Flags().IntVar(&thisCmd.Virtual, "virtual", 0, T("The identifier for a virtual server to attach"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *AttachDeviceTicketCommand) Run(args []string) error {
	ticketid, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidUsageError(T("The ticket id must be a number."))
	}

	if cmd.Hardware != 0 && cmd.Virtual != 0 {
		return errors.NewInvalidUsageError(T("hardware and virtual flags cannot be set at the same time."))
	} else if cmd.Hardware == 0 && cmd.Virtual == 0 {
		return errors.NewInvalidUsageError(T("either the hardware or virtual flag must be set."))
	}

	if cmd.Hardware != 0 {
		deviceid := cmd.Hardware
		err = cmd.TicketManager.AttachDeviceToTicket(ticketid, deviceid, true)
	} else if cmd.Virtual != 0 {
		deviceid := cmd.Virtual
		err = cmd.TicketManager.AttachDeviceToTicket(ticketid, deviceid, false)
	}

	if err != nil {
		return cli.NewExitError(T("Error: {{.Error}}", map[string]interface{}{"Error": err.Error()}), 2)
	} else {
		cmd.UI.Ok()
		return nil
	}
}

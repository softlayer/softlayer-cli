package virtual

import (
	"time"

	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ReadyCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Wait                 int
}

func NewReadyCommand(sl *metadata.SoftlayerCommand) (cmd *ReadyCommand) {
	thisCmd := &ReadyCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "ready " + T("IDENTIFIER"),
		Short: T("Check if a virtual server instance is ready for use"),
		Long: T(`${COMMAND_NAME} sl vs ready IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs ready 12345678 --wait 30
   This command checks virtual server instance with ID 12345678 status to see if it is ready for use continuously and waits up to 30 seconds.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().IntVar(&thisCmd.Wait, "hard", 30, T("Wait until the virtual server is finished provisioning for up to X seconds before returning"))
	return thisCmd
}

func (cmd *ReadyCommand) Run(args []string) error {
	vsID, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	until := time.Now().Add(time.Duration(cmd.Wait) * time.Second)
	ready, message, err := cmd.VirtualServerManager.InstanceIsReady(vsID, until)
	subs := map[string]interface{}{"VsID": vsID, "VsId": vsID}
	if err != nil {
		return slErrors.NewAPIError(T("Failed to check virtual server instance {{.VsID}} is ready.\n", subs), err.Error(), 2)
	}
	if ready {
		cmd.UI.Print(T("Virtual server instance: {{.VsId}} is ready.", subs))
	} else {
		cmd.UI.Print(T("Not ready: {{.Message}}", map[string]interface{}{"Message": message}))
	}
	return nil
}

package dedicatedhost

import (
	"strconv"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CancellHostCommand struct {
	*metadata.SoftlayerCommand
	DedicatedHostManager managers.DedicatedHostManager
	Command              *cobra.Command
}

func NewCancelHostCommand(sl *metadata.SoftlayerCommand) *CancellHostCommand {
	thisCmd := &CancellHostCommand{
		SoftlayerCommand:     sl,
		DedicatedHostManager: managers.NewDedicatedhostManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "cancel " + T("IDENTIFIER"),
		Short: T("Cancel a dedicated host server immediately."),
		Long: T(`If there are any guests on this Dedicated Host, this command will fail until those guests are deleted.
Use 'sl dedicatedhost cancel-guests [IDENTIFIER]' to remove all guests from a host.
Use 'sl vs delete' to remove a specific guest.`),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CancellHostCommand) Run(args []string) error {
	hardwareID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("IDENTIFIER")
	}
	err = cmd.DedicatedHostManager.DeleteHost(hardwareID)
	if err != nil {
		return err
	}
	cmd.UI.Print(T("Dedicated Host {{.ID}} was cancelled", map[string]interface{}{"ID": hardwareID}))
	return nil
}

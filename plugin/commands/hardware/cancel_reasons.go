package hardware

import (
	"github.com/spf13/cobra"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CancelReasonsCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
}

func NewCancelReasonsCommand(sl *metadata.SoftlayerCommand) (cmd *CancelReasonsCommand) {
	thisCmd := &CancelReasonsCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "cancel-reasons",
		Short: T("Display a list of cancellation reasons"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CancelReasonsCommand) Run(args []string) error {
	reasons := cmd.HardwareManager.GetCancellationReasons()
	table := cmd.UI.Table([]string{T("Code"), T("Reason")})
	for key, value := range reasons {
		table.Add(key, value)
	}
	table.Print()
	return nil
}

package hardware

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

// ToggleIPMICommand is the implementation of Softlayer `hardware toggle-ipmi` command
type ToggleIPMICommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
	Enable          bool
	Disable         bool
	QuietFlag       bool
}

// NewToggleIPMICommand will create an instance of ToggleIPMICommand
func NewToggleIPMICommand(sl *metadata.SoftlayerCommand) (cmd *ToggleIPMICommand) {
	thisCmd := &ToggleIPMICommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "toggle-ipmi " + T("IDENTIFIER"),
		Short: T("Toggle the IPMI interface on and off. This command is asynchronous."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVar(&thisCmd.Enable, "enable", false, T("Enable the IPMI interface."))
	cobraCmd.Flags().BoolVar(&thisCmd.Disable, "disable", false, T("Disable the IPMI interface."))
	cobraCmd.Flags().BoolVarP(&thisCmd.QuietFlag, "quiet", "q", false, T("Suppress verbose output"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

// Run will execute ToggleIPMICommand
func (cmd *ToggleIPMICommand) Run(args []string) error {
	hardwareID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	if cmd.Enable && cmd.Disable {
		return errors.NewExclusiveFlagsErrorWithDetails([]string{"--enable", "--disable"}, "")
	}

	if !cmd.Enable && !cmd.Disable {
		return errors.NewInvalidUsageError(T("Either '--enable' or '--disable' is required."))
	}

	enabled := !cmd.Disable
	if err = cmd.HardwareManager.ToggleIPMI(hardwareID, enabled); err != nil {
		return errors.NewAPIError(T("Failed to toggle IPMI interface of hardware server '{{.ID}}'.\n", map[string]interface{}{"ID": hardwareID}), err.Error(), 2)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Successfully send request to toggle IPMI interface of hardware server '{{.ID}}'.", map[string]interface{}{"ID": hardwareID}))
	return nil
}

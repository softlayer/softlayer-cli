package vlan

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CancelCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	Force          bool
}

func NewCancelCommand(sl *metadata.SoftlayerCommand) *CancelCommand {
	thisCmd := &CancelCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "cancel " + T("IDENTIFIER"),
		Short: T("Cancel a VLAN."),
		Long: T(`${COMMAND_NAME} sl vlan cancel IDENTIFIER [OPTIONS]
	
EXAMPLE:
	${COMMAND_NAME} sl vlan cancel 12345678 -f
	This command cancels vlan with ID 12345678 without asking for confirmation.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CancelCommand) Run(args []string) error {
	vlanID, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("VLAN ID")
	}
	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will cancel the VLAN: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": vlanID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	// See if the API will just tell us if this VLAN can't be cancelled for a specific reason
	reasons := cmd.NetworkManager.GetCancelFailureReasons(vlanID)
	if len(reasons) > 0 {
		for _, reason := range reasons {
			cmd.UI.Print(reason)
		}
		return cli.NewExitError(T("Failed to cancel VLAN {{.ID}}.\n", map[string]interface{}{"ID": vlanID}), 2)

	}
	err = cmd.NetworkManager.CancelVLAN(vlanID)
	if err != nil {
		if strings.Contains(err.Error(), errors.SL_EXP_OBJ_NOT_FOUND) {
			return errors.NewAPIError(T("Unable to find VLAN with ID {{.ID}}.\n", map[string]interface{}{"ID": vlanID}), err.Error(), 0)
		}
		return errors.NewAPIError(T("Failed to cancel VLAN {{.ID}}.\n", map[string]interface{}{"ID": vlanID}), err.Error(), 2)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("VLAN {{.ID}} was cancelled.", map[string]interface{}{"ID": vlanID}))
	return nil
}

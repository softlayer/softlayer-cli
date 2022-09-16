package subnet

import (
	"strconv"
	"strings"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/spf13/cobra"
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
		Short: T("Cancel a subnet"),
		Long: T(`${COMMAND_NAME} sl subnet cancel IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl subnet cancel 12345678 -f
   This command cancels subnet with ID 12345678 without asking for confirmation.`),
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

	subnetID, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Subnet ID")
	}
	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will cancel the subnet: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": subnetID}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.NetworkManager.CancelSubnet(subnetID)
	if err != nil {
		if strings.Contains(err.Error(), errors.SL_EXP_OBJ_NOT_FOUND) {
			return errors.NewAPIError(T("Unable to find subnet with ID: {{.ID}}.\n", map[string]interface{}{"ID": subnetID}), err.Error(), 0)
		}
		return errors.NewAPIError(T("Failed to cancel subnet: {{.ID}}.\n", map[string]interface{}{"ID": subnetID}), err.Error(), 2)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Subnet {{.ID}} was cancelled.", map[string]interface{}{"ID": subnetID}))
	return nil
}

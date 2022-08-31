package globalip

import (
	"strings"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
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
		Short: T("Cancel a global IP."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CancelCommand) Run(args []string) error {

	globalIPID, err := utils.ResolveGloablIPId(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Globalip ID")
	}
	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will cancel the IP address: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": globalIPID}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.NetworkManager.CancelGlobalIP(globalIPID)
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return errors.NewAPIError(T("Unable to find global IP with ID: {{.ID}}.", map[string]interface{}{"ID": globalIPID}), err.Error(), 0)
		}
		return errors.NewAPIError(T("Failed to cancel global IP: {{.ID}}.", map[string]interface{}{"ID": globalIPID}), err.Error(), 2)

	}

	cmd.UI.Ok()
	cmd.UI.Print(T("IP address {{.ID}} was cancelled.", map[string]interface{}{"ID": globalIPID}))
	return nil
}

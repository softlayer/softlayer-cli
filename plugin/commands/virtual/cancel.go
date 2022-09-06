package virtual

import (
	"strings"

	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CancelCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Force                bool
}

func NewCancelCommand(sl *metadata.SoftlayerCommand) (cmd *CancelCommand) {
	thisCmd := &CancelCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "cancel " + T("IDENTIFIER"),
		Short: T("Cancel virtual server instance"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	return thisCmd
}

func (cmd *CancelCommand) Run(args []string) error {

	VsID, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	subs := map[string]interface{}{"VsID": VsID, "VsId": VsID}
	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will cancel the virtual server instance: {{.VsID}} and cannot be undone. Continue?", subs))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.VirtualServerManager.CancelInstance(VsID)
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return slErrors.NewAPIError(T("Unable to find virtual server instance with ID: {{.VsID}}.\n", subs), err.Error(), 0)
		}
		return slErrors.NewAPIError(T("Failed to cancel virtual server instance: {{.VsID}}.\n", subs), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Virtual server instance: {{.VsId}} was cancelled.", subs))
	return nil
}

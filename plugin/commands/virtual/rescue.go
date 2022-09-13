package virtual

import (
	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type RescueCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Force                bool
}

func NewRescueCommand(sl *metadata.SoftlayerCommand) (cmd *RescueCommand) {
	thisCmd := &RescueCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "rescue " + T("IDENTIFIER"),
		Short: T("Reboot a virtual server instance into a rescue image"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	return thisCmd
}

func (cmd *RescueCommand) Run(args []string) error {

	vsID, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}
	subs := map[string]interface{}{"VsId": vsID, "VsID": vsID}
	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will reboot virtual server instance: {{.VsId}}. Continue?", subs))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.VirtualServerManager.RescueInstance(vsID)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to rescue virtual server instance: {{.VsID}}.\n", subs), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Virtual server instance: {{.VsId}} was rebooted.", subs))
	return nil
}

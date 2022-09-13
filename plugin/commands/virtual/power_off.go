package virtual

import (
	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PowerOffCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Hard                 bool
	Soft                 bool
	Force                bool
}

func NewPowerOffCommand(sl *metadata.SoftlayerCommand) (cmd *PowerOffCommand) {
	thisCmd := &PowerOffCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "power-off " + T("IDENTIFIER"),
		Short: T("Power off an active virtual server instance"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().BoolVar(&thisCmd.Hard, "hard", false, T("Perform a hard shutdown"))
	cobraCmd.Flags().BoolVar(&thisCmd.Soft, "soft", false, T("Perform a soft shutdown"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	return thisCmd
}

func (cmd *PowerOffCommand) Run(args []string) error {

	vsID, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	if cmd.Hard && cmd.Soft {
		return slErrors.NewExclusiveFlagsError("--hard", "--soft")
	}

	subs := map[string]interface{}{"VsId": vsID, "VsID": vsID}
	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will power off virtual server instance: {{.VsId}}. Continue?", subs))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.VirtualServerManager.PowerOffInstance(vsID, cmd.Soft, cmd.Hard)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to power off virtual server instance: {{.VsID}}.\n", subs), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Virtual server instance: {{.VsId}} was power off.", subs))
	return nil
}

package virtual

import (
	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ReloadCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Postinstall          string
	Image                int
	Key                  []int
	Force                bool
}

func NewReloadCommand(sl *metadata.SoftlayerCommand) (cmd *ReloadCommand) {
	thisCmd := &ReloadCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "reload " + T("IDENTIFIER"),
		Short: T("Reload operating system on a virtual server instance"),
		Long: T(`${COMMAND_NAME} sl vs reload IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl vs reload 12345678
   This command reloads current operating system for virtual server instance with ID 12345678.
   ${COMMAND_NAME} sl vs reload 12345678 --image 1234
   This command reloads operating system from image with ID 1234 for virtual server instance with ID 12345678.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().StringVarP(&thisCmd.Postinstall, "postinstall", "i", "", T("Post-install script to download"))
	cobraCmd.Flags().IntVar(&thisCmd.Image, "image", 0, T("Image ID. The default is to use the current operating system.\nSee: '${COMMAND_NAME} sl image list' for reference"))
	cobraCmd.Flags().IntSliceVarP(&thisCmd.Key, "key", "k", []int{}, T("The IDs of the SSH keys to add to the root user (multiple occurrence permitted)"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))

	return thisCmd
}

func (cmd *ReloadCommand) Run(args []string) error {

	vsID, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	subs := map[string]interface{}{"VsId": vsID, "VsID": vsID, "CommandName": "ibmcloud"}
	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will reload operating system of virtual server instance: {{.VsId}} and cannot be undone. Continue?", subs))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.VirtualServerManager.ReloadInstance(vsID, cmd.Postinstall, cmd.Key, cmd.Image)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to reload virtual server instance: {{.VsID}}.\n", subs), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("System reloading for virtual server instance: {{.VsId}} is in progress. Run '{{.CommandName}} sl vs ready {{.VsId}}' to check whether it is ready later on.", subs))
	return nil
}

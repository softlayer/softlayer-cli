package dedicatedhost

import (
	"github.com/spf13/cobra"
	"github.com/urfave/cli"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CancelCommand struct {
	*metadata.SoftlayerCommand
	DedicatedHostManager managers.DedicatedHostManager
	Command              *cobra.Command
	Force                bool
}

func NewCancelCommand(sl *metadata.SoftlayerCommand) *CancelCommand {
	thisCmd := &CancelCommand{
		SoftlayerCommand:     sl,
		DedicatedHostManager: managers.NewDedicatedhostManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "cancel-guests",
		Short: T("Cancel all virtual guests of the dedicated host immediately."),
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
	HostID, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Host ID")
	}

	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will cancel all virtual server instances in the dedicatedhost: {{.HostID}} and cannot be undone. Continue?", map[string]interface{}{"HostID": HostID}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	listGuest, err := cmd.DedicatedHostManager.CancelGuests(HostID)
	if err != nil {
		return slErr.NewAPIError(T("Failed to cancel all guests in the dedicatedhost: {{.HostID}}.", map[string]interface{}{"HostID": HostID}), err.Error(), 2)
	}

	if len(listGuest) > 0 {
		table := cmd.UI.Table([]string{T("Id"), T("Server Name"), T("Status")})
		for _, guest := range listGuest {
			table.Add(utils.FormatIntPointer(&guest.Id), utils.FormatStringPointer(&guest.Fqdn), utils.FormatStringPointer(&guest.Status))
		}
		table.Print()
		return nil
	}

	return cli.NewExitError(T("There is not any guest into the dedicated host {{.ID}}.", map[string]interface{}{"ID": HostID}), 2)
}

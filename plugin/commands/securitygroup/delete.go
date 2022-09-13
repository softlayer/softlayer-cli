package securitygroup

import (
	"strconv"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/spf13/cobra"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type DeleteCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	ForceFlag      bool
}

func NewDeleteCommand(sl *metadata.SoftlayerCommand) (cmd *DeleteCommand) {
	thisCmd := &DeleteCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "delete " + T("SECURITYGROUP_ID"),
		Short: T("Delete the given security group"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVarP(&thisCmd.ForceFlag, "force", "f", false, T("Force operation without confirmation"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DeleteCommand) Run(args []string) error {
	groupID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group ID")
	}
	if !cmd.ForceFlag {
		confirm, err := cmd.UI.Confirm(T("This will delete security group {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": groupID}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.NetworkManager.DeleteSecurityGroup(groupID)
	if err != nil {
		return errors.NewAPIError(T("Failed to delete security group {{.ID}}.\n", map[string]interface{}{"ID": groupID}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Security group {{.ID}} is deleted.", map[string]interface{}{"ID": groupID}))
	return nil
}

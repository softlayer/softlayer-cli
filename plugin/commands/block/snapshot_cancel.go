package block

import (
	"strconv"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type SnapshotCancelCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	Reason         string
	Immediate      bool
	Force          bool
}

func NewSnapshotCancelCommand(sl *metadata.SoftlayerStorageCommand) *SnapshotCancelCommand {
	thisCmd := &SnapshotCancelCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "snapshot-cancel " + T("IDENTIFIER"),
		Short: T("Cancel existing snapshot space for a given volume"),
		Long: T(`${COMMAND_NAME} sl {{.storageType}} snapshot-cancel SNAPSHOT_ID [OPTIONS]

Cancel existing snapshot space for a given volume.

EXAMPLE:
   ${COMMAND_NAME} sl {{.storageType}} snapshot-cancel 12345678 --immediate -f 
   This command cancels snapshot with ID 12345678 immediately without asking for confirmation.`, sl.StorageI18n),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Reason, "reason", "", T("An optional reason for cancellation"))
	cobraCmd.Flags().BoolVar(&thisCmd.Immediate, "immediate", false, T("Cancel the snapshot space immediately instead of on the billing anniversary"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SnapshotCancelCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	subs := map[string]interface{}{"ID": volumeID}
	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will cancel the block volume snapshot space: {{.ID}} and cannot be undone. Continue?", subs))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.StorageManager.CancelSnapshotSpace("block", volumeID, cmd.Reason, cmd.Immediate)
	if err != nil {
		return slErr.NewAPIError(T("Failed to cancel snapshot space for volume {{.ID}}.\n", subs), err.Error(), 2)
	}

	cmd.UI.Ok()
	if cmd.Immediate {
		cmd.UI.Print(T("Block volume {{.ID}} has been marked for immediate snapshot cancellation.", subs))
	} else {
		cmd.UI.Print(T("Block volume {{.ID}} has been marked for snapshot cancellation.", subs))
	}
	return nil
}

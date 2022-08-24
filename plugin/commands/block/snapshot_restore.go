package block

import (
	"strconv"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type SnapshotRestoreCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
}

func NewSnapshotRestoreCommand(sl *metadata.SoftlayerStorageCommand) *SnapshotRestoreCommand {
	thisCmd := &SnapshotRestoreCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "snapshot-restore " + T("IDENTIFIER") + " " + T("SNAPSHOT_ID"),
		Short: T("Restore block volume using a given snapshot"),
		Long: T(`${COMMAND_NAME} sl block snapshot-restore VOLUME_ID SNAPSHOT_ID
	
EXAMPLE:
   ${COMMAND_NAME} sl block snapshot-restore 12345678 87654321
   This command restores volume with ID 12345678 from snapshot with ID 87654321.`),
		Args: metadata.TwoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SnapshotRestoreCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	snapshotID, err := strconv.Atoi(args[1])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Snapshot ID")
	}
	err = cmd.StorageManager.RestoreFromSnapshot(volumeID, snapshotID)
	// Easier to have 2 volumeIds than to change the translation strings, might fix one day.
	subs := map[string]interface{}{"SnapshotId": snapshotID, "VolumeID": volumeID, "VolumeId": volumeID}
	if err != nil {
		return slErr.NewAPIError(T("Failed to restore volume {{.VolumeID}} from snapshot {{.SnapshotId}}.\n", subs), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Block volume {{.VolumeId}} is being restored using snapshot {{.SnapshotId}}.", subs))
	return nil
}

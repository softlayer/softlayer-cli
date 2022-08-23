package block

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type SnapshotDeleteCommand struct {
	*metadata.SoftlayerCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
}

func NewSnapshotDeleteCommand(sl *metadata.SoftlayerCommand) *SnapshotDeleteCommand {
	thisCmd := &SnapshotDeleteCommand{
		SoftlayerCommand: sl,
		StorageManager:   managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "snapshot-delete " + T("IDENTIFIER"),
		Short: T("Delete a snapshot on a given volume"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SnapshotDeleteCommand) Run(args []string) error {

	snapshotID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Snapshot ID")
	}
	err = cmd.StorageManager.DeleteSnapshot(snapshotID)
	subs := map[string]interface{}{"ID": snapshotID, "SnapshotId": snapshotID}
	if err != nil {
		if strings.Contains(err.Error(), slErr.SL_EXP_OBJ_NOT_FOUND) {
			return slErr.NewAPIError(T("Unable to find snapshot with ID {{.ID}}.\n", subs), err.Error(), 0)
		}
		return slErr.NewAPIError(T("Failed to delete snapshot {{.ID}}.\n", subs), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Snapshot {{.SnapshotId}} was deleted.", subs))
	return nil
}

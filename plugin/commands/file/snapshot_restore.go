package file

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
)

type SnapshotRestoreCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewSnapshotRestoreCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *SnapshotRestoreCommand) {
	return &SnapshotRestoreCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func (cmd *SnapshotRestoreCommand) Run(c *cli.Context) error {
	if c.NArg() != 2 {
		return errors.NewInvalidUsageError(T("This command requires two arguments."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	snapshotID, err := strconv.Atoi(c.Args()[1])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Snapshot ID")
	}
	err = cmd.StorageManager.RestoreFromSnapshot(volumeID, snapshotID)
	if err != nil {
		return cli.NewExitError(T("Failed to restore volume {{.VolumeID}} from snapshot {{.SnapshotId}}.\n",
			map[string]interface{}{"SnapshotId": snapshotID, "VolumeID": volumeID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("File volume {{.VolumeId}} is being restored using snapshot {{.SnapshotId}}.",
		map[string]interface{}{"SnapshotId": snapshotID, "VolumeId": volumeID}))
	return nil
}

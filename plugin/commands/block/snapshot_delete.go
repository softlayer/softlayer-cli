package block

import (
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
)

type SnapshotDeleteCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewSnapshotDeleteCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *SnapshotDeleteCommand) {
	return &SnapshotDeleteCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func (cmd *SnapshotDeleteCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	snapshotID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Snapshot ID")
	}
	err = cmd.StorageManager.DeleteSnapshot(snapshotID)
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return cli.NewExitError(T("Unable to find snapshot with ID {{.ID}}.\n", map[string]interface{}{"ID": snapshotID})+err.Error(), 0)
		}
		return cli.NewExitError(T("Failed to delete snapshot {{.ID}}.\n", map[string]interface{}{"ID": snapshotID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Snapshot {{.SnapshotId}} was deleted.", map[string]interface{}{"SnapshotId": snapshotID}))
	return nil
}

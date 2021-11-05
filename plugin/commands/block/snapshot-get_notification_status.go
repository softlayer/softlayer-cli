package block

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type SnapshotGetNotificationStatusCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewSnapshotGetNotificationStatusCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *SnapshotGetNotificationStatusCommand) {
	return &SnapshotGetNotificationStatusCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func (cmd *SnapshotGetNotificationStatusCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	var snapshot string
	snapshot, err = cmd.StorageManager.GetSnapshotNotificationStatus(volumeID)
	if err != nil {
		return cli.NewExitError(T("Failed to get the snapshot notification status for the volume '{{.ID}}'.\n", map[string]interface{}{"ID": volumeID})+err.Error(), 2)
	}

	if snapshot == "" {
		cmd.UI.Print(T("Snapshots space usage threshold warning flag setting is null. Set to default value enable. For volume '{{.ID}}'.", map[string]interface{}{"ID": volumeID}))
	} else {
		cmd.UI.Print(T("Snapshots space usage threshold warning flag setting is enabled for volume '{{.ID}}'.", map[string]interface{}{"ID": volumeID}))
		cmd.UI.Print(T("Snapshots Notification Status: '{{.STATUS}}'.", map[string]interface{}{"STATUS": snapshot}))
	}

	return nil
}

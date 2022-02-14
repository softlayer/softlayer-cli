package block

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
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

func BlockVolumeSnapshotGetNotificationStatusMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "snapshot-get-notification-status",
		Description: T("Get snapshots space usage threshold warning flag setting for a given volume."),
		Usage: T(`${COMMAND_NAME} sl block snapshot-get-notification-status VOLUME_ID

EXAMPLE:
	${COMMAND_NAME} sl block snapshot-get-notification-status VOLUME_ID
	Get snapshots space usage threshold warning flag setting for a given volume.`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
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

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	enabled, err := cmd.StorageManager.GetSnapshotNotificationStatus(volumeID)
	if err != nil {
		return cli.NewExitError(T("Failed to get the snapshot notification status for volume '{{.ID}}'.\n", map[string]interface{}{"ID": volumeID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, enabled)
	}

	if enabled == 0 {
		cmd.UI.Print(T("Disabled: Snapshots space usage threshold is disabled for volume '{{.ID}}'.", map[string]interface{}{"ID": volumeID}))
	} else {
		cmd.UI.Print(T("Enabled: Snapshots space usage threshold is enabled for volume '{{.ID}}'.", map[string]interface{}{"ID": volumeID}))
	}

	return nil
}

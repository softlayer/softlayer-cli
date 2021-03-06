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

type SnapshotSetNotificationCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewSnapshotSetNotificationCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *SnapshotSetNotificationCommand) {
	return &SnapshotSetNotificationCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func BlockVolumeSnapshotSetNotificationMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "snapshot-set-notification",
		Description: T("Enables/Disables snapshot space usage threshold warning for a given volume."),
		Usage: T(`${COMMAND_NAME} sl block  snapshot-set-notification VOLUME_ID

EXAMPLE:
	${COMMAND_NAME} sl block snapshot-set-notification --enable 1234567
	Enables/Disables snapshot space usage threshold warning for a given volume.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "enable",
				Usage: T("Enable snapshot notification."),
			},
			cli.BoolFlag{
				Name:  "disable",
				Usage: T("Disable snapshot notification."),
			},
		},
	}
}

func (cmd *SnapshotSetNotificationCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	if c.IsSet("enable") && c.IsSet("disable") {
		return errors.NewExclusiveFlagsErrorWithDetails([]string{"--enable", "--disable"}, "")
	}

	if !c.IsSet("enable") && !c.IsSet("disable") {
		return errors.NewInvalidUsageError(T("Either '--enable' or '--disable' is required."))
	}

	enabled := !c.IsSet("disable")
	if err = cmd.StorageManager.SetSnapshotNotification(volumeID, enabled); err != nil {
		return cli.NewExitError(T("Failed to set the snapshort notification  for volume '{{.ID}}'.\n", map[string]interface{}{"ID": volumeID})+err.Error(), 2)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Snapshots space usage threshold warning notification has been set to '{{.ENABLE}}' for volume '{{.ID}}'.", map[string]interface{}{"ID": volumeID, "ENABLE": enabled}))
	return nil
}

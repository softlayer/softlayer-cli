package file

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type SnapshotDisableCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewSnapshotDisableCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *SnapshotDisableCommand) {
	return &SnapshotDisableCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func FileSnapshotDisableMetaData() cli.Command {
	return cli.Command{
		Category:    "file",
		Name:        "snapshot-disable",
		Description: T("Disable snapshots on the specified schedule for a given volume"),
		Usage: T(`${COMMAND_NAME} sl file snapshot-disable VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file snapshot-disable 12345678 -s DAILY
   This command disables daily snapshot for volume with ID 12345678.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "s,schedule-type",
				Usage: T("Snapshot schedule [required], options are: HOURLY,DAILY,WEEKLY"),
			},
		},
	}
}

func (cmd *SnapshotDisableCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	if !c.IsSet("schedule-type") {
		return errors.NewInvalidUsageError(T("[--schedule-type] is required, options are: HOURLY, DAILY, WEEKLY."))
	}
	scheduleType := c.String("schedule-type")
	if scheduleType != "HOURLY" && scheduleType != "DAILY" && scheduleType != "WEEKLY" {
		return errors.NewInvalidUsageError(T("[--schedule-type] must be HOURLY, DAILY, or WEEKLY."))
	}
	err = cmd.StorageManager.DisableSnapshots(volumeID, scheduleType)
	if err != nil {
		return cli.NewExitError(T("Failed to disable {{.ScheduleType}} snapshot for volume {{.VolumeID}}.\n",
			map[string]interface{}{"ScheduleType": scheduleType, "VolumeID": volumeID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("{{.ScheduleType}} snapshots have been disabled for volume {{.VolumeID}}.",
		map[string]interface{}{"ScheduleType": scheduleType, "VolumeID": volumeID}))
	return nil
}

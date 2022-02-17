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

type SnapshotEnableCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewSnapshotEnableCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *SnapshotEnableCommand) {
	return &SnapshotEnableCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

var DAY_OF_WEEK = map[int]string{
	0: "SUNDAY",
	1: "MONDAY",
	2: "TUESDAY",
	3: "WEDNESDAY",
	4: "THURSDAY",
	5: "FRIDAY",
	6: "SATURDAY",
}

func FileSnapshotEnableMetaData() cli.Command {
	return cli.Command{
		Category:    "file",
		Name:        "snapshot-enable",
		Description: T("Enable snapshots for a given volume on the specified schedule"),
		Usage: T(`${COMMAND_NAME} sl file snapshot-enable VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file snapshot-enable 12345678 -s WEEKLY -c 5 -m 0 --hour 2 -d 0
   This command enables snapshot for volume with ID 12345678, snapshot is taken weekly on every Sunday at 2:00, and up to 5 snapshots are retained.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "s,schedule-type",
				Usage: T("Snapshot schedule [required], options are: HOURLY,DAILY,WEEKLY"),
			},
			cli.IntFlag{
				Name:  "c,retention-count",
				Usage: T("Number of snapshots to retain [required]"),
			},
			cli.IntFlag{
				Name:  "m,minute",
				Usage: T("Minute of the hour when snapshots should be taken, integer between 0 to 59"),
			},
			cli.IntFlag{
				Name:  "r,hour",
				Usage: T("Hour of the day when snapshots should be taken, integer between 0 to 23"),
			},
			cli.IntFlag{
				Name:  "d,day-of-week",
				Usage: T("Day of the week when snapshots should be taken, integer between 0 to 6. \n      0 means Sunday,1 means Monday,2 means Tuesday,3 means Wendesday,4 means Thursday,5 means Friday,6 means Saturday"),
			},
		},
	}
}

func (cmd *SnapshotEnableCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	if !c.IsSet("schedule-type") {
		return errors.NewInvalidUsageError(T("[-s|--schedule-type] is required, options are: HOURLY, DAILY, WEEKLY."))
	}
	scheduleType := c.String("schedule-type")
	if scheduleType != "HOURLY" && scheduleType != "DAILY" && scheduleType != "WEEKLY" {
		return errors.NewInvalidUsageError(T("[-s|--schedule-type] must be HOURLY, DAILY, or WEEKLY."))
	}
	if !c.IsSet("retention-count") {
		return errors.NewMissingInputError("-c|--retention-count")
	}
	retentionCount := c.Int("retention-count")
	minute := c.Int("m")
	if minute < 0 || minute > 59 {
		return errors.NewInvalidUsageError(T("[-m|--minute] value must be between 0 and 59."))
	}
	hour := c.Int("r")
	if hour < 0 || hour > 23 {
		return errors.NewInvalidUsageError(T("[-r|--hour] value must be between 0 and 23."))
	}
	dayOfWeek := c.Int("d")
	if dayOfWeek < 0 || dayOfWeek > 6 {
		return errors.NewInvalidUsageError(T("[-d|--day-of-week] value must be between 0 and 6."))
	}
	err = cmd.StorageManager.EnableSnapshot(volumeID, scheduleType, retentionCount, minute, hour, DAY_OF_WEEK[dayOfWeek])
	if err != nil {
		return cli.NewExitError(T("Failed to enable {{.ScheduleType}} snapshot for volume {{.VolumeID}}.\n",
			map[string]interface{}{"ScheduleType": scheduleType, "VolumeID": volumeID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("{{.ScheduleType}} snapshots have been enabled for volume {{.VolumeID}}.",
		map[string]interface{}{"ScheduleType": scheduleType, "VolumeID": volumeID}))
	return nil
}

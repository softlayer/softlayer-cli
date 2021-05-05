package block

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
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

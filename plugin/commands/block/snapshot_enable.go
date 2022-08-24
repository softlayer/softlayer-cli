package block

import (
	"strconv"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

var DAY_OF_WEEK = map[int]string{
	0: "SUNDAY",
	1: "MONDAY",
	2: "TUESDAY",
	3: "WEDNESDAY",
	4: "THURSDAY",
	5: "FRIDAY",
	6: "SATURDAY",
}

type SnapshotEnableCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	ScheduleType   string
	RetentionCount int
	Minute         int
	Hour           int
	DayOfWeek      int
}

func NewSnapshotEnableCommand(sl *metadata.SoftlayerStorageCommand) *SnapshotEnableCommand {
	thisCmd := &SnapshotEnableCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "snapshot-enable " + T("IDENTIFIER"),
		Short: T("Enable snapshots for a given volume on the specified schedule"),
		Long: T(`${COMMAND_NAME} sl block snapshot-enable VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl block snapshot-enable 12345678 -s WEEKLY -c 5 -m 0 --hour 2 -d 0
   This command enables snapshot for volume with ID 12345678, snapshot is taken weekly on every Sunday at 2:00, and up to 5 snapshots are retained.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.ScheduleType, "schedule-type", "s", "", T("Snapshot schedule [required], options are: HOURLY,DAILY,WEEKLY"))
	cobraCmd.Flags().IntVarP(&thisCmd.RetentionCount, "retention-count", "c", 0, T("Number of snapshots to retain [required]"))
	cobraCmd.Flags().IntVarP(&thisCmd.Minute, "minute", "m", 0, T("Minute of the hour when snapshots should be taken, integer between 0 to 59"))
	cobraCmd.Flags().IntVarP(&thisCmd.Hour, "hour", "r", 0, T("Hour of the day when snapshots should be taken, integer between 0 to 23"))
	cobraCmd.Flags().IntVarP(&thisCmd.DayOfWeek, "day-of-week", "d", 0, T("Day of the week when snapshots should be taken, integer between 0 to 6. \n      0 means Sunday,1 means Monday,2 means Tuesday,3 means Wendesday,4 means Thursday,5 means Friday,6 means Saturday"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SnapshotEnableCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	if cmd.ScheduleType == "" {
		return slErr.NewInvalidUsageError(T("[-s|--schedule-type] is required, options are: HOURLY, DAILY, WEEKLY."))
	}
	scheduleType := cmd.ScheduleType
	if scheduleType != "HOURLY" && scheduleType != "DAILY" && scheduleType != "WEEKLY" {
		return slErr.NewInvalidUsageError(T("[-s|--schedule-type] must be HOURLY, DAILY, or WEEKLY."))
	}
	retentionCount := cmd.RetentionCount

	if retentionCount == 0 {
		return slErr.NewMissingInputError("-c|--retention-count")
	}

	minute := cmd.Minute
	if minute < 0 || minute > 59 {
		return slErr.NewInvalidUsageError(T("[-m|--minute] value must be between 0 and 59."))
	}
	hour := cmd.Hour
	if hour < 0 || hour > 23 {
		return slErr.NewInvalidUsageError(T("[-r|--hour] value must be between 0 and 23."))
	}
	dayOfWeek := cmd.DayOfWeek
	if dayOfWeek < 0 || dayOfWeek > 6 {
		return slErr.NewInvalidUsageError(T("[-d|--day-of-week] value must be between 0 and 6."))
	}
	err = cmd.StorageManager.EnableSnapshot(volumeID, scheduleType, retentionCount, minute, hour, DAY_OF_WEEK[dayOfWeek])
	subs := map[string]interface{}{"ScheduleType": scheduleType, "VolumeID": volumeID}
	if err != nil {
		return slErr.NewAPIError(T("Failed to enable {{.ScheduleType}} snapshot for volume {{.VolumeID}}.\n", subs), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("{{.ScheduleType}} snapshots have been enabled for volume {{.VolumeID}}.", subs))
	return nil
}

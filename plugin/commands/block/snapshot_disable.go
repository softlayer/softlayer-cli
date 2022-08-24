package block

import (
	"strconv"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type SnapshotDisableCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	Schedule_type  string
}

func NewSnapshotDisableCommand(sl *metadata.SoftlayerStorageCommand) *SnapshotDisableCommand {
	thisCmd := &SnapshotDisableCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "snapshot-disable " + T("IDENTIFIER"),
		Short: T("Disable snapshots on the specified schedule for a given volume"),
		Long: T(`${COMMAND_NAME} sl {{.storageType}} snapshot-disable VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl {{.storageType}} snapshot-disable 12345678 -s DAILY
   This command disables daily snapshot for volume with ID 12345678.`, sl.StorageI18n),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.Schedule_type, "schedule-type", "s", "", T("Snapshot schedule [required], options are: HOURLY,DAILY,WEEKLY"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SnapshotDisableCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	if cmd.Schedule_type == "" {
		return slErr.NewInvalidUsageError(T("[--schedule-type] is required, options are: HOURLY, DAILY, WEEKLY."))
	}
	scheduleType := cmd.Schedule_type
	if scheduleType != "HOURLY" && scheduleType != "DAILY" && scheduleType != "WEEKLY" {
		return slErr.NewInvalidUsageError(T("[--schedule-type] must be HOURLY, DAILY, or WEEKLY."))
	}
	err = cmd.StorageManager.DisableSnapshots(volumeID, scheduleType)
	subs := map[string]interface{}{"ScheduleType": scheduleType, "VolumeID": volumeID}
	if err != nil {
		return slErr.NewAPIError(T("Failed to disable {{.ScheduleType}} snapshot for volume {{.VolumeID}}.\n", subs), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("{{.ScheduleType}} snapshots have been disabled for volume {{.VolumeID}}.", subs))
	return nil
}

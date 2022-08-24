package block

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SnapshotScheduleListCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
}

func NewSnapshotScheduleListCommand(sl *metadata.SoftlayerStorageCommand) (cmd *SnapshotScheduleListCommand) {
	thisCmd := &SnapshotScheduleListCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "snapshot-schedule-list " + T("IDENTIFIER"),
		Short: T("List snapshot schedules for a given volume"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SnapshotScheduleListCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	outputFormat := cmd.GetOutputFlag()

	snapshotSchedules, err := cmd.StorageManager.GetVolumeSnapshotSchedules(volumeID)
	if err != nil {
		return slErr.NewAPIError(T("Failed to get snapshot schedule list on your account.\n"), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, snapshotSchedules.Schedules)
	}

	table := cmd.UI.Table([]string{T("id"),
		T("active"),
		T("type"),
		T("replication"),
		T("date_created"),
		T("minute"),
		T("hour"),
		T("day"),
		T("week"),
		T("day_of_week"),
		T("date_of_month"),
		T("month_of_year"),
		T("maximum_snapshots")})
	for _, sps := range snapshotSchedules.Schedules {
		var replication, fileScheduleType string

		if sps.Type != nil && sps.Type.Keyname != nil && strings.Contains(*sps.Type.Keyname, "REPLICATION") {
			replication = "*"
		} else {
			replication = "-"
		}

		if sps.Type != nil && sps.Type.Keyname != nil {
			fileScheduleType = strings.Replace(*sps.Type.Keyname, "REPLICATION_", "", -1)
		}
		fileScheduleType = strings.Replace(fileScheduleType, "SNAPSHOT_", "", -1)

		propertyList := []string{"MINUTE", "HOUR", "DAY", "WEEK",
			"DAY_OF_WEEK", "DAY_OF_MONTH",
			"MONTH_OF_YEAR", "SNAPSHOT_LIMIT"}
		scheduleProperties := []string{}
		for _, propKey := range propertyList {
			var item string
			item = "-"
			for _, scheduleProperty := range sps.Properties {
				if scheduleProperty.Type != nil && scheduleProperty.Type.Keyname != nil && *scheduleProperty.Type.Keyname == propKey {
					if scheduleProperty.Value != nil {
						if *scheduleProperty.Value == "-1" {
							item = "*"
						} else {
							item = *scheduleProperty.Value
						}
					}
					break
				}
			}
			scheduleProperties = append(scheduleProperties, item)
		}
		var active string
		if sps.Active != nil {
			active = "*"
		} else {
			active = ""
		}

		tableItem := []string{utils.FormatIntPointer(sps.Id),
			active,
			fileScheduleType,
			replication,
			utils.FormatSLTimePointer(sps.CreateDate),
		}

		tableItem = append(tableItem, scheduleProperties...)

		table.Add(tableItem...)

	}
	table.Print()
	return nil
}

package file

import (
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SnapshotScheduleListCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewSnapshotScheduleListCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *SnapshotScheduleListCommand) {
	return &SnapshotScheduleListCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func FileSnapshotScheduleListMetaData() cli.Command {
	return cli.Command{
		Category:    "file",
		Name:        "snapshot-schedule-list",
		Description: T("List snapshot schedules for a given volume"),
		Usage: T(`${COMMAND_NAME} sl file snapshot-schedule-list VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file snapshot-schedule-list 12345678
   This command list snapshot schedules for volume with ID 12345678`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}


func (cmd *SnapshotScheduleListCommand) Run(c *cli.Context) error {
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

	snapshotSchedules, err := cmd.StorageManager.GetVolumeSnapshotSchedules(volumeID)
	if err != nil {
		return cli.NewExitError(T("Failed to get snapshot schedule list on your account.\n")+err.Error(), 2)
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

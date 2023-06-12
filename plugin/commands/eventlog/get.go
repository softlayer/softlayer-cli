package eventlog

import (
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type GetCommand struct {
	*metadata.SoftlayerCommand
	EventLogManager managers.EventLogManager
	Command         *cobra.Command
	DateMin         string
	DateMax         string
	ObjEvent        string
	ObjId           string
	ObjType         string
	UtcOffset       string
	Metadata        bool
	Limit           int
}

func NewGetCommand(sl *metadata.SoftlayerCommand) (cmd *GetCommand) {
	thisCmd := &GetCommand{
		SoftlayerCommand: sl,
		EventLogManager:  managers.NewEventLogManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "get",
		Short: T("Get Event Logs"),
		Long: T(`${COMMAND_NAME} sl event-log get [OPTIONS]

EXAMPLE: 
   ${COMMAND_NAME} sl event-log get 
   ${COMMAND_NAME} sl event-log get --limit 5 --obj-id 123456 --obj-event Create --metadata
   ${COMMAND_NAME} sl event-log get --date-min 2021-03-31 --date-max 2021-04-31`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVarP(&thisCmd.DateMin, "date-min", "d", "", T("The earliest date we want to search for event logs [YYYY-MM-DD]."))
	cobraCmd.Flags().StringVarP(&thisCmd.DateMax, "date-max", "D", "", T("The latest date we want to search for event logs [YYYY-MM-DD]."))
	cobraCmd.Flags().StringVarP(&thisCmd.ObjEvent, "obj-event", "e", "", T("The event we want to get event logs for"))
	cobraCmd.Flags().StringVarP(&thisCmd.ObjId, "obj-id", "i", "", T("The id of the object we want to get event logs for"))
	cobraCmd.Flags().StringVarP(&thisCmd.ObjType, "obj-type", "t", "", T("The type of the object we want to get event logs for"))
	cobraCmd.Flags().StringVarP(&thisCmd.UtcOffset, "utc-offset", "z", "", T("UTC Offset for searching with dates. +/-HHMM format  [default: -0000]"))
	cobraCmd.Flags().BoolVar(&thisCmd.Metadata, "metadata", false, T("Display metadata if present  [default: no-metadata]"))
	cobraCmd.Flags().IntVarP(&thisCmd.Limit, "limit", "l", 50, T("Total number of result to return. -1 to return ALL, there may be a LOT of these.  [default: 50]"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *GetCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	limit := cmd.Limit

	dateMin := cmd.DateMin
	if dateMin != "" {
		time, err := time.Parse(time.RFC3339, dateMin+"T00:00:00Z")
		if err != nil {
			return errors.NewInvalidUsageError(T("Invalid format date to --date-min."))
		}
		dateMin = time.Format("2006-01-02T15:04:05")
	}

	dateMax := cmd.DateMax
	if dateMax != "" {
		time, err := time.Parse(time.RFC3339, dateMax+"T00:00:00Z")
		if err != nil {
			return errors.NewInvalidUsageError(T("Invalid format date to --date-max."))
		}
		dateMax = time.Format("2006-01-02T15:04:05")
	}

	objEvent := cmd.ObjEvent

	objId := cmd.ObjId

	objType := cmd.ObjType

	utcOffset := cmd.UtcOffset

	metadata := cmd.Metadata

	filter := buildFilter(dateMin, dateMax, objEvent, objId, objType, utcOffset)

	mask := "mask[eventName,label,objectName,eventCreateDate,userId,userType,objectId,metaData,user[username]]"
	logs, err := cmd.EventLogManager.GetEventLogs(mask, filter, limit)
	if err != nil {
		return errors.NewAPIError(T("Failed to get Event Logs.\n"), err.Error(), 2)
	}
	if logs[0].EventName != nil {
		var table terminal.Table
		if metadata {
			table = cmd.UI.Table([]string{T("Event"), T("Object"), T("Type"), T("Date"), T("Username"), T("Metadata")})
		} else {
			table = cmd.UI.Table([]string{T("Event"), T("Object"), T("Type"), T("Date"), T("Username")})
		}

		for _, log := range logs {
			user := "CUSTOMER"
			if log.UserId != nil && log.User != nil && log.User.Username != nil {
				user = utils.FormatStringPointer(log.User.Username)
			} else if log.UserType != nil {
				user = utils.FormatStringPointer(log.UserType)
			}

			if metadata {
				table.Add(
					utils.FormatStringPointer(log.EventName),
					utils.FormatStringPointer(log.Label),
					utils.FormatStringPointer(log.ObjectName),
					utils.FormatSLTimePointer(log.EventCreateDate),
					user,
					utils.FormatStringPointer(log.MetaData),
				)
			} else {
				table.Add(
					utils.FormatStringPointer(log.EventName),
					utils.FormatStringPointer(log.Label),
					utils.FormatStringPointer(log.ObjectName),
					utils.FormatSLTimePointer(log.EventCreateDate),
					user,
				)
			}
		}

		utils.PrintTable(cmd.UI, table, outputFormat)
	} else {
		cmd.UI.Print(T("No logs available for filter {{.filter}}", map[string]interface{}{"filter": filter}))
	}

	return nil
}

func buildFilter(dateMin string, dateMax string, objEvent string, objId string, objType string, utcOffset string) string {
	filters := filter.New()
	if dateMin != "" && dateMax != "" {
		filters = append(filters, filter.Path("eventCreateDate").DateBetween(setUtcOffset(dateMin, utcOffset), setUtcOffset(dateMax, utcOffset)))
	} else if dateMin != "" {
		filters = append(filters, filter.Path("eventCreateDate").DateAfter(setUtcOffset(dateMin, utcOffset)))
	} else if dateMax != "" {
		filters = append(filters, filter.Path("eventCreateDate").DateBefore(setUtcOffset(dateMax, utcOffset)))
	}

	if objEvent != "" {
		filters = append(filters, filter.Path("eventName").Eq(objEvent))
	}

	if objId != "" {
		filters = append(filters, filter.Path("objectId").Eq(objId))
	}

	if objType != "" {
		filters = append(filters, filter.Path("objectName").Eq(objType))
	}

	return filters.Build()
}

func setUtcOffset(date string, utcOffset string) string {
	if utcOffset == "" {
		date = date + ".000000-00:00"
	} else {
		date = date + ".000000" + utcOffset[:3] + ":" + utcOffset[3:]
	}
	return date
}

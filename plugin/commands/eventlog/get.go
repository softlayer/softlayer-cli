package eventlog

import (
	"fmt"
	"strconv"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type GetCommand struct {
	UI              terminal.UI
	EventLogManager managers.EventLogManager
}

func NewGetCommand(ui terminal.UI, eventLogManagerManager managers.EventLogManager) (cmd *GetCommand) {
	return &GetCommand{
		UI:              ui,
		EventLogManager: eventLogManagerManager,
	}
}

func (cmd *GetCommand) Run(c *cli.Context) error {

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	limit := 50
	if c.IsSet("limit") {
		newLimit, err := strconv.Atoi(c.String("limit"))
		if err != nil {
			return errors.NewInvalidSoftlayerIdInputError("limit")
		}
		limit = newLimit
	}

	dateMin := ""
	if c.IsSet("date-min") {
		date := c.String("date-min")
		time, err := time.Parse(time.RFC3339, date+"T00:00:00Z")
		if err != nil {
			return errors.NewInvalidUsageError(T("Invalid format date to --date-min."))
		}
		dateMin = time.Format("2006-01-02T15:04:05")
	}

	dateMax := ""
	if c.IsSet("date-max") {
		date := c.String("date-max")
		time, err := time.Parse(time.RFC3339, date+"T00:00:00Z")
		if err != nil {
			return errors.NewInvalidUsageError(T("Invalid format date to --date-max."))
		}
		dateMax = time.Format("2006-01-02T15:04:05")
	}

	objEvent := ""
	if c.IsSet("obj-event") {
		objEvent = c.String("obj-event")
	}

	objId := ""
	if c.IsSet("obj-id") {
		objId = c.String("obj-id")
	}

	objType := ""
	if c.IsSet("obj-type") {
		objType = c.String("obj-type")
	}

	utcOffset := ""
	if c.IsSet("utc-offset") {
		utcOffset = c.String("utc-offset")
	}

	metadata := false
	if c.IsSet("metadata") {
		metadata = c.Bool("metadata")
	}

	filter := buildFilter(dateMin, dateMax, objEvent, objId, objType, utcOffset)

	mask := "mask[eventName,label,objectName,eventCreateDate,userId,userType,objectId,metaData,user[username]]"
	logs, err := cmd.EventLogManager.GetEventLogs(mask, filter, limit)
	if err != nil {
		return cli.NewExitError(T("Failed to get Event Logs.\n")+err.Error(), 2)
	}
	if logs[0].EventName != nil {
		var table terminal.Table
		if metadata {
			table = cmd.UI.Table([]string{T("Event"), T("Object"), T("Type"), T("Date"), T("Username"), T("Metadata")})
		} else {
			table = cmd.UI.Table([]string{T("Event"), T("Object"), T("Type"), T("Date"), T("Username")})
		}

		for _, log := range logs {
			user := ""
			if log.UserId != nil {
				user = *log.User.Username
			} else {
				user = *log.UserType
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
		fmt.Println(filter)
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

func EventLogGetMetaData() cli.Command {
	return cli.Command{
		Category:    "event-log",
		Name:        "get",
		Description: T("Get Event Logs"),
		Usage: T(`${COMMAND_NAME} sl event-log get [OPTIONS]

EXAMPLE: 
   ${COMMAND_NAME} sl event-log get 
   ${COMMAND_NAME} sl event-log get --limit 5 --obj-id 123456 --obj-event Create --metadata
   ${COMMAND_NAME} sl event-log get --date-min 2021-03-31 --date-max 2021-04-31`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "d,date-min",
				Usage: T("The earliest date we want to search for event logs [YYYY-MM-DD]."),
			},
			cli.StringFlag{
				Name:  "D,date-max",
				Usage: T("The latest date we want to search for event logs [YYYY-MM-DD]."),
			},
			cli.StringFlag{
				Name:  "e,obj-event",
				Usage: T("The event we want to get event logs for"),
			},
			cli.StringFlag{
				Name:  "i,obj-id",
				Usage: T("The id of the object we want to get event logs for"),
			},
			cli.StringFlag{
				Name:  "t,obj-type",
				Usage: T("The type of the object we want to get event logs for"),
			},
			cli.StringFlag{
				Name:  "z,utc-offset",
				Usage: T("UTC Offset for searching with dates. +/-HHMM format  [default: -0000]"),
			},
			cli.BoolFlag{
				Name:  "metadata",
				Usage: T("Display metadata if present  [default: no-metadata]"),
			},
			cli.StringFlag{
				Name:  "l,limit",
				Usage: T("Total number of result to return. -1 to return ALL, there may be a LOT of these.  [default: 50]"),
			},
			metadata.OutputFlag(),
		},
	}
}

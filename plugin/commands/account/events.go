package account

import (
	"bytes"
	"time"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type EventsCommand struct {
	*metadata.SoftlayerCommand
	AccountManager managers.AccountManager
	Command        *cobra.Command
	DateMin        string
	Ack            bool
	Planned        bool
	Unplanned      bool
	Announcement   bool
}

func NewEventsCommand(sl *metadata.SoftlayerCommand) *EventsCommand {
	thisCmd := &EventsCommand{
		SoftlayerCommand: sl,
		AccountManager:   managers.NewAccountManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "events",
		Short: T("Summary and acknowledgement of upcoming and ongoing maintenance events."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.DateMin, "date-min", "d", "", T("Earliest date to retrieve events for [YYYY-MM-DD]. Default: 2 days ago."))
	cobraCmd.Flags().BoolVar(&thisCmd.Ack, "ack-all", false, T("Acknowledge every upcoming event. Doing so will turn off the popup in the control portal."))
	cobraCmd.Flags().BoolVar(&thisCmd.Planned, "planned", false, T("Show only planned events."))
	cobraCmd.Flags().BoolVar(&thisCmd.Unplanned, "unplanned", false, T("Show only unplanned events."))
	cobraCmd.Flags().BoolVar(&thisCmd.Announcement, "announcement", false, T("Show only announcement events."))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *EventsCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	datefilter := ""
	if cmd.DateMin != "" {
		date := cmd.DateMin
		time, err := time.Parse(time.RFC3339, date+"T00:00:00Z")
		if err != nil {
			return errors.NewInvalidUsageError(T("Invalid format date."))
		}
		datefilter = time.Format("01/02/2006")
	} else {
		time := time.Now()
		time = time.AddDate(0, 0, -2) // rest 2 days
		datefilter = time.Format("01/02/2006")
	}

	mask := "mask[id, subject, startDate, endDate, modifyDate, statusCode, acknowledgedFlag, impactedResourceCount, updateCount, systemTicketId, notificationOccurrenceEventType[keyName]]"

	// Gets three specific types of events
	plannedEvents, err := cmd.AccountManager.GetEvents("PLANNED", mask, datefilter)
	if err != nil {
		return errors.NewAPIError(T("Failed to get planned events."), err.Error(), 2)
	}

	unplannedEvents, err := cmd.AccountManager.GetEvents("UNPLANNED_INCIDENT", mask, datefilter)
	if err != nil {
		return errors.NewAPIError(T("Failed to get unplanned events."), err.Error(), 2)
	}

	announcement, err := cmd.AccountManager.GetEvents("ANNOUNCEMENT", mask, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get announcement events."), err.Error(), 2)
	}

	if cmd.Ack {
		ackAll(plannedEvents, cmd.AccountManager)
		ackAll(unplannedEvents, cmd.AccountManager)
		ackAll(announcement, cmd.AccountManager)
	}

	if cmd.Planned {
		PrintPlannedEvents(plannedEvents, cmd.UI, outputFormat)
	}
	if cmd.Unplanned {
		PrintUnplannedEvents(unplannedEvents, cmd.UI, outputFormat)
	}
	if cmd.Announcement {
		PrintAnnouncementEvents(announcement, cmd.UI, outputFormat)
	}

	if !cmd.Planned && !cmd.Unplanned && !cmd.Announcement {
		PrintPlannedEvents(plannedEvents, cmd.UI, outputFormat)
		PrintUnplannedEvents(unplannedEvents, cmd.UI, outputFormat)
		PrintAnnouncementEvents(announcement, cmd.UI, outputFormat)
	}

	return nil
}

func PrintPlannedEvents(events []datatypes.Notification_Occurrence_Event, ui terminal.UI, outputFormat string) {
	bufEvent := new(bytes.Buffer)
	table := terminal.NewTable(bufEvent, []string{
		T("Event Data"),
		T("Id"),
		T("Event ID"),
		T("Subject"),
		T("Status"),
		T("Items"),
		T("Start Date"),
		T("End Date"),
		T("Acknowledged"),
		T("Updates"),
	})
	for _, event := range events {
		table.Add(
			utils.FormatSLTimePointer(event.StartDate),
			utils.FormatIntPointer(event.Id),
			utils.FormatIntPointer(event.SystemTicketId),
			utils.ShortenString(utils.FormatStringPointer(event.Subject)),
			utils.FormatStringPointer(event.StatusCode.Name),
			utils.FormatUIntPointer(event.ImpactedResourceCount),
			utils.FormatSLTimePointer(event.StartDate),
			utils.FormatSLTimePointer(event.EndDate),
			utils.FormatBoolPointer(event.AcknowledgedFlag),
			utils.FormatUIntPointer(event.UpdateCount),
		)
	}
	utils.PrintTableWithTitle(ui, table, bufEvent, "Planned", outputFormat)
}

func PrintUnplannedEvents(events []datatypes.Notification_Occurrence_Event, ui terminal.UI, outputFormat string) {
	bufEvent := new(bytes.Buffer)
	table := terminal.NewTable(bufEvent, []string{
		T("Id"),
		T("Event ID"),
		T("Subject"),
		T("Status"),
		T("Items"),
		T("Start Date"),
		T("Last Updated"),
		T("Acknowledged"),
		T("Updates"),
	})
	for _, event := range events {
		table.Add(
			utils.FormatIntPointer(event.Id),
			utils.FormatIntPointer(event.SystemTicketId),
			utils.ShortenString(utils.FormatStringPointer(event.Subject)),
			utils.FormatStringPointer(event.StatusCode.Name),
			utils.FormatUIntPointer(event.ImpactedResourceCount),
			utils.FormatSLTimePointer(event.StartDate),
			utils.FormatSLTimePointer(event.ModifyDate),
			utils.FormatBoolPointer(event.AcknowledgedFlag),
			utils.FormatUIntPointer(event.UpdateCount),
		)
	}
	utils.PrintTableWithTitle(ui, table, bufEvent, "Unplanned", outputFormat)
}

func PrintAnnouncementEvents(events []datatypes.Notification_Occurrence_Event, ui terminal.UI, outputFormat string) {
	bufEvent := new(bytes.Buffer)
	table := terminal.NewTable(bufEvent, []string{
		T("Id"),
		T("Event ID"),
		T("Subject"),
		T("Status"),
		T("Items"),
		T("Acknowledged"),
		T("Updates"),
	})
	for _, event := range events {
		table.Add(
			utils.FormatIntPointer(event.Id),
			utils.FormatIntPointer(event.SystemTicketId),
			utils.ShortenString(utils.FormatStringPointer(event.Subject)),
			utils.FormatStringPointer(event.StatusCode.Name),
			utils.FormatUIntPointer(event.ImpactedResourceCount),
			utils.FormatBoolPointer(event.AcknowledgedFlag),
			utils.FormatUIntPointer(event.UpdateCount),
		)
	}
	utils.PrintTableWithTitle(ui, table, bufEvent, "Announcement", outputFormat)
}

func ackAll(events []datatypes.Notification_Occurrence_Event, accountManager managers.AccountManager) {
	for _, event := range events {
		if event.Id != nil {
			accountManager.AckEvent(*event.Id)
		}
	}
}

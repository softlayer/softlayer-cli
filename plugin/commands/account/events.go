package account

import (
	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type EventsCommand struct {
	UI             terminal.UI
	AccountManager managers.AccountManager
}

func NewEventsCommand(ui terminal.UI, accountManager managers.AccountManager) (cmd *EventsCommand) {
	return &EventsCommand{
		UI:             ui,
		AccountManager: accountManager,
	}
}

func EventsMetaData() cli.Command {
	return cli.Command{
		Category:    "account",
		Name:        "events",
		Description: T("Summary and acknowledgement of upcoming and ongoing maintenance events"),
		Usage:       T(`${COMMAND_NAME} sl account events`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

func (cmd *EventsCommand) Run(c *cli.Context) error {
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	// Gets three specific types of events
	plannedEvents, err := cmd.AccountManager.GetEvents("PLANNED")
	if err != nil {
		return cli.NewExitError(T("Failed to get planned events.")+err.Error(), 2)
	}
	unplannedEvents, err := cmd.AccountManager.GetEvents("UNPLANNED_INCIDENT")
	if err != nil {
		return cli.NewExitError(T("Failed to get unplanned events.")+err.Error(), 2)
	}
	announcement, err := cmd.AccountManager.GetEvents("ANNOUNCEMENT")
	if err != nil {
		return cli.NewExitError(T("Failed to get announcement events.")+err.Error(), 2)
	}

	// Print All events with keyname specific: PLANNED, UNPLANNED AND ANNOUNCEMENT
	PrintPlannedEvents(plannedEvents, cmd.UI, outputFormat)
	PrintUnplannedEvents(unplannedEvents, cmd.UI, outputFormat)
	PrintAnnouncementEvents(announcement, cmd.UI, outputFormat)
	return nil
}

func PrintPlannedEvents(events []datatypes.Notification_Occurrence_Event, ui terminal.UI, outputFormat string) {
	table := ui.Table([]string{
		"Event Data",
		"Id",
		"Event ID",
		"Subject",
		"Status",
		"Items",
		"Start Date",
		"End Date",
		"Acknowledged",
		"Updates",
	})
	for _, event := range events {
		table.Add(
			utils.FormatSLTimePointer(event.StartDate),
			utils.FormatIntPointer(event.Id),
			utils.FormatIntPointer(event.SystemTicketId),
			utils.FormatStringPointer(event.Subject),
			utils.FormatStringPointer(event.StatusCode.Name),
			utils.FormatUIntPointer(event.ImpactedResourceCount),
			utils.FormatSLTimePointer(event.StartDate),
			utils.FormatSLTimePointer(event.EndDate),
			utils.FormatBoolPointer(event.AcknowledgedFlag),
			utils.FormatUIntPointer(event.UpdateCount),
		)
	}
	utils.PrintTableWithTitle(ui, table, "Planned", outputFormat)
}

func PrintUnplannedEvents(events []datatypes.Notification_Occurrence_Event, ui terminal.UI, outputFormat string) {
	table := ui.Table([]string{
		"Id",
		"Event ID",
		"Subject",
		"Status",
		"Items",
		"Start Date",
		"Last Updated",
		"Acknowledged",
		"Updates",
	})
	for _, event := range events {
		table.Add(
			utils.FormatIntPointer(event.Id),
			utils.FormatIntPointer(event.SystemTicketId),
			utils.FormatStringPointer(event.Subject),
			utils.FormatStringPointer(event.StatusCode.Name),
			utils.FormatUIntPointer(event.ImpactedResourceCount),
			utils.FormatSLTimePointer(event.StartDate),
			utils.FormatSLTimePointer(event.ModifyDate),
			utils.FormatBoolPointer(event.AcknowledgedFlag),
			utils.FormatUIntPointer(event.UpdateCount),
		)
	}
	utils.PrintTableWithTitle(ui, table, "Unplanned", outputFormat)
}

func PrintAnnouncementEvents(events []datatypes.Notification_Occurrence_Event, ui terminal.UI, outputFormat string) {
	table := ui.Table([]string{
		"Id",
		"Event ID",
		"Subject",
		"Status",
		"Items",
		"Acknowledged",
		"Updates",
	})
	for _, event := range events {
		table.Add(
			utils.FormatIntPointer(event.Id),
			utils.FormatIntPointer(event.SystemTicketId),
			utils.FormatStringPointer(event.Subject),
			utils.FormatStringPointer(event.StatusCode.Name),
			utils.FormatUIntPointer(event.ImpactedResourceCount),
			utils.FormatBoolPointer(event.AcknowledgedFlag),
			utils.FormatUIntPointer(event.UpdateCount),
		)
	}
	utils.PrintTableWithTitle(ui, table, "Announcement", outputFormat)
}

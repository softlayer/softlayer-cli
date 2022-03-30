package account

import (
	"bytes"

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

	plannedEvents, err := cmd.AccountManager.GetEvents("PLANNED")
	unplannedEvents, err := cmd.AccountManager.GetEvents("UNPLANNED_INCIDENT")
	announcement, err := cmd.AccountManager.GetEvents("ANNOUNCEMENT")
	if err != nil {
		return cli.NewExitError(T("Failed to get events.")+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		allEvents := []datatypes.Notification_Occurrence_Event{}
		allEvents = append(allEvents, plannedEvents...)
		allEvents = append(allEvents, unplannedEvents...)
		allEvents = append(allEvents, announcement...)
		return utils.PrintPrettyJSON(cmd.UI, allEvents)
	}

	PrintPlannedEvents(plannedEvents, cmd.UI)
	PrintUnplannedEvents(unplannedEvents, cmd.UI)
	PrintAnnouncement(announcement, cmd.UI)
	return nil
}

func PrintPlannedEvents(events []datatypes.Notification_Occurrence_Event, ui terminal.UI) {
	tableTitle := ui.Table([]string{T("Planned")})
	bufEvent := new(bytes.Buffer)
	table := terminal.NewTable(bufEvent, []string{
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
	table.Print()
	tableTitle.Add(bufEvent.String())
	tableTitle.Print()
}

func PrintUnplannedEvents(events []datatypes.Notification_Occurrence_Event, ui terminal.UI) {
	tableTitle := ui.Table([]string{T("Unplanned")})
	bufEvent := new(bytes.Buffer)
	table := terminal.NewTable(bufEvent, []string{
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
	table.Print()
	tableTitle.Add(bufEvent.String())
	tableTitle.Print()
}

func PrintAnnouncement(events []datatypes.Notification_Occurrence_Event, ui terminal.UI) {
	tableTitle := ui.Table([]string{T("Announcement")})
	bufEvent := new(bytes.Buffer)
	table := terminal.NewTable(bufEvent, []string{
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
	table.Print()
	tableTitle.Add(bufEvent.String())
	tableTitle.Print()
}

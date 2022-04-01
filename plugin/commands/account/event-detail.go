package account

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type EventDetailCommand struct {
	UI             terminal.UI
	AccountManager managers.AccountManager
}

func NewEventDetailCommand(ui terminal.UI, accountManager managers.AccountManager) (cmd *EventDetailCommand) {
	return &EventDetailCommand{
		UI:             ui,
		AccountManager: accountManager,
	}
}

func EventDetailMetaData() cli.Command {
	return cli.Command{
		Category:    "account",
		Name:        "event-detail",
		Description: T("Details of a specific event, and ability to acknowledge event."),
		Usage:       T(`${COMMAND_NAME} sl account event-detail`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

func (cmd *EventDetailCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError("Event ID is required.")
	}

	eventID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return errors.NewInvalidUsageError(T("The event ID has to be a positive integer."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	event, err := cmd.AccountManager.GetEventDetail(eventID)
	if err != nil {
		return cli.NewExitError(T("Failed to get the event {{.eventID}}. ", map[string]interface{}{"eventID": eventID})+err.Error(), 2)
	}

	BasicEventTable(event, cmd.UI, outputFormat)
	ImpactedTable(event, cmd.UI, outputFormat)
	UpdateTable(event, cmd.UI, outputFormat)
	return nil
}

func BasicEventTable(event datatypes.Notification_Occurrence_Event, ui terminal.UI, outputFormat string) {
	tableTitle := ui.Table([]string{T(utils.FormatStringPointer(event.Subject))})
	bufEvent := new(bytes.Buffer)
	table := terminal.NewTable(bufEvent, []string{
		"Id",
		"Status",
		"Type",
		"Start",
		"End",
	})
	table.Add(
		utils.FormatIntPointer(event.Id),
		utils.FormatStringPointer(event.StatusCode.Name),
		utils.FormatStringPointer(event.NotificationOccurrenceEventType.KeyName),
		utils.FormatSLTimePointer(event.StartDate),
		utils.FormatSLTimePointer(event.EndDate),
	)
	if outputFormat == "JSON" {
		table.PrintJson()
		tableTitle.Add(bufEvent.String())
		tableTitle.PrintJson()
	} else {
		table.Print()
		tableTitle.Add(bufEvent.String())
		tableTitle.Print()
	}
}

func ImpactedTable(event datatypes.Notification_Occurrence_Event, ui terminal.UI, outputFormat string) {
	table := ui.Table([]string{
		T("Id"),
		T("Hostname"),
		T("Label"),
	})
	for _, resources := range event.ImpactedResources {
		table.Add(
			utils.FormatIntPointer(resources.ResourceTableId),
			utils.FormatStringPointer(resources.ResourceName),
			utils.FormatStringPointer(resources.FilterLabel),
		)
	}
	if outputFormat == "JSON" {
		table.PrintJson()
	} else {
		table.Print()
	}
}

func UpdateTable(event datatypes.Notification_Occurrence_Event, ui terminal.UI, outputFormat string) {
	updateStartDate := ""
	text := ""
	for _, update := range event.Updates {
		updateStartDate = utils.FormatSLTimePointer(update.StartDate)
		text += utils.FormatStringPointer(update.Contents)
	}
	header := fmt.Sprintf("======= Update #%d on %s =======", len(event.Updates), updateStartDate)

	ui.Print(T(header))
	ui.Print((text))
}

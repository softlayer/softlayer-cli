package account

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type EventDetailCommand struct {
	*metadata.SoftlayerCommand
	AccountManager managers.AccountManager
	Command *cobra.Command
}

func NewEventDetailCommand(sl *metadata.SoftlayerCommand) *EventDetailCommand {
	thisCmd := &EventDetailCommand{
		SoftlayerCommand: sl,
		AccountManager: managers.NewAccountManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use: "event-detail " + T("IDENTIFIER"),
		Short: T("Details of a specific event, and ability to acknowledge event."),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}


func (cmd *EventDetailCommand) Run(args []string) error {

	eventID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError(T("Event ID"))
	}

	outputFormat := cmd.GetOutputFlag()

	mask := "mask[acknowledgedFlag,attachments,impactedResources,statusCode,updates,notificationOccurrenceEventType]"
	event, err := cmd.AccountManager.GetEventDetail(eventID, mask)
	if err != nil {
		sub := map[string]interface{}{"eventID": eventID}
		return slErr.NewAPIError(T("Failed to get the event {{.eventID}}. ", sub), err.Error(), 2)
	}

	BasicEventTable(event, cmd.UI, outputFormat)
	ImpactedTable(event, cmd.UI, outputFormat)
	UpdateTable(event, cmd.UI, outputFormat)
	return nil
}

func BasicEventTable(event datatypes.Notification_Occurrence_Event, ui terminal.UI, outputFormat string) {
	bufEvent := new(bytes.Buffer)
	table := terminal.NewTable(bufEvent, []string{
		T("Id"),
		T("Status"),
		T("Type"),
		T("Start"),
		T("End"),
	})
	table.Add(
		utils.FormatIntPointer(event.Id),
		utils.FormatStringPointer(event.StatusCode.Name),
		utils.FormatStringPointer(event.NotificationOccurrenceEventType.KeyName),
		utils.FormatSLTimePointer(event.StartDate),
		utils.FormatSLTimePointer(event.EndDate),
	)
	utils.PrintTableWithTitle(ui, table, bufEvent, utils.FormatStringPointer(event.Subject), outputFormat)
	
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
	utils.PrintTable(ui, table, outputFormat)
}

func UpdateTable(event datatypes.Notification_Occurrence_Event, ui terminal.UI, outputFormat string) {
	updateStartDate := ""
	text := ""
	for _, update := range event.Updates {
		updateStartDate = utils.FormatSLTimePointer(update.StartDate)
		text += utils.FormatStringPointerName(update.Contents)
	}
	header := fmt.Sprintf("======= Update #%d on %s =======", len(event.Updates), updateStartDate)

	if outputFormat == "JSON" {
		table := ui.Table([]string{
			T("Updates"),
		})
		table.Add(header)
		table.Add(text)
		table.PrintJson()
	} else {
		ui.Print(header)
		ui.Print(text)
	}
}

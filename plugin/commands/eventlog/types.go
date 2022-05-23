package eventlog

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type TypesCommand struct {
	UI              terminal.UI
	EventLogManager managers.EventLogManager
}

func NewTypesCommand(ui terminal.UI, eventLogManagerManager managers.EventLogManager) (cmd *TypesCommand) {
	return &TypesCommand{
		UI:              ui,
		EventLogManager: eventLogManagerManager,
	}
}

func (cmd *TypesCommand) Run(c *cli.Context) error {

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	types, err := cmd.EventLogManager.GetEventLogTypes()
	if err != nil {
		return cli.NewExitError(T("Failed to get Event Log types.\n")+err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Types")})

	for _, typeEvent := range types {
		table.Add(typeEvent)
	}

	utils.PrintTable(cmd.UI, table, outputFormat)

	return nil
}

func EventLogTypesMetaData() cli.Command {
	return cli.Command{
		Category:    "event-log",
		Name:        "types",
		Description: T("Get Event Log types"),
		Usage: T(`${COMMAND_NAME} sl event-log types

EXAMPLE: 
   ${COMMAND_NAME} sl event-log types`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

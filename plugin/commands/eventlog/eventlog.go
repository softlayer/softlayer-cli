package eventlog

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	eventLogManager := managers.NewEventLogManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"event-log-get": func(c *cli.Context) error {
			return NewGetCommand(ui, eventLogManager).Run(c)
		},
		"event-log-types": func(c *cli.Context) error {
			return NewTypesCommand(ui, eventLogManager).Run(c)
		},
	}

	return CommandActionBindings
}

func EventLogNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "event-log",
		Description: T("Classic infrastructure Event Log Group"),
	}
}

func EventLogMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "event-log",
		Description: T("Classic infrastructure Event Log Group"),
		Usage:       "${COMMAND_NAME} sl event-log",
		Subcommands: []cli.Command{
			EventLogGetMetaData(),
			EventLogTypesMetaData(),
		},
	}
}

package reports

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

/*
This account package follows a slightly different pattern than the other CLI commands
because I'd like to eventually adpot this pattern throughout to get away from having metadata files
for every command.
*/

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	CommandActionBindings := map[string]func(c *cli.Context) error{
		"report-datacenter-closures": func(c *cli.Context) error {
			return NewDCClosuresCommand(ui, session).Run(c)
		},
	}

	return CommandActionBindings
}

func ReportsNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "report",
		Description: T("Classic Infrastructure Reports"),
	}
}

func ReportsMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "report",
		Description: T("Classic Infrastructure Reports"),
		Usage:       "${COMMAND_NAME} sl reports",
		Subcommands: []cli.Command{
			DCClosuresMetaData(),
		},
	}
}

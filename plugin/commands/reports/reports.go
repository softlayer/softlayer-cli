package reports

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

/*
This account package follows a slightly different pattern than the other CLI commands
because I'd like to eventually adpot this pattern throughout to get away from having metadata files
for every command.
*/

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	reportManager := managers.NewReportManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"report-datacenter-closures": func(c *cli.Context) error {
			return NewDCClosuresCommand(ui, session).Run(c)
		},
		"report-bandwidth": func(c *cli.Context) error {
			return NewBandwidthCommand(ui, reportManager).Run(c)
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
		Usage:       "${COMMAND_NAME} sl report",
		Subcommands: []cli.Command{
			DCClosuresMetaData(),
			ReportBandwidthMetaData(),
		},
	}
}

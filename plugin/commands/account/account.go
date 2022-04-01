package account

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
	accountManager := managers.NewAccountManager(session)
	CommandActionBindings := map[string]func(c *cli.Context) error{
		"account-bandwidth-pools": func(c *cli.Context) error {
			return NewBandwidthPoolsCommand(ui, session).Run(c)
		},
		"account-event-detail": func(c *cli.Context) error {
			return NewEventDetailCommand(ui, accountManager).Run(c)
		},
	}

	return CommandActionBindings
}

func AccountNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "account",
		Description: T("Classic infrastructure Account commands"),
	}
}

func AccountMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "account",
		Description: T("Classic infrastructure Account commands"),
		Usage:       "${COMMAND_NAME} sl account",
		Subcommands: []cli.Command{
			BandwidthPoolsMetaData(),
			EventDetailMetaData(),
		},
	}
}

package nas

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	nasNetworkStorageManager := managers.NewNasNetworkStorageManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"nas-list": func(c *cli.Context) error {
			return NewListCommand(ui, nasNetworkStorageManager).Run(c)
		},
		"nas-credentials": func(c *cli.Context) error {
			return NewCredentialsCommand(ui, nasNetworkStorageManager).Run(c)
		},
	}

	return CommandActionBindings
}

func NasNetworkStorageNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "nas",
		Description: T("Classic infrastructure Network Attached Storage"),
	}
}

func NasNetworkStorageMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "nas",
		Description: T("Classic infrastructure Network Attached Storage"),
		Usage:       "${COMMAND_NAME} sl nas",
		Subcommands: []cli.Command{
			NasListMetaData(),
			NasCredentialsMetaData(),
		},
	}
}

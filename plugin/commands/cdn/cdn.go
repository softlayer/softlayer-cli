package cdn

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	cdnManager := managers.NewCdnManager(session)
	CommandActionBindings := map[string]func(c *cli.Context) error{
		"cdn-list": func(c *cli.Context) error {
			return NewListCommand(ui, cdnManager).Run(c)
		},
		"cdn-detail": func(c *cli.Context) error {
			return NewDetailCommand(ui, cdnManager).Run(c)
		},
		"cdn-edit": func(c *cli.Context) error {
			return NewEditCommand(ui, cdnManager).Run(c)
		},
	}

	return CommandActionBindings
}

func CdnNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "cdn",
		Description: T("Classic infrastructure CDN commands"),
	}
}

func CdnMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "cdn",
		Description: T("Classic infrastructure CDN commands"),
		Usage:       "${COMMAND_NAME} sl cdn",
		Subcommands: []cli.Command{
			ListMetaData(),
			DetailMetaData(),
			EditMetaData(),
		},
	}
}

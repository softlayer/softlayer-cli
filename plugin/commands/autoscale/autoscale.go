package autoscale

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	autoScaleManager := managers.NewAutoScaleManager(session)
	virtualServerManager := managers.NewVirtualServerManager(session)
	securityManager := managers.NewSecurityManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"autoscale-edit": func(c *cli.Context) error {
			return NewEditCommand(ui, autoScaleManager).Run(c)
		},
		"autoscale-tag": func(c *cli.Context) error {
			return NewTagCommand(ui, autoScaleManager, virtualServerManager).Run(c)
		},
		"autoscale-logs": func(c *cli.Context) error {
			return NewLogsCommand(ui, autoScaleManager, securityManager).Run(c)
		},
		"autoscale-detail": func(c *cli.Context) error {
			return NewDetailCommand(ui, autoScaleManager, securityManager).Run(c)
		},
		"autoscale-list": func(c *cli.Context) error {
			return NewListCommand(ui, autoScaleManager).Run(c)
		},
		"autoscale-scale": func(c *cli.Context) error {
			return NewScaleCommand(ui, autoScaleManager).Run(c)
		},
		"autoscale-delete": func(c *cli.Context) error {
			return NewDeleteCommand(ui, autoScaleManager).Run(c)
		},
	}

	return CommandActionBindings
}

func AutoScaleNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "autoscale",
		Description: T("Classic infrastructure Autoscale Group"),
	}
}

func AutoScaleMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "autoscale",
		Description: T("Classic infrastructure Autoscale Group"),
		Usage:       "${COMMAND_NAME} sl autoscale",
		Subcommands: []cli.Command{
			AutoScaleEditMetaData(),
			AutoScaleTagMetaData(),
			AutoScaleLogsMetaData(),
			AutoScaleDetailMetaData(),
			AutoScaleListMetaData(),
			AutoScaleScaleMetaData(),
			AutoScaleDeleteMetaData(),
		},
	}
}

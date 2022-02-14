package placementgroup

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	placeGroupManager := managers.NewPlaceGroupManager(session)
	virtualServerManager := managers.NewVirtualServerManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"placement-group-create": func(c *cli.Context) error {
			return NewPlacementGroupCreateCommand(ui, placeGroupManager).Run(c)
		},
		"placement-group-list": func(c *cli.Context) error {
			return NewPlacementGroupListCommand(ui, placeGroupManager).Run(c)
		},
		"placement-group-delete": func(c *cli.Context) error {
			return NewPlacementGroupDeleteCommand(ui, placeGroupManager, virtualServerManager).Run(c)
		},
		"placement-group-create-options": func(c *cli.Context) error {
			return NewPlacementGroupCreateOptionsCommand(ui, placeGroupManager).Run(c)
		},
		"placement-group-detail": func(c *cli.Context) error {
			return NewPlacementGroupDetailCommand(ui, placeGroupManager).Run(c)
		},
	}

	return CommandActionBindings
}

func PlacementGroupNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "placement-group",
		Description: T("Classic infrastructure Placement Group"),
	}
}

func PlacementGroupMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "placement-group",
		Description: T("Classic infrastructure Placement Group"),
		Usage:       "${COMMAND_NAME} sl placement-group",
		Subcommands: []cli.Command{
			PlacementGroupCreateMetaData(),
			PlacementGroupCreateOptionsMetaData(),
			PlacementGroupListMetaData(),
			PlacementGroupDeleteMetaData(),
			PlacementGroupDetailMetaData(),
		},
	}
}

package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	NS_PLACEMENT_GROUP_NAME  = "placement-group"
	CMD_PLACEMENT_GROUP_NAME = "placement-group"

	//sl security
	CMD_PLACEMENT_GROUP_CREATE_NAME         = "create"
	CMD_PLACEMENT_GROUP_CREATE_OPTIONS_NAME = "create-options"
	CMD_PLACEMENT_GROUP_DELETE_NAME         = "delete"
	CMD_PLACEMENT_GROUP_DETAIL_NAME         = "detail"
	CMD_PLACEMENT_GROUP_LIST_NAME           = "list"
)

func PlacementGroupNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_PLACEMENT_GROUP_NAME,
		Description: T("Classic infrastructure Placement Group"),
	}
}

func PlacementGroupMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_PLACEMENT_GROUP_NAME,
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

func PlacementGroupCreateMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_PLACEMENT_GROUP_NAME,
		Name:        CMD_PLACEMENT_GROUP_CREATE_NAME,
		Description: T("Create a placement group"),
		Usage:       "${COMMAND_NAME} sl placement-group create (--name NAME) (-b, --backend-router-id BACKENDROUTER) (-r, --rule-id RULE) [--output FORMAT]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("Name for this new placement group. [required]"),
			},
			cli.IntFlag{
				Name:  "b,backend-router-id",
				Usage: T("Backend router ID. [required]"),
			},
			cli.IntFlag{
				Name:  "r,rule-id",
				Usage: T("The ID of the rule to govern this placement group. [required]"),
			},
			OutputFlag(),
		},
	}
}

func PlacementGroupCreateOptionsMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_PLACEMENT_GROUP_NAME,
		Name:        CMD_PLACEMENT_GROUP_CREATE_OPTIONS_NAME,
		Description: T("List options for creating a placement group"),
		Usage:       "${COMMAND_NAME} sl placement-group create-options",
		Flags:       []cli.Flag{},
	}
}

func PlacementGroupListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_PLACEMENT_GROUP_NAME,
		Name:        CMD_PLACEMENT_GROUP_LIST_NAME,
		Description: T("List placement groups"),
		Usage:       "${COMMAND_NAME} sl placement-group list [--output FORMAT]",
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

func PlacementGroupDeleteMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_PLACEMENT_GROUP_NAME,
		Name:        CMD_PLACEMENT_GROUP_DELETE_NAME,
		Description: T("Delete a placement group"),
		Usage:       "${COMMAND_NAME} sl placement-group delete (--id PLACEMENTGROUP_ID) [-f, --force]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "id",
				Usage: T("ID for the placement group. [required]"),
			},
			//cli.BoolFlag{   # tmp disable this option. because the placement can't be deleted if the VSI status is delete pending.
			//	Name:  "purge",
			//	Usage: T("Delete all guests in this placement group. The group itself can be deleted once all VMs are fully reclaimed"),
			//},
			ForceFlag(),
		},
	}
}

func PlacementGroupDetailMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_PLACEMENT_GROUP_NAME,
		Name:        CMD_PLACEMENT_GROUP_DETAIL_NAME,
		Description: T("View details of a placement group"),
		Usage:       "${COMMAND_NAME} sl placement-group detail (--id PLACEMENTGROUP_ID) [--output FORMAT]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "id",
				Usage: T("ID for the placement group. [required]"),
			},
			OutputFlag(),
		},
	}
}

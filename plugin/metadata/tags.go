package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	NS_TAGS_NAME         = "tags"
	CMD_TAGS_NAME        = "tags"
	CMD_TAGS_LIST_NAME   = "list"
	CMD_TAGS_DETAIL_NAME = "detail"
	CMD_TAGS_DELETE_NAME = "delete"
)

func TagsNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_TAGS_NAME,
		Description: T("Classic infrastructure Tag management"),
	}
}

// MAIN command
func TagsMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_TAGS_NAME,
		Description: T("Classic infrastructure Tag management"),
		Usage:       "${COMMAND_NAME} sl tags",
		Subcommands: []cli.Command{
			TagsListMetaData(),
			TagsDetailsMetaData(),
			TagsDeleteMetaData(),
		},
	}
}

// sl tags list
func TagsListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_TAGS_NAME,
		Name:        CMD_TAGS_LIST_NAME,
		Description: T("List all tags currently on your account"),
		Usage: T(`${COMMAND_NAME} sl tags list [OPTIONS]

EXAMPLE:
	${COMMAND_NAME} sl tags list
	Shows all tags and a count of devices associated with that tag.

	${COMMAND_NAME} sl tags list -d
	Shows all tags with devices, and some basic information about devices using this tag.
`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "d,detail",
				Usage: T("List information about devices using the tag."),
			},
			OutputFlag(),
		},
	}
}

// sl tags detail
func TagsDetailsMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_TAGS_NAME,
		Name:        CMD_TAGS_DETAIL_NAME,
		Description: T("Get information about the resources using the selected tag."),
		Usage: T(`${COMMAND_NAME} sl tags detail [TAG NAME]

EXAMPLE:
	${COMMAND_NAME} sl tags detail tag1
	Shows all items that are tagged with 'tag1'
`),
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

// sl tags delete
func TagsDeleteMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_TAGS_NAME,
		Name:        CMD_TAGS_DELETE_NAME,
		Description: T("Removes an empty tag from your account."),
		Usage: T(`${COMMAND_NAME} sl tags delete [TAG NAME]

EXAMPLE:
	${COMMAND_NAME} sl tags delete tag1
	Removes "tag" from your account.
`),
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

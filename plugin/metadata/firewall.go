package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
)

var (
	NS_FIREWALL_NAME = "firewall"
	NS_FIREWALL_DESC = T("Classic infrastructure Firewalls")

	CMD_FIREWALL_NAME  = "firewall"
	CMD_FIREWALL_DESC  = "Classic infrastructure Firewalls"
	CMD_FIREWALL_USAGE = "${COMMAND_NAME} sl firewall"
	//sl-firewall
	CMD_FW_ADD_NAME    = "add"
	CMD_FW_CANCEL_NAME = "cancel"
	CMD_FW_DETAIL_NAME = "detail"
	CMD_FW_EDIT_NAME   = "edit"
	CMD_FW_LIST_NAME   = "list"

	CMD_FW_ADD_DESC    = T("Create a new firewall")
	CMD_FW_CANCEL_DESC = T("Cancels a firewall")
	CMD_FW_DETAIL_DESC = T("Detail information about a firewall")
	CMD_FW_EDIT_DESC   = T("Edit firewall rules")
	CMD_FW_LIST_DESC   = T("List all firewalls on your account")

	CMD_FW_ADD_USAGE    = "${COMMAND_NAME} sl firewall add TARGET [OPTIONS]"
	CMD_FW_CANCEL_USAGE = "${COMMAND_NAME} sl firewall cancel IDENTIFIER [OPTIONS]"
	CMD_FW_DETAIL_USAGE = "${COMMAND_NAME} sl firewall detail  IDENTIFIER [OPTIONS]"
	CMD_FW_EDIT_USAGE   = "${COMMAND_NAME} sl firewall edit IDENTIFIER [OPTIONS]"
	CMD_FW_LIST_USAGE   = "${COMMAND_NAME} sl firewall list [OPTIONS]"

	CMD_FW_ADD_OPT1      = "type"
	CMD_FW_ADD_OPT1_DESC = T("Firewall type  [required]. Options are: vlan,vs,hardware")
	CMD_FW_ADD_OPT2      = "ha,high-availability"
	CMD_FW_ADD_OPT2_DESC = T("High available firewall option")
)

var NS_FIREWALL = plugin.Namespace{
	ParentName:  NS_SL_NAME,
	Name:        NS_FIREWALL_NAME,
	Description: NS_FIREWALL_DESC,
}

var CMD_FW = cli.Command{
	Category:    NS_SL_NAME,
	Name:        CMD_FIREWALL_NAME,
	Description: CMD_FIREWALL_DESC,
	Usage:       CMD_FIREWALL_USAGE,
	Subcommands: []cli.Command{
		CMD_FW_ADD,
		CMD_FW_CANCEL,
		CMD_FW_DETAIL,
		CMD_FW_EDIT,
		CMD_FW_LIST,
	},
}

var CMD_FW_ADD = cli.Command{
	Category:    CMD_FIREWALL_NAME,
	Name:        CMD_FW_ADD_NAME,
	Description: CMD_FW_ADD_DESC,
	Usage:       CMD_FW_ADD_USAGE,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  CMD_FW_ADD_OPT1,
			Usage: CMD_FW_ADD_OPT1_DESC,
		},
		cli.BoolFlag{
			Name:  CMD_FW_ADD_OPT2,
			Usage: CMD_FW_ADD_OPT2_DESC,
		},
		ForceFlag(),
	},
}

var CMD_FW_CANCEL = cli.Command{
	Category:    CMD_FIREWALL_NAME,
	Name:        CMD_FW_CANCEL_NAME,
	Description: CMD_FW_CANCEL_DESC,
	Usage:       CMD_FW_CANCEL_USAGE,
	Flags: []cli.Flag{
		ForceFlag(),
	},
}

var CMD_FW_DETAIL = cli.Command{
	Category:    CMD_FIREWALL_NAME,
	Name:        CMD_FW_DETAIL_NAME,
	Description: CMD_FW_DETAIL_DESC,
	Usage:       CMD_FW_DETAIL_USAGE,
}

var CMD_FW_EDIT = cli.Command{
	Category:    CMD_FIREWALL_NAME,
	Name:        CMD_FW_EDIT_NAME,
	Description: CMD_FW_EDIT_DESC,
	Usage:       CMD_FW_EDIT_USAGE,
}

var CMD_FW_LIST = cli.Command{
	Category:    CMD_FIREWALL_NAME,
	Name:        CMD_FW_LIST_NAME,
	Description: CMD_FW_LIST_DESC,
	Usage:       CMD_FW_LIST_USAGE,
}

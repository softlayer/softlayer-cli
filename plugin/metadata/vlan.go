package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	NS_VLAN_NAME  = "vlan"
	CMD_VLAN_NAME = "vlan"

	CMD_VLAN_CREATE_NAME  = "create"
	CMD_VLAN_CANCEL_NAME  = "cancel"
	CMD_VLAN_DETAIL_NAME  = "detail"
	CMD_VLAN_EDIT_NAME    = "edit"
	CMD_VLAN_LIST_NAME    = "list"
	CMD_VLAN_OPTIONS_NAME = "options"
)

func VlanNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_VLAN_NAME,
		Description: T("Classic infrastructure Network VLANs"),
	}
}

func VlanMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_VLAN_NAME,
		Description: T("Classic infrastructure Network VLANs"),
		Usage:       "${COMMAND_NAME} sl vlan",
		Subcommands: []cli.Command{
			VlanCreateMetaData(),
			VlanCancelMetaData(),
			VlanDetailMetaData(),
			VlanEditMetaData(),
			VlanListMetaData(),
			VlanOptionsMetaData(),
		},
	}
}

func VlanCreateMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VLAN_NAME,
		Name:        CMD_VLAN_CREATE_NAME,
		Description: T("Create a new VLAN"),
		Usage: T(`${COMMAND_NAME} sl vlan create [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vlan create -t public -d dal09 -n myvlan
   This command creates a public vlan located in datacenter dal09 named "myvlan".
   ${COMMAND_NAME} sl vlan create -r bcr01a.dal09 -n myvlan
   This command creates a vlan on router bcr01a.dal09 named "myvlan".`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "t,vlan-type",
				Usage: T("The type of the VLAN, either public or private"),
			},
			cli.StringFlag{
				Name:  "r,router",
				Usage: T("The hostname of the router"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("The short name of the datacenter"),
			},
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("The name of the VLAN"),
			},
			ForceFlag(),
			OutputFlag(),
		},
	}
}

func VlanCancelMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VLAN_NAME,
		Name:        CMD_VLAN_CANCEL_NAME,
		Description: T("Cancel a VLAN"),
		Usage: T(`${COMMAND_NAME} sl vlan cancel IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vlan cancel 12345678 -f
   This command cancels vlan with ID 12345678 without asking for confirmation.`),
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}

func VlanDetailMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VLAN_NAME,
		Name:        CMD_VLAN_DETAIL_NAME,
		Description: T("Get details about a VLAN"),
		Usage: T(`${COMMAND_NAME} sl vlan detail IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl vlan detail 12345678	--no-vs --no-hardware
   This command shows details of vlan with ID 12345678, and not list virtual server or hardware server.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "no-vs",
				Usage: T("Hide virtual server listing"),
			},
			cli.BoolFlag{
				Name:  "no-hardware",
				Usage: T("Hide hardware listing"),
			},
			OutputFlag(),
		},
	}
}

func VlanEditMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VLAN_NAME,
		Name:        CMD_VLAN_EDIT_NAME,
		Description: T("Edit the details about a VLAN"),
		Usage: T(`${COMMAND_NAME} sl vlan edit IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vlan edit 12345678 -n myvlan-rename
   This command updates vlan with ID 12345678 and gives it a new name "myvlan-rename".`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("The name of the VLAN"),
			},
		},
	}
}

func VlanListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VLAN_NAME,
		Name:        CMD_VLAN_LIST_NAME,
		Description: T("List all the VLANs on your account"),
		Usage: T(`${COMMAND_NAME} sl vlan list [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vlan list -d dal09 --sortby number
   This commands lists all vlans on current account filtering by datacenter equals to dal09, and sort them by vlan number.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by. Options are: id,number,name,firewall,datacenter,hardware,virtual_servers,public_ips"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Filter by datacenter shortname"),
			},
			cli.IntFlag{
				Name:  "n,number",
				Usage: T("Filter by VLAN number"),
			},
			cli.StringFlag{
				Name:  "name",
				Usage: T("Filter by VLAN name"),
			},
			cli.IntFlag{
				Name:  "order",
				Usage: T("Filter by ID of the order that purchased the VLAN"),
			},
			OutputFlag(),
		},
	}
}

func VlanOptionsMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VLAN_NAME,
		Name:        CMD_VLAN_OPTIONS_NAME,
		Description: T("List all the options for creating VLAN"),
		Usage: T(`${COMMAND_NAME} sl vlan options
	
EXAMPLE:
   ${COMMAND_NAME} sl vlan options
   This command lists all options for creating a vlan, eg. vlan type, datacenters, subnet size, routers, etc.`),
	}
}

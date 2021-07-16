package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	NS_SUBNET_NAME  = "subnet"
	CMD_SUBNET_NAME = "subnet"

	CMD_SUBNET_CANCEL_NAME = "cancel"
	CMD_SUBNET_CREATE_NAME = "create"
	CMD_SUBNET_DETAIL_NAME = "detail"
	CMD_SUBNET_LIST_NAME   = "list"
	CMD_SUBNET_LOOKUP_NAME = "lookup"
)

func SubnetNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_SUBNET_NAME,
		Description: T("Classic infrastructure Network subnets"),
	}
}

func SubnetMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_SUBNET_NAME,
		Description: T("Classic infrastructure Network subnets"),
		Usage:       "${COMMAND_NAME} sl subnet",
		Subcommands: []cli.Command{
			SubnetCancelMetaData(),
			SubnetCreateMetaData(),
			SubnetDetailMetaData(),
			SubnetListMetaData(),
			SubnetLookupMetaData(),
		},
	}
}

func SubnetCancelMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SUBNET_NAME,
		Name:        CMD_SUBNET_CANCEL_NAME,
		Description: T("Cancel a subnet"),
		Usage: T(`${COMMAND_NAME} sl subnet cancel IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl subnet cancel 12345678 -f
   This command cancels subnet with ID 12345678 without asking for confirmation.`),
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}

func SubnetCreateMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SUBNET_NAME,
		Name:        CMD_SUBNET_CREATE_NAME,
		Description: T("Add a new subnet to your account"),
		Usage: T(`${COMMAND_NAME} sl subnet create NETWORK QUANTITY VLAN_ID [OPTIONS]
	
	Add a new subnet to your account. Valid quantities vary by type.
	
	Type    - Valid Quantities (IPv4)
  	public  - 4, 8, 16, 32
  	private - 4, 8, 16, 32, 64

  	Type    - Valid Quantities (IPv6)
	public  - 64

EXAMPLE:
   ${COMMAND_NAME} sl subnet create public 16 567 
   This command creates a public subnet with 16 IPv4 addresses and places it on vlan with ID 567.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "v6,ipv6",
				Usage: T("Order IPv6 Addresses"),
			},
			cli.BoolFlag{
				Name:  "test",
				Usage: T("Do not order the subnet; just get a quote"),
			},
			ForceFlag(),
			OutputFlag(),
		},
	}
}

func SubnetDetailMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SUBNET_NAME,
		Name:        CMD_SUBNET_DETAIL_NAME,
		Description: T("Get details of a subnet"),
		Usage: T(`${COMMAND_NAME} sl subnet detail IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl subnet detail 12345678 
   This command shows detailed information about subnet with ID 12345678, including virtual servers and hardware servers information.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "no-vs",
				Usage: T("Hide virtual server listing"),
			},
			cli.BoolFlag{
				Name:  "no-hardware",
				Usage: T("Hide hardware listing"),
			},cli.BoolFlag{
				Name:  "no-ip-address",
				Usage: T("Hide IP address listing"),
			},cli.BoolFlag{
				Name:  "no-Tag",
				Usage: T("Hide Tag listing"),
			},
			OutputFlag(),
		},
	}
}

func SubnetListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SUBNET_NAME,
		Name:        CMD_SUBNET_LIST_NAME,
		Description: T("List all subnets on your account"),
		Usage: T(`${COMMAND_NAME} sl subnet list [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl subnet list -d dal09 -t PRIMARY --network-space PUBLIC --v4
   This command lists IPv4 subnets on the current account, and filters by datacenter is dal09, subnet type is PRIMARY, and network space is PUBLIC.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by. Options are: id,identifier,type,network_space,datacenter,vlan_id,IPs,hardware,vs"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Filter by datacenter shortname"),
			},
			cli.StringFlag{
				Name:  "identifier",
				Usage: T("Filter by network identifier"),
			},
			cli.StringFlag{
				Name:  "t,subnet-type",
				Usage: T("Filter by subnet type"),
			},
			cli.StringFlag{
				Name:  "network-space",
				Usage: T("Filter by network space"),
			},
			cli.BoolFlag{
				Name:  "v4,ipv4",
				Usage: T("Display IPv4 subnets only"),
			},
			cli.BoolFlag{
				Name:  "v6,ipv6",
				Usage: T("Display IPv6 subnets only"),
			},
			cli.IntFlag{
				Name:  "order",
				Usage: T("Filter by the ID of order that purchased the subnets"),
			},
			OutputFlag(),
		},
	}
}

func SubnetLookupMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SUBNET_NAME,
		Name:        CMD_SUBNET_LOOKUP_NAME,
		Description: T("Find an IP address and display its subnet and device information"),
		Usage: T(`${COMMAND_NAME} sl subnet lookup IP_ADDRESS [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl subnet lookup 9.125.235.255
   This command finds the IP address record with IP address 9.125.235.255 and displays its subnet and device information.`),
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	NS_GLOBALIP_NAME  = "globalip"
	CMD_GLOBALIP_NAME = "globalip"

	CMD_GP_ASSIGN_NAME   = "assign"
	CMD_GP_CANCEL_NAME   = "cancel"
	CMD_GP_CREATE_NAME   = "create"
	CMD_GP_LIST_NAME     = "list"
	CMD_GP_UNASSIGN_NAME = "unassign"
)

func GlobalIpNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_GLOBALIP_NAME,
		Description: T("Classic infrastructure Global IP addresses"),
	}
}

func GlobalIpMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_GLOBALIP_NAME,
		Description: T("Classic infrastructure Global IP addresses"),
		Usage:       "${COMMAND_NAME} sl globalip",
		Subcommands: []cli.Command{
			GlobalIpCreateMetaData(),
			GlobalIpAssignMetaData(),
			GlobalIpCancelMetaData(),
			GlobalIpListMetaData(),
			GlobalIpUnassignMetaData(),
		},
	}
}

func GlobalIpCreateMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_GLOBALIP_NAME,
		Name:        CMD_GP_CREATE_NAME,
		Description: T("Create a global IP"),
		Usage: T(`${COMMAND_NAME} sl globalip create [OPTIONS]

EXAMPLE:
    ${COMMAND_NAME} sl globalip create --v6 
	This command creates an IPv6 address.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "v6",
				Usage: T("Order an IPv6 IP address"),
			},
			cli.BoolFlag{
				Name:  "test",
				Usage: T("Test order"),
			},
			ForceFlag(),
			OutputFlag(),
		},
	}
}

func GlobalIpAssignMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_GLOBALIP_NAME,
		Name:        CMD_GP_ASSIGN_NAME,
		Description: T("Assign a global IP to a target router or device"),
		Usage: T(`${COMMAND_NAME} sl globalip assign IDENTIFIER TARGET [OPTIONS]

EXAMPLE:
    ${COMMAND_NAME} sl globalip assign 12345678 9.111.123.456
	This command assigns IP address with ID 12345678 to a target device whose IP address is 9.111.123.456.`),
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

func GlobalIpCancelMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_GLOBALIP_NAME,
		Name:        CMD_GP_CANCEL_NAME,
		Description: T("Cancel a global IP"),
		Usage: T(`${COMMAND_NAME} sl globalip cancel IDENTIFIER [OPTIONS]

EXAMPLE:
    ${COMMAND_NAME} sl globalip cancel 12345678
	This command cancels IP address with ID 12345678.`),
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}

func GlobalIpListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_GLOBALIP_NAME,
		Name:        CMD_GP_LIST_NAME,
		Description: T("List all global IPs on your account"),
		Usage: T(`${COMMAND_NAME} sl globalip list [OPTIONS]

EXAMPLE:
    ${COMMAND_NAME} sl globalip list --v4 
	This command lists all IPv4 addresses on the current account.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "v4",
				Usage: T("Display IPv4 IPs only"),
			},
			cli.BoolFlag{
				Name:  "v6",
				Usage: T("Display IPv6 IPs only"),
			},
			cli.IntFlag{
				Name:  "order",
				Usage: T("Filter by the ID of order that purchased this IP address"),
			},
			OutputFlag(),
		},
	}
}

func GlobalIpUnassignMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_GLOBALIP_NAME,
		Name:        CMD_GP_UNASSIGN_NAME,
		Description: T("Unassign a global IP from a target router or device"),
		Usage: T(`${COMMAND_NAME} sl globalip unassign IDENTIFIER [OPTIONS]

EXAMPLE:
    ${COMMAND_NAME} sl globalip unassign 12345678
	This command unassigns IP address with ID 12345678 from the target device.`),
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

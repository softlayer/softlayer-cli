package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	NS_DEDICATEDHOST_NAME  = "dedicatedhost"
	CMD_DEDICATEDHOST_NAME = "dedicatedhost"

	//sl dedicatedhost
	CMD_DEDICATEDHOST_LIST_NAME          = "list"
	CMD_DEDICATEDHOST_LIST_GUESTS_NAME   = "list-guests"
	CMD_DEDICATEDHOST_CREATE_NAME        = "create"
	CMD_DEDICATEDHOST_CANCEL_GUESTS_NAME = "cancel-guests"
)

func DedicatedhostNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_DEDICATEDHOST_NAME,
		Description: T("Classic infrastructure Dedicatedhost"),
	}
}

func DedicatedhostMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_DEDICATEDHOST_NAME,
		Description: T("Classic infrastructure Dedicatedhost"),
		Usage:       "${COMMAND_NAME} sl dedicatedhost",
		Subcommands: []cli.Command{
			DedicatedhostListMetaData(),
			DedicatedhostListGuestsMetaData(),
			DedicatedhostCreateMetaData(),
			DedicatedhostCancelGuestsMetaData(),
		},
	}
}

func DedicatedhostListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_DEDICATEDHOST_NAME,
		Name:        CMD_DEDICATEDHOST_LIST_NAME,
		Description: T("List Dedicated Host."),
		Usage: T(`${COMMAND_NAME} sl dedicatedhost list [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl dedicatedhost list -d dal09 --sortby diskCapacity
   This command list all Dedicated Host in the Account.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "c,cpu",
				Usage: T("Filter by the number of CPU cores"),
			},
			cli.StringSliceFlag{
				Name:  "t,tag",
				Usage: T("Filter by tags"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Filter by Datacenter shortname"),
			},
			cli.StringFlag{
				Name:  "H,name",
				Usage: T("Filter by host portion of the FQDN"),
			},
			cli.IntFlag{
				Name:  "m,memory",
				Usage: T("Filter by Memory capacity in mebibytes"),
			},
			cli.StringFlag{
				Name:  "d,disk",
				Usage: T("Filter by Disk capacity"),
			},
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by, default:id, options are: id,name,cpuCount,createDate,diskCapacity,memoryCapacity,datacenter"),
			},
			cli.StringSliceFlag{
				Name:  "column",
				Usage: T("Column to display. Options are: id,name,cpuCount,createDate,diskCapacity,memoryCapacity,datacenter,modifyDate,guestCount,notes,billingItem. This option can be specified multiple times"),
			},
			cli.StringSliceFlag{
				Name:   "columns",
				Hidden: true,
			},
			OutputFlag(),
		},
	}
}

func DedicatedhostCancelGuestsMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_DEDICATEDHOST_NAME,
		Name:        CMD_DEDICATEDHOST_CANCEL_GUESTS_NAME,
		Description: T("Cancel all virtual guests of the dedicated host immediately."),
		Usage: T(`${COMMAND_NAME} sl dedicatedhost cancel-guests IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl dedicatedhost cancel-guests 1234567
   This command cancel all virtual guests of the dedicated host immediately.`),
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}

func DedicatedhostListGuestsMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_DEDICATEDHOST_NAME,
		Name:        CMD_DEDICATEDHOST_LIST_GUESTS_NAME,
		Description: T("List Dedicated Host Guests."),
		Usage: T(`${COMMAND_NAME} sl dedicatedhost list-guests IDENTIFIER[OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl dedicatedhost list-guests -d dal09 --sortby hostname 1234567
   This command list all Dedicated Host guests in the Account.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "c,cpu",
				Usage: T("Filter by the number of CPU cores"),
			},
			cli.StringSliceFlag{
				Name:  "t,tag",
				Usage: T("Filter by tags"),
			},
			cli.StringFlag{
				Name:  "d,domain",
				Usage: T("Filter by domain portion of the FQDN"),
			},
			cli.StringFlag{
				Name:  "H,hostname",
				Usage: T("Filter by host portion of the FQDN"),
			},
			cli.IntFlag{
				Name:  "m,memory",
				Usage: T("Filter by Memory capacity in megabytes"),
			},
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by, default:hostname"),
			},
			cli.StringSliceFlag{
				Name:  "column",
				Usage: T("Column to display. [Options are: guid, cpu, memory, datacenter, primary_ip, backend_ip, created_by, power_state, tags] [default: id,hostname,domain,primary_ip,backend_ip,power_state]"),
			},
			cli.StringSliceFlag{
				Name:   "columns",
				Hidden: true,
			},
			OutputFlag(),
		},
	}
}

func DedicatedhostCreateMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_DEDICATEDHOST_NAME,
		Name:        CMD_DEDICATEDHOST_CREATE_NAME,
		Description: T("Create a dedicatedhost"),
		Usage:       "${COMMAND_NAME} sl dedicatedhost create [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "H,hostname",
				Usage: T("Host portion of the FQDN [required]"),
			},
			cli.StringFlag{
				Name:  "D,domain",
				Usage: T("Domain portion of the FQDN [required]"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Datacenter shortname [required]"),
			},
			cli.StringFlag{
				Name:  "s,size",
				Usage: T("Size of the dedicated host, currently only one size is available: 56_CORES_X_242_RAM_X_1_4_TB"),
			},
			cli.StringFlag{
				Name:  "b,billing",
				Usage: T("Billing rate. Default is: hourly. Options are: hourly, monthly"),
			},
			cli.StringFlag{
				Name:  "v,vlan-private",
				Usage: T("The ID of the private VLAN on which you want the dedicated host placed. See: '${COMMAND_NAME} sl vlan list' for reference"),
			},
			cli.BoolFlag{
				Name:  "test",
				Usage: T("Do not actually create the dedicatedhost"),
			},
			ForceFlag(),
			OutputFlag(),
		},
	}
}

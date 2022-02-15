package dedicatedhost

/*

// TODO Actually implement this. 


import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)


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

*/
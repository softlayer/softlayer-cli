package subnet

import (
	"sort"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewListCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *ListCommand) {
	return &ListCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *ListCommand) Run(c *cli.Context) error {
	version := 0
	if c.IsSet("v4") {
		version = 4
	} else if c.IsSet("v6") {
		version = 6
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	mask := "hardware,datacenter,ipAddressCount,virtualGuests,networkVlan[id,networkSpace],subnetType,id,networkIdentifier"
	subnets, err := cmd.NetworkManager.ListSubnets(c.String("identifier"), c.String("d"), version, c.String("t"), c.String("network-space"), c.Int("order"), mask)
	if err != nil {
		return cli.NewExitError(T("Failed to list subnets on your account.\n")+err.Error(), 2)
	}
	sortby := c.String("sortby")
	if sortby == "" || sortby == "id" || sortby == "ID" {
		sort.Sort(utils.SubnetById(subnets))
	} else if sortby == "identifier" {
		sort.Sort(utils.SubnetByIdentifier(subnets))
	} else if sortby == "type" {
		sort.Sort(utils.SubnetByType(subnets))
	} else if sortby == "network_space" {
		sort.Sort(utils.SubnetByNetworkSpace(subnets))
	} else if sortby == "datacenter" {
		sort.Sort(utils.SubnetByDatacenter(subnets))
	} else if sortby == "vlan_id" {
		sort.Sort(utils.SubnetByVlanId(subnets))
	} else if sortby == "IPs" {
		sort.Sort(utils.SubnetByIpCount(subnets))
	} else if sortby == "hardware" {
		sort.Sort(utils.SubnetByHardwareCount(subnets))
	} else if sortby == "vs" {
		sort.Sort(utils.SubnetByVSCount(subnets))
	} else {
		return errors.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.", map[string]interface{}{"Column": sortby}))
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, subnets)
	}

	headers := []string{T("ID"), T("identifier"), T("type"), T("network_space"), T("datacenter"), T("vlan_id"), T("IPs"), T("hardware"), T("virtual_servers")}
	table := cmd.UI.Table(headers)
	for _, subnet := range subnets {

		var networktype, datacenter, vlanID string

		if subnet.NetworkVlan != nil {
			networktype = utils.FormatStringPointer(subnet.NetworkVlan.NetworkSpace)
		}
		if subnet.Datacenter != nil {
			datacenter = utils.FormatStringPointer(subnet.Datacenter.Name)
		}
		if subnet.NetworkVlan != nil {
			vlanID = utils.FormatIntPointer(subnet.NetworkVlan.Id)
		}

		table.Add(utils.FormatIntPointer(subnet.Id),
			utils.FormatStringPointer(subnet.NetworkIdentifier),
			utils.FormatStringPointer(subnet.SubnetType),
			networktype, datacenter, vlanID,
			utils.FormatUIntPointer(subnet.IpAddressCount),
			strconv.Itoa(len(subnet.Hardware)),
			strconv.Itoa(len(subnet.VirtualGuests)))
	}
	table.Print()
	return nil
}

func SubnetListMetaData() cli.Command {
	return cli.Command{
		Category:    "subnet",
		Name:        "list",
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
			metadata.OutputFlag(),
		},
	}
}

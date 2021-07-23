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

	mask := "hardware,datacenter,ipAddressCount,virtualGuests,networkVlan[id,networkSpace],subnetType,id,networkIdentifier,ipAddresses[id, ipAddress,note]"
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

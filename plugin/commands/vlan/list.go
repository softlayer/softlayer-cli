package vlan

import (
	"sort"

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
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	vlans, err := cmd.NetworkManager.ListVlans(c.String("d"), c.Int("n"), c.String("name"), c.Int("order"), "")
	if err != nil {
		return cli.NewExitError(T("Failed to list VLANs on your account.\n")+err.Error(), 2)
	}

	sortby := c.String("sortby")
	if sortby == "id" || sortby == "ID" {
		sort.Sort(utils.VlanById(vlans))
	} else if sortby == "number" {
		sort.Sort(utils.VlanByNumber(vlans))
	} else if sortby == "name" {
		sort.Sort(utils.VlanByName(vlans))
	} else if sortby == "firewall" {
		sort.Sort(utils.VlanByFirewall(vlans))
	} else if sortby == "datacenter" {
		sort.Sort(utils.VlanByDatacenter(vlans))
	} else if sortby == "hardware" {
		sort.Sort(utils.VlanByHardwareCount(vlans))
	} else if sortby == "virtual_servers" {
		sort.Sort(utils.VlanByVirtualServerCount(vlans))
	} else if sortby == "public_ips" {
		sort.Sort(utils.VlanByPublicIPCount(vlans))
	} else if sortby == "" {
		//do nothing
	} else {
		return errors.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.", map[string]interface{}{"Column": sortby}))
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, vlans)
	}

	headers := []string{T("ID"), T("number"), T("name"), T("firewall"), T("network_space"), T("primary_router"), T("hardware"), T("virtual_servers"), T("public_ips")}
	table := cmd.UI.Table(headers)

	for _, vlan := range vlans {
		var firewall string
		if len(vlan.FirewallInterfaces) > 0 {
			firewall = T("Yes")
		} else {
			firewall = T("No")
		}
		hostName := ""
		if vlan.PrimaryRouter != nil && vlan.PrimaryRouter.Hostname != nil {
			hostName = *vlan.PrimaryRouter.Hostname
		}

		table.Add(utils.FormatIntPointer(vlan.Id),
			utils.FormatIntPointer(vlan.VlanNumber),
			utils.FormatStringPointer(vlan.Name),
			firewall,
			utils.FormatStringPointer(vlan.NetworkSpace),
			utils.FormatStringPointer(&hostName),
			utils.FormatUIntPointer(vlan.HardwareCount),
			utils.FormatUIntPointer(vlan.VirtualGuestCount),
			utils.FormatUIntPointer(vlan.TotalPrimaryIpAddressCount))
	}

	table.Print()
	return nil
}

func VlanListMetaData() cli.Command {
	return cli.Command{
		Category:    "vlan",
		Name:        "list",
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
			metadata.OutputFlag(),
		},
	}
}

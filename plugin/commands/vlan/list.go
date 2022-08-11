package vlan

import (
	"sort"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

	mask := ""
	pods, err := cmd.NetworkManager.GetPods(mask)
	if err != nil {
		return cli.NewExitError(T("Failed to get Pods.\n")+err.Error(), 2)
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

	headers := []string{T("Id"), T("Number"), T("Fully qualified name"), T("Name"), T("Network"), T("Data center"), T("Pod"), T("Gateway/Firewall"), T("Hardware"), T("Virtual servers"), T("Public ips"), T("Premium"), T("Tags")}
	table := cmd.UI.Table(headers)

	for _, vlan := range vlans {
		var premium string
		if vlan.BillingItem != nil {
			premium = T("Yes")
		} else {
			premium = T("No")
		}
		table.Add(
			utils.FormatIntPointer(vlan.Id),
			utils.FormatIntPointer(vlan.VlanNumber),
			utils.FormatStringPointer(vlan.FullyQualifiedName),
			utils.FormatStringPointer(vlan.Name),
			cases.Title(language.Und).String(utils.FormatStringPointer(vlan.NetworkSpace)),
			utils.FormatStringPointer(vlan.PrimaryRouter.Datacenter.Name),
			getPodWithClosedAnnouncement(vlan, pods),
			getFirewallGateway(vlan),
			utils.FormatUIntPointer(vlan.HardwareCount),
			utils.FormatUIntPointer(vlan.VirtualGuestCount),
			utils.FormatUIntPointer(vlan.TotalPrimaryIpAddressCount),
			premium,
			utils.TagRefsToString(vlan.TagReferences),
		)
	}

	table.Print()
	return nil
}

func getFirewallGateway(vlan datatypes.Network_Vlan) string {
	if vlan.NetworkVlanFirewall != nil {
		return utils.FormatStringPointer(vlan.NetworkVlanFirewall.FullyQualifiedDomainName)
	}
	if vlan.AttachedNetworkGateway != nil {
		return utils.FormatStringPointer(vlan.AttachedNetworkGateway.Name)
	}
	return "-"
}

func getPodWithClosedAnnouncement(vlan datatypes.Network_Vlan, pods []datatypes.Network_Pod) string {
	for _, pod := range pods {
		if *pod.BackendRouterId == *vlan.PrimaryRouter.Id || *pod.FrontendRouterId == *vlan.PrimaryRouter.Id {
			namePod := cases.Title(language.Und).String(strings.Split(utils.FormatStringPointer(pod.Name), ".")[1])
			if utils.WordInList(pod.Capabilities, "CLOSURE_ANNOUNCED") {
				return namePod + "*"
			} else {
				return namePod
			}
		}
	}
	return ""
}

func VlanListMetaData() cli.Command {
	return cli.Command{
		Category:    "vlan",
		Name:        "list",
		Description: T("List all the VLANs on your account"),
		Usage: T(`${COMMAND_NAME} sl vlan list [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vlan list -d dal09 --sortby number
   This commands lists all vlans on current account filtering by datacenter equals to dal09, and sort them by vlan number.
 
Note: In field Pod, if add (*) indicated that closed soon`),
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

package subnet

import (
	"sort"
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ListCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	Sortby         string
	Datacenter     string
	Identifier     string
	SubnetType     string
	NetworkSpace   string
	Ipv4           bool
	Ipv6           bool
	Order          int
}

func NewListCommand(sl *metadata.SoftlayerCommand) *ListCommand {
	thisCmd := &ListCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List all subnets on your account"),
		Long: T(`${COMMAND_NAME} sl subnet list [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl subnet list -d dal09 -t PRIMARY --network-space PUBLIC --v4
   This command lists IPv4 subnets on the current account, and filters by datacenter is dal09, subnet type is PRIMARY, and network space is PUBLIC.`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Sortby, "sortby", "", T("Column to sort by. Options are: id,identifier,type,network_space,datacenter,vlan_id,IPs,hardware,vs"))
	cobraCmd.Flags().StringVarP(&thisCmd.Datacenter, "datacenter", "d", "", T("Filter by datacenter shortname"))
	cobraCmd.Flags().StringVar(&thisCmd.Identifier, "identifier", "", T("Filter by network identifier"))
	cobraCmd.Flags().StringVarP(&thisCmd.SubnetType, "subnet-type", "t", "", T("Filter by subnet type"))
	cobraCmd.Flags().StringVar(&thisCmd.NetworkSpace, "network-space", "", T("Filter by network space"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Ipv4, "ipv4", "4", false, T("Display IPv4 subnets only"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Ipv6, "ipv6", "6", false, T("Display IPv6 subnets only"))
	cobraCmd.Flags().IntVar(&thisCmd.Order, "order", 0, T("Filter by the ID of order that purchased the subnets"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ListCommand) Run(args []string) error {
	version := 0
	if cmd.Ipv4 {
		version = 4
	} else if cmd.Ipv6 {
		version = 6
	}

	outputFormat := cmd.GetOutputFlag()

	mask := "hardware,datacenter,ipAddressCount,virtualGuests,networkVlan[id,networkSpace,fullyQualifiedName],subnetType,id,networkIdentifier,addressSpace,endPointIpAddress,note,tagReferences[tag]"
	subnets, err := cmd.NetworkManager.ListSubnets(cmd.Identifier, cmd.Datacenter, version, cmd.SubnetType, cmd.NetworkSpace, cmd.Order, mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to list subnets on your account.\n"), err.Error(), 2)
	}
	sortby := cmd.Sortby
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

	headers := []string{T("ID"), T("Identifier"), T("Network"), T("Type"), T("VLAN"), T("Location"), T("Target"), T("IPs"), T("Hardware"), T("Vs"), T("Tags"), T("Note")}
	table := cmd.UI.Table(headers)
	for _, subnet := range subnets {

		var networktype, datacenter, subnetType, note string
		target := "-"
		vlan := ""

		networktype = cases.Title(language.Und).String(utils.FormatStringPointer(subnet.AddressSpace))

		if subnet.NetworkVlan != nil {
			vlan = utils.FormatStringPointer(subnet.NetworkVlan.FullyQualifiedName)
		}
		if vlan == "" {
			vlan = "Unrouted"
		}

		subnetType = utils.FormatStringPointer(subnet.SubnetType)
		if subnetType == "PRIMARY" || subnetType == "PRIMARY_6" || subnetType == "ADDITIONAL_PRIMARY" {
			subnetType = "Primary"
		}
		if subnetType == "SECONDARY_ON_VLAN" {
			subnetType = "Portable"
		}
		if subnetType == "STATIC_IP_ROUTED" {
			subnetType = "Static"
		}
		if subnetType == "GLOBAL_IP" {
			subnetType = "Global"
		}

		if subnet.Datacenter != nil {
			datacenter = utils.FormatStringPointer(subnet.Datacenter.Name)
		}

		if subnet.EndPointIpAddress != nil {
			target = utils.FormatStringPointer(subnet.EndPointIpAddress.IpAddress)
		}
		note = utils.TagRefsToString(subnet.TagReferences)
		if note == "" {
			note = "-"
		}

		table.Add(utils.FormatIntPointer(subnet.Id),
			utils.FormatStringPointer(subnet.NetworkIdentifier),
			networktype,
			subnetType,
			vlan,
			datacenter,
			target,
			utils.FormatUIntPointer(subnet.IpAddressCount),
			strconv.Itoa(len(subnet.Hardware)),
			strconv.Itoa(len(subnet.VirtualGuests)),
			note,
			utils.FormatStringPointer(subnet.Note),
		)
	}
	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

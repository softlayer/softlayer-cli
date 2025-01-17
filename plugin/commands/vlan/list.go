package vlan

import (
	"sort"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
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
	Number         int
	Name           string
	Order          int
}

func NewListCommand(sl *metadata.SoftlayerCommand) *ListCommand {
	thisCmd := &ListCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List all the VLANs on your account"),
		Long: T(`EXAMPLE:
   ${COMMAND_NAME} sl vlan list -d dal09 --sortby number
   This commands lists all vlans on current account filtering by datacenter equals to dal09, and sort them by vlan number.
 
Note: In field Pod, if add (*) indicated that closed soon`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Sortby, "sortby", "", T("Column to sort by. Options are: id,number,name,firewall,datacenter,hardware,virtual_servers,public_ips"))
	cobraCmd.Flags().StringVarP(&thisCmd.Datacenter, "datacenter", "d", "", T("Filter by datacenter shortname"))
	cobraCmd.Flags().IntVarP(&thisCmd.Number, "number", "n", 0, T("Filter by VLAN number"))
	cobraCmd.Flags().StringVar(&thisCmd.Name, "name", "", T("Filter by VLAN name"))
	cobraCmd.Flags().IntVar(&thisCmd.Order, "order", 0, T("Filter by ID of the order that purchased the VLAN"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ListCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	vlans, err := cmd.NetworkManager.ListVlans(cmd.Datacenter, cmd.Number, cmd.Name, cmd.Order, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to list VLANs on your account.\n"), err.Error(), 2)
	}

	mask := ""
	pods, err := cmd.NetworkManager.GetPods(mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get Pods.\n"), err.Error(), 2)
	}

	sortby := cmd.Sortby
	if sortby == "id" || sortby == "ID" {
		sort.Slice(vlans, func(i, j int) bool {
			return *vlans[i].Id < *vlans[j].Id
		})
	} else if sortby == "number" {
		sort.Slice(vlans, func(i, j int) bool {
			return *vlans[i].VlanNumber < *vlans[j].VlanNumber
		})
	} else if sortby == "name" {
		sort.Slice(vlans, func(i, j int) bool {
			return *vlans[i].VlanNumber < *vlans[j].VlanNumber
		})
	} else if sortby == "firewall" {
		sort.Slice(vlans, func(i, j int) bool {
			return len(vlans[i].FirewallInterfaces) < len(vlans[j].FirewallInterfaces)
		})
	} else if sortby == "datacenter" {
		sort.Slice(vlans, func(i, j int) bool {
			if vlans[i].Datacenter != nil && vlans[i].Datacenter.Name != nil &&
				vlans[j].Datacenter != nil && vlans[j].Datacenter.Name != nil {
				return *vlans[i].Datacenter.Name < *vlans[j].Datacenter.Name
			}
			return false
		})
	} else if sortby == "hardware" {
		sort.Slice(vlans, func(i, j int) bool {
			return *vlans[i].HardwareCount < *vlans[j].HardwareCount
		})
	} else if sortby == "virtual_servers" {
		sort.Slice(vlans, func(i, j int) bool {
			return *vlans[i].VirtualGuestCount < *vlans[j].VirtualGuestCount
		})
	} else if sortby == "public_ips" {
		sort.Slice(vlans, func(i, j int) bool {
			return *vlans[i].TotalPrimaryIpAddressCount < *vlans[j].TotalPrimaryIpAddressCount
		})
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
		datacenterName := ""
		if vlan.Datacenter != nil && vlan.Datacenter.Name != nil {
			datacenterName = utils.FormatStringPointer(vlan.Datacenter.Name)
		}
		table.Add(
			utils.FormatIntPointer(vlan.Id),
			utils.FormatIntPointer(vlan.VlanNumber),
			utils.FormatStringPointer(vlan.FullyQualifiedName),
			utils.FormatStringPointer(vlan.Name),
			cases.Title(language.Und).String(utils.FormatStringPointer(vlan.NetworkSpace)),
			datacenterName,
			getPodWithClosedAnnouncement(vlan, pods),
			getFirewallGateway(vlan),
			utils.FormatUIntPointer(vlan.HardwareCount),
			utils.FormatUIntPointer(vlan.VirtualGuestCount),
			utils.FormatUIntPointer(vlan.TotalPrimaryIpAddressCount),
			premium,
			utils.TagRefsToString(vlan.TagReferences),
		)
	}
	utils.PrintTable(cmd.UI, table, outputFormat)
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
		if *pod.Name == *vlan.PodName {
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

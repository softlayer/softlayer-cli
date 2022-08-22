package firewall

import (
	"fmt"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	*metadata.SoftlayerCommand
	FirewallManager managers.FirewallManager
	Command         *cobra.Command
}

func NewListCommand(sl *metadata.SoftlayerCommand) (cmd *ListCommand) {
	thisCmd := &ListCommand{
		SoftlayerCommand: sl,
		FirewallManager:  managers.NewFirewallManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List all firewalls on your account."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ListCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	table := cmd.UI.Table([]string{T("Firewall ID"), T("Type"), T("Features"), T("Server/Vlan Id")})
	fwvlans, err := cmd.FirewallManager.GetFirewalls()
	if err != nil {
		return errors.NewAPIError(T("Failed to get firewalls on your account."), err.Error(), 2)
	}
	multiVlanFirewalls, err := cmd.FirewallManager.GetMultiVlanFirewalls("")
	if err != nil {
		return errors.NewAPIError(T("Failed to get multi vlan firewalls on your account."), err.Error(), 2)
	}

	//dedicated firewalls
	for _, vlan := range fwvlans {
		if vlan.NetworkVlanFirewall == nil {
			continue
		}
		if (vlan.DedicatedFirewallFlag != nil && *vlan.DedicatedFirewallFlag == 0) || vlan.DedicatedFirewallFlag == nil {
			continue
		}
		features := "-"
		if vlan.HighAvailabilityFirewallFlag != nil && *vlan.HighAvailabilityFirewallFlag == true {
			features = "HA"
		}
		var firewallID int
		if vlan.NetworkVlanFirewall != nil && vlan.NetworkVlanFirewall.Id != nil {
			firewallID = *vlan.NetworkVlanFirewall.Id
		}
		table.Add(fmt.Sprintf("vlan:%d",
			firewallID),
			"VLAN - dedicated",
			features,
			utils.FormatIntPointer(vlan.Id))
	}

	//shared firewalls
	for _, vlan := range fwvlans {
		if vlan.DedicatedFirewallFlag != nil && *vlan.DedicatedFirewallFlag == 1 {
			continue
		}
		for _, vs := range vlan.FirewallGuestNetworkComponents {
			if hasFirewallComponent(vs) {
				table.Add(fmt.Sprintf("vs:%d", utils.IntPointertoInt(vs.Id)),
					"Virtual Server - standard", "-",
					utils.FormatIntPointer(vs.GuestNetworkComponent.GuestId))
			}
		}

		for _, hw := range vlan.FirewallNetworkComponents {
			if hasFirewallComponent(hw) {
				table.Add(fmt.Sprintf("server:%d", utils.IntPointertoInt(hw.Id)),
					"Hardware Server - standard", "-",
					utils.FormatIntPointer(hw.NetworkComponent.DownlinkComponent.HardwareId))
			}
		}
	}

	utils.PrintTable(cmd.UI, table, outputFormat)

	cmd.UI.Print("\n")
	table = cmd.UI.Table([]string{T("Firewall ID"), T("Firewall"), T("Type"), T("Hostname"), T("Location"), T("Public Ip"), T("Private Ip"), T("Associated VLANs"), T("Status")})
	//multi vlan firewalls
	for _, firewall := range multiVlanFirewalls {
		table.Add(
			fmt.Sprintf("multiVlan:%d", *firewall.NetworkFirewall.Id),
			utils.FormatStringPointer(firewall.Name),
			utils.FormatStringPointer(firewall.NetworkFirewall.FirewallType),
			utils.FormatStringPointer(firewall.Members[0].Hardware.Hostname),
			utils.FormatStringPointer(firewall.NetworkFirewall.Datacenter.Name),
			utils.FormatStringPointer(firewall.PublicIpAddress.IpAddress),
			utils.FormatStringPointer(firewall.PrivateIpAddress.IpAddress),
			fmt.Sprintf("%d VLANs", len(firewall.InsideVlans)),
			utils.FormatStringPointer(firewall.Status.KeyName),
		)
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

func hasFirewallComponent(component datatypes.Network_Component_Firewall) bool {
	if component.Status != nil && *component.Status != "no_edit" {
		return true
	}
	return false
}

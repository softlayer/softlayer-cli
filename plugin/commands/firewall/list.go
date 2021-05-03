package firewall

import (
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	UI              terminal.UI
	FirewallManager managers.FirewallManager
}

func NewListCommand(ui terminal.UI, firewallManager managers.FirewallManager) (cmd *ListCommand) {
	return &ListCommand{
		UI:              ui,
		FirewallManager: firewallManager,
	}
}

func (cmd *ListCommand) Run(c *cli.Context) error {
	table := cmd.UI.Table([]string{T("firewall id"), T("type"), T("features"), T("server/vlan id")})
	fwvlans, err := cmd.FirewallManager.GetFirewalls()
	if err != nil {
		return cli.NewExitError(T("Failed to get firewalls on your account.\n")+err.Error(), 2)
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
	table.Print()
	return nil
}

func hasFirewallComponent(component datatypes.Network_Component_Firewall) bool {
	if component.Status != nil && *component.Status != "no_edit" {
		return true
	}
	return false
}

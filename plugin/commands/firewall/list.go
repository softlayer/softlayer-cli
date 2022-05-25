package firewall

import (
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
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

func FirewallListMetaData() cli.Command {
	return cli.Command{
		Category:    "firewall",
		Name:        "list",
		Description: T("List all firewalls on your account"),
		Usage: T(`${COMMAND_NAME} sl firewall list [OPTIONS]

EXAMPLE: 
   ${COMMAND_NAME} sl firewall list`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

func (cmd *ListCommand) Run(c *cli.Context) error {

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	table := cmd.UI.Table([]string{T("Firewall ID"), T("Type"), T("Features"), T("Server/Vlan Id")})
	fwvlans, err := cmd.FirewallManager.GetFirewalls()
	if err != nil {
		return cli.NewExitError(T("Failed to get firewalls on your account.\n")+err.Error(), 2)
	}
	multiVlanFirewalls, err := cmd.FirewallManager.GetMultiVlanFirewalls("")
	if err != nil {
		return cli.NewExitError(T("Failed to get multi vlan firewalls on your account.\n")+err.Error(), 2)
	}
	//multi vlan firewalls
	for _, firewall := range multiVlanFirewalls {
		features := "-"
		if *firewall.MemberCount > 1 {
			features = "HA"
		}
		vlans := "-"
		if firewall.InsideVlans != nil && len(firewall.InsideVlans) != 0 {
			for _, vlan := range firewall.InsideVlans {
				vlans = vlans + strconv.Itoa(*vlan.NetworkVlanId) + "\n"
			}
		}
		table.Add(
			fmt.Sprintf("multiVlan:%d", *firewall.NetworkFirewall.Id),
			utils.FormatStringPointer(firewall.NetworkFirewall.FirewallType),
			features,
			vlans)
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
	return nil
}

func hasFirewallComponent(component datatypes.Network_Component_Firewall) bool {
	if component.Status != nil && *component.Status != "no_edit" {
		return true
	}
	return false
}

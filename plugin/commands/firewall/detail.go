package firewall

import (
	"bytes"
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	*metadata.SoftlayerCommand
	FirewallManager managers.FirewallManager
	Command         *cobra.Command
	Credentials     bool
}

func NewDetailCommand(sl *metadata.SoftlayerCommand) (cmd *DetailCommand) {
	thisCmd := &DetailCommand{
		SoftlayerCommand: sl,
		FirewallManager:  managers.NewFirewallManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "detail " + T("IDENTIFIER"),
		Short: T("Detail information about a firewall"),
		Long: T(`${COMMAND_NAME} sl firewall detail IDENTIFIER [OPTIONS]
		
EXAMPLE: 
${COMMAND_NAME} sl firewall detail vs:12345
${COMMAND_NAME} sl firewall detail server:234567
${COMMAND_NAME} sl firewall detail vlan:345678
${COMMAND_NAME} sl firewall detail multiVlan:456789`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVar(&thisCmd.Credentials, "credentials", false,
		T("Display FortiGate username and FortiGate password to multi vlans"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DetailCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	firewallType, firewallID, err := cmd.FirewallManager.ParseFirewallID(args[0])
	if err != nil {
		return errors.NewAPIError(T("Failed to parse firewall ID : {{.FirewallID}}.\n", map[string]interface{}{"FirewallID": args[0]}), err.Error(), 1)
	}

	var table terminal.Table
	if firewallType == "multiVlan" {
		firewall, err := cmd.FirewallManager.GetMultiVlanFirewall(firewallID, "")
		if err != nil {
			return errors.NewAPIError(T("Failed to get multi vlan firewall.\n"), err.Error(), 2)
		}
		table = cmd.UI.Table([]string{T("Name"), T("Value")})
		table.Add(T("Name"), utils.FormatStringPointer(firewall.NetworkGateway.Name))
		table.Add(T("Datacenter"), utils.FormatStringPointer(firewall.Datacenter.LongName))
		table.Add(T("Public IP"), utils.FormatStringPointer(firewall.NetworkGateway.PublicIpAddress.IpAddress))
		table.Add(T("Private IP"), utils.FormatStringPointer(firewall.NetworkGateway.PrivateIpAddress.IpAddress))
		table.Add(T("Public IPv6"), utils.FormatStringPointer(firewall.NetworkGateway.PublicIpv6Address.IpAddress))
		table.Add(T("Public VLAN"), utils.FormatIntPointer(firewall.NetworkGateway.PublicVlan.VlanNumber))
		table.Add(T("Private VLAN"), utils.FormatIntPointer(firewall.NetworkGateway.PrivateVlan.VlanNumber))
		table.Add(T("Type"), utils.FormatStringPointer(firewall.FirewallType))

		if cmd.Credentials {
			table.Add(T("FortiGate username"), utils.FormatStringPointer(firewall.ManagementCredentials.Username))
			table.Add(T("FortiGate password"), utils.FormatStringPointer(firewall.ManagementCredentials.Password))
		}

		if rules := firewall.Rules; len(rules) > 0 {
			buf := new(bytes.Buffer)
			ruleTable := terminal.NewTable(buf, []string{"#", T("action"), T("protocol"), T("src_ip"), T("src_mask"), T("dest"), T("dest_mask")})
			for _, rule := range rules {
				ruleTable.Add(utils.FormatIntPointer(rule.OrderValue),
					utils.FormatStringPointer(rule.Action),
					utils.FormatStringPointer(rule.Protocol),
					utils.FormatStringPointer(rule.SourceIpAddress),
					utils.FormatStringPointer(rule.SourceIpSubnetMask),
					fmt.Sprintf("%s:%s-%s", utils.FormatStringPointer(rule.DestinationIpAddress), utils.FormatIntPointer(rule.DestinationPortRangeStart), utils.FormatIntPointer(rule.DestinationPortRangeEnd)),
					utils.FormatStringPointer(rule.DestinationIpSubnetMask),
				)
			}
			ruleTable.Print()
			table.Add("Rules", buf.String())
		}

	} else {
		table = cmd.UI.Table([]string{"#", T("action"), T("protocol"), T("src_ip"), T("src_mask"), T("dest"), T("dest_mask")})
		if firewallType == "vlan" {
			firewallRules, err := cmd.FirewallManager.GetDedicatedFirewallRules(firewallID)
			if err != nil {
				return errors.NewAPIError(T("Failed to get dedicated firewall rules.\n"), err.Error(), 2)
			}
			for _, rule := range firewallRules {
				table.Add(utils.FormatIntPointer(rule.OrderValue),
					utils.FormatStringPointer(rule.Action),
					utils.FormatStringPointer(rule.Protocol),
					utils.FormatStringPointer(rule.SourceIpAddress),
					utils.FormatStringPointer(rule.SourceIpSubnetMask),
					fmt.Sprintf("%s:%s-%s", utils.FormatStringPointer(rule.DestinationIpAddress), utils.FormatIntPointer(rule.DestinationPortRangeStart), utils.FormatIntPointer(rule.DestinationPortRangeEnd)),
					utils.FormatStringPointer(rule.DestinationIpSubnetMask))
			}
		} else {
			firewallRules, err := cmd.FirewallManager.GetStandardFirewallRules(firewallID)
			if err != nil {
				return errors.NewAPIError(T("Failed to get standard firewall rules.\n"), err.Error(), 2)
			}
			for _, rule := range firewallRules {
				table.Add(utils.FormatIntPointer(rule.OrderValue),
					utils.FormatStringPointer(rule.Action),
					utils.FormatStringPointer(rule.Protocol),
					utils.FormatStringPointer(rule.SourceIpAddress),
					utils.FormatStringPointer(rule.SourceIpSubnetMask),
					fmt.Sprintf("%s:%s-%s", utils.FormatStringPointer(rule.DestinationIpAddress), utils.FormatIntPointer(rule.DestinationPortRangeStart), utils.FormatIntPointer(rule.DestinationPortRangeEnd)),
					utils.FormatStringPointer(rule.DestinationIpSubnetMask))
			}
		}
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

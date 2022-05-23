package firewall

import (
	"bytes"
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	UI              terminal.UI
	FirewallManager managers.FirewallManager
}

func NewDetailCommand(ui terminal.UI, firewallManager managers.FirewallManager) (cmd *DetailCommand) {
	return &DetailCommand{
		UI:              ui,
		FirewallManager: firewallManager,
	}
}

func FirewallDetailMetaData() cli.Command {
	return cli.Command{
		Category:    "firewall",
		Name:        "detail",
		Description: T("Detail information about a firewall"),
		Usage: T(`${COMMAND_NAME} sl firewall detail IDENTIFIER [OPTIONS]
		
EXAMPLE: 
${COMMAND_NAME} sl firewall detail vs:12345
${COMMAND_NAME} sl firewall detail server:234567
${COMMAND_NAME} sl firewall detail vlan:345678
${COMMAND_NAME} sl firewall detail multiVlan:456789`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "credentials",
				Usage: T("Display FortiGate username and FortiGate password to multi vlans"),
			},
			metadata.OutputFlag(),
		},
	}
}

func (cmd *DetailCommand) Run(c *cli.Context) error {

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	firewallType, firewallID, err := cmd.FirewallManager.ParseFirewallID(c.Args()[0])
	if err != nil {
		return cli.NewExitError(T("Failed to parse firewall ID : {{.FirewallID}}.\n", map[string]interface{}{"FirewallID": c.Args()[0]})+err.Error(), 1)
	}

	var table terminal.Table
	if firewallType == "multiVlan" {
		firewall, err := cmd.FirewallManager.GetMultiVlanFirewall(firewallID, "")
		if err != nil {
			return cli.NewExitError(T("Failed to get multi vlan firewall.\n")+err.Error(), 2)
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

		if c.IsSet("credentials") && c.Bool("credentials") {
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
				return cli.NewExitError(T("Failed to get dedicated firewall rules.\n")+err.Error(), 2)
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
				return cli.NewExitError(T("Failed to get standard firewall rules.\n")+err.Error(), 2)
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

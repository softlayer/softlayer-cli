package firewall

import (
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
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
		Usage:       "${COMMAND_NAME} sl firewall detail  IDENTIFIER [OPTIONS]",
	}
}


func (cmd *DetailCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	firewallType, firewallID, err := cmd.FirewallManager.ParseFirewallID(c.Args()[0])
	if err != nil {
		return cli.NewExitError(T("Failed to parse firewall ID : {{.FirewallID}}.\n", map[string]interface{}{"FirewallID": c.Args()[0]})+err.Error(), 1)
	}

	table := cmd.UI.Table([]string{"#", T("action"), T("protocol"), T("src_ip"), T("src_mask"), T("dest"), T("dest_mask")})
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
	table.Print()
	return nil
}

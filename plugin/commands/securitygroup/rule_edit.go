package securitygroup

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type RuleEditCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewRuleEditCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *RuleEditCommand) {
	return &RuleEditCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *RuleEditCommand) Run(c *cli.Context) error {
	if c.NArg() != 2 {
		return errors.NewInvalidUsageError(T("This command requires two arguments."))
	}
	groupID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group ID")
	}
	ruleID, err := strconv.Atoi(c.Args()[1])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group rule ID")
	}

	if c.IsSet("d") {
		direction := c.String("d")
		if direction != "egress" && direction != "ingress" {
			return errors.NewInvalidUsageError(T("-d|--direction has to be either egress or ingress."))
		}
	}
	if c.IsSet("e") {
		etherType := c.String("e")
		if etherType != "IPv4" && etherType != "IPv6" {
			return errors.NewInvalidUsageError(T("-e|--ether-type has to be either IPv4 or IPv6."))
		}
	}
	if c.IsSet("p") {
		protocol := c.String("p")
		if protocol != "icmp" && protocol != "tcp" && protocol != "udp" {
			return errors.NewInvalidUsageError(T("Options for -p|--protocol are: icmp,tcp,udp"))
		}
	}
	err = cmd.NetworkManager.EditSecurityGroupRule(groupID, ruleID, c.String("r"), c.Int("s"), c.String("d"), c.String("e"), c.Int("M"), c.Int("m"), c.String("p"))
	if err != nil {
		return cli.NewExitError(T("Failed to edit rule {{.RuleId}} in security group {{.GroupID}}.\n",
			map[string]interface{}{"RuleId": ruleID, "GroupID": groupID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Rule {{.RuleId}} in security group {{.GroupID}} is updated.", map[string]interface{}{"RuleId": ruleID, "GroupID": groupID}))
	return nil
}

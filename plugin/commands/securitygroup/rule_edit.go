package securitygroup

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
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

func SecurityGroupRuleEditMetaData() cli.Command {
	return cli.Command{
		Category:    "securitygroup",
		Name:        "rule-edit",
		Description: T("Edit a security group rule in a security group"),
		Usage:       "${COMMAND_NAME} sl securitygroup rule-edit SECURITYGROUP_ID RULE_ID [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "r,remote-ip",
				Usage: T("The remote IP/CIDR to enforce"),
			},
			cli.IntFlag{
				Name:  "s,remote-group",
				Usage: T("The ID of the remote security group to enforce"),
			},
			cli.StringFlag{
				Name:  "d,direction",
				Usage: T("The direction of traffic to enforce (ingress or egress), required"),
			},
			cli.StringFlag{
				Name:  "e,ether-type",
				Usage: T("The ethertype (IPv4 or IPv6) to enforce, default is IPv4 if not specified"),
			},
			cli.IntFlag{
				Name:  "M,port-max",
				Usage: T("The upper port bound to enforce"),
			},
			cli.IntFlag{
				Name:  "m,port-min",
				Usage: T("The lower port bound to enforce"),
			},
			cli.StringFlag{
				Name:  "p,protocol",
				Usage: T("The protocol (icmp, tcp, udp) to enforce"),
			},
		},
	}
}

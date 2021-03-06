package securitygroup

import (
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type RuleAddCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewRuleAddCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *RuleAddCommand) {
	return &RuleAddCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *RuleAddCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	groupID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group ID")
	}
	direction := c.String("d")
	if direction != "egress" && direction != "ingress" {
		return errors.NewInvalidUsageError(T("-d|--direction has to be either egress or ingress."))
	}
	etherType := "IPv4"
	if c.IsSet("e") {
		etherType = c.String("e")
	}
	if etherType != "IPv4" && etherType != "IPv6" {
		return errors.NewInvalidUsageError(T("-e|--ether-type has to be either IPv4 or IPv6."))
	}
	portMax := c.Int("M")
	portMin := c.Int("m")
	protocol := c.String("p")
	if portMax > 0 || portMin > 0 {
		if protocol == "" {
			return errors.NewInvalidUsageError(T("-p|--protocal must be set when -M or -m is specified."))
		}
		protocol = strings.ToLower(protocol)
		if protocol != "icmp" && protocol != "tcp" && protocol != "udp" {
			return errors.NewInvalidUsageError(T("Options for -p|--protocol are: icmp,tcp,udp"))
		}
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	resp, err := cmd.NetworkManager.AddSecurityGroupRule(groupID, c.String("r"), c.Int("s"), direction, etherType, portMax, portMin, protocol)
	if err != nil {
		return cli.NewExitError(T("Failed to add rule to security group {{.GroupID}}.\n", map[string]interface{}{"GroupID": groupID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Rule is added to security group {{.GroupID}}.", map[string]interface{}{"GroupID": groupID}))
	return nil
}

func SecurityGroupRuleAddMetaData() cli.Command {
	return cli.Command{
		Category:    "securitygroup",
		Name:        "rule-add",
		Description: T("Add a security group rule to a security group"),
		Usage:       "${COMMAND_NAME} sl securitygroup rule-add SECURITYGROUP_ID [OPTIONS]",
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
			metadata.OutputFlag(),
		},
	}
}

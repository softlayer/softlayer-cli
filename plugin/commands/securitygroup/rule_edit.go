package securitygroup

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type RuleEditCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	RemoteIp       string
	RemoteGroup    int
	Direction      string
	EtherType      string
	PortMax        int
	PortMin        int
	protocol       string
}

func NewRuleEditCommand(sl *metadata.SoftlayerCommand) (cmd *RuleEditCommand) {
	thisCmd := &RuleEditCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "rule-edit " + T("SECURITYGROUP_ID") + " " + T("RULE_ID"),
		Short: T("Edit a security group rule in a security group"),
		Args:  metadata.TwoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVarP(&thisCmd.RemoteIp, "remote-ip", "r", "", T("The remote IP/CIDR to enforce"))
	cobraCmd.Flags().IntVarP(&thisCmd.RemoteGroup, "remote-group", "s", 0, T("The ID of the remote security group to enforce"))
	cobraCmd.Flags().StringVarP(&thisCmd.Direction, "direction", "d", "", T("The direction of traffic to enforce (ingress or egress), required"))
	cobraCmd.Flags().StringVarP(&thisCmd.EtherType, "ether-type", "e", "", T("The ethertype (IPv4 or IPv6) to enforce, default is IPv4 if not specified"))
	cobraCmd.Flags().IntVarP(&thisCmd.PortMax, "port-max", "M", 0, T("The upper port bound to enforce"))
	cobraCmd.Flags().IntVarP(&thisCmd.PortMin, "port-min", "m", 0, T("The lower port bound to enforce"))
	cobraCmd.Flags().StringVarP(&thisCmd.protocol, "protocol", "p", "", T("The protocol (icmp, tcp, udp) to enforce"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *RuleEditCommand) Run(args []string) error {
	groupID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group ID")
	}
	ruleID, err := strconv.Atoi(args[1])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group rule ID")
	}

	if cmd.Direction != "" {
		direction := cmd.Direction
		if direction != "egress" && direction != "ingress" {
			return errors.NewInvalidUsageError(T("-d|--direction has to be either egress or ingress."))
		}
	}
	if cmd.EtherType != "" {
		etherType := cmd.EtherType
		if etherType != "IPv4" && etherType != "IPv6" {
			return errors.NewInvalidUsageError(T("-e|--ether-type has to be either IPv4 or IPv6."))
		}
	}
	if cmd.protocol != "" {
		protocol := cmd.protocol
		if protocol != "icmp" && protocol != "tcp" && protocol != "udp" {
			return errors.NewInvalidUsageError(T("Options for -p|--protocol are: icmp,tcp,udp"))
		}
	}
	err = cmd.NetworkManager.EditSecurityGroupRule(groupID, ruleID, cmd.RemoteIp, cmd.RemoteGroup, cmd.Direction, cmd.EtherType, cmd.PortMax, cmd.PortMin, cmd.protocol)
	if err != nil {
		return errors.NewAPIError(T("Failed to edit rule {{.RuleId}} in security group {{.GroupID}}.\n",
			map[string]interface{}{"RuleId": ruleID, "GroupID": groupID}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Rule {{.RuleId}} in security group {{.GroupID}} is updated.", map[string]interface{}{"RuleId": ruleID, "GroupID": groupID}))
	return nil
}

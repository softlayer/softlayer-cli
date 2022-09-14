package securitygroup

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type RuleAddCommand struct {
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

func NewRuleAddCommand(sl *metadata.SoftlayerCommand) (cmd *RuleAddCommand) {
	thisCmd := &RuleAddCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "rule-add " + T("SECURITYGROUP_ID"),
		Short: T("Add a security group rule to a security group"),
		Args:  metadata.OneArgs,
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

func (cmd *RuleAddCommand) Run(args []string) error {
	groupID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group ID")
	}
	direction := cmd.Direction
	if direction != "egress" && direction != "ingress" {
		return errors.NewInvalidUsageError(T("-d|--direction has to be either egress or ingress."))
	}
	etherType := "IPv4"
	if cmd.EtherType != "" {
		etherType = cmd.EtherType
	}
	if etherType != "IPv4" && etherType != "IPv6" {
		return errors.NewInvalidUsageError(T("-e|--ether-type has to be either IPv4 or IPv6."))
	}
	portMax := cmd.PortMax
	portMin := cmd.PortMin
	protocol := cmd.protocol
	if portMax > 0 || portMin > 0 {
		if protocol == "" {
			return errors.NewInvalidUsageError(T("-p|--protocal must be set when -M or -m is specified."))
		}
		protocol = strings.ToLower(protocol)
		if protocol != "icmp" && protocol != "tcp" && protocol != "udp" {
			return errors.NewInvalidUsageError(T("Options for -p|--protocol are: icmp,tcp,udp"))
		}
	}

	outputFormat := cmd.GetOutputFlag()

	resp, err := cmd.NetworkManager.AddSecurityGroupRule(groupID, cmd.RemoteIp, cmd.RemoteGroup, direction, etherType, portMax, portMin, protocol)
	if err != nil {
		return errors.NewAPIError(T("Failed to add rule to security group {{.GroupID}}.\n", map[string]interface{}{"GroupID": groupID}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Rule is added to security group {{.GroupID}}.", map[string]interface{}{"GroupID": groupID}))
	return nil
}

package securitygroup

import (
	"sort"
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type RuleListCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	Sortby         string
}

func NewRuleListCommand(sl *metadata.SoftlayerCommand) (cmd *RuleListCommand) {
	thisCmd := &RuleListCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "rule-list " + T("SECURITYGROUP_ID"),
		Short: T("List security group rules"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Sortby, "sortby", "", T("Column to sort by. Options are: id,remoteIp,remoteGroupId,direction,ethertype,portRangeMin,portRangeMax,protocol"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *RuleListCommand) Run(args []string) error {
	groupID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group ID")
	}

	outputFormat := cmd.GetOutputFlag()

	sortColumns := []string{"id", "remoteIp", "remoteGroupId", "direction", "ethertype", "portRangeMin", "portRangeMax", "protocol"}
	sortby := cmd.Sortby
	if sortby != "" && utils.StringInSlice(sortby, sortColumns) == -1 {
		return errors.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.", map[string]interface{}{"Column": sortby}))
	}
	rules, err := cmd.NetworkManager.ListSecurityGroupRules(groupID)
	if err != nil {
		return errors.NewAPIError(T("Failed to get rules of security group {{.GroupID}}.\n", map[string]interface{}{"GroupID": groupID}), err.Error(), 2)
	}

	if sortby == "" || sortby == "id" {
		sort.Sort(utils.RuleById(rules))
	} else if sortby == "remoteIp" {
		sort.Sort(utils.RuleByRemoteIp(rules))
	} else if sortby == "remoteGroupId" {
		sort.Sort(utils.RuleByRemoteGroupId(rules))
	} else if sortby == "direction" {
		sort.Sort(utils.RuleByDirection(rules))
	} else if sortby == "ethertype" {
		sort.Sort(utils.RuleByEtherType(rules))
	} else if sortby == "portRangeMin" {
		sort.Sort(utils.RuleByMinPort(rules))
	} else if sortby == "portRangeMax" {
		sort.Sort(utils.RuleByMaxPort(rules))
	} else if sortby == "protocol" {
		sort.Sort(utils.RuleByProtocol(rules))
	} else {
		return errors.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.", map[string]interface{}{"Column": sortby}))
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, rules)
	}

	if len(rules) == 0 {
		cmd.UI.Print(T("No rules are found for security group {{.GroupID}}.", map[string]interface{}{"GroupID": groupID}))
		return nil
	}

	ruleTable := cmd.UI.Table([]string{T("ID"), T("Remote IP"), T("Remote Group ID"), T("Direction"), T("Ether Type"), T("Port Range Min"), T("Port Range Max"), T("Protocol")})
	for _, rule := range rules {
		ruleTable.Add(utils.FormatIntPointer(rule.Id),
			utils.FormatStringPointer(rule.RemoteIp),
			utils.FormatIntPointer(rule.RemoteGroupId),
			utils.FormatStringPointer(rule.Direction),
			utils.FormatStringPointer(rule.Ethertype),
			utils.FormatIntPointer(rule.PortRangeMin),
			utils.FormatIntPointer(rule.PortRangeMax),
			utils.FormatStringPointer(rule.Protocol),
		)
	}
	ruleTable.Print()
	return nil
}

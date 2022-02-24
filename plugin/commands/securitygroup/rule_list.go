package securitygroup

import (
	"sort"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type RuleListCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewRuleListCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *RuleListCommand) {
	return &RuleListCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *RuleListCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	groupID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	sortColumns := []string{"id", "remoteIp", "remoteGroupId", "direction", "ethertype", "portRangeMin", "portRangeMax", "protocol"}
	sortby := c.String("sortby")
	if sortby != "" && utils.StringInSlice(sortby, sortColumns) == -1 {
		return errors.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.", map[string]interface{}{"Column": sortby}))
	}
	rules, err := cmd.NetworkManager.ListSecurityGroupRules(groupID)
	if err != nil {
		return cli.NewExitError(T("Failed to get rules of security group {{.GroupID}}.\n", map[string]interface{}{"GroupID": groupID})+err.Error(), 2)
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

func SecurityGroupRuleListMetaData() cli.Command {
	return cli.Command{
		Category:    "securitygroup",
		Name:        "rule-list",
		Description: T("List security group rules"),
		Usage:       "${COMMAND_NAME} sl securitygroup rule-list SECURITYGROUP_ID [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by. Options are: id,remoteIp,remoteGroupId,direction,ethertype,portRangeMin,portRangeMax,protocol"),
			},
			metadata.OutputFlag(),
		},
	}
}

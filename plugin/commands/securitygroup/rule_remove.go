package securitygroup

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type RuleRemoveCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewRuleRemoveCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *RuleRemoveCommand) {
	return &RuleRemoveCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *RuleRemoveCommand) Run(c *cli.Context) error {
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
	if !c.IsSet("f") {
		confirm, err := cmd.UI.Confirm(T("This will remove rule {{.RuleId}} in security group {{.GroupId}} and cannot be undone. Continue?",
			map[string]interface{}{"RuleId": ruleID, "GroupId": groupID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.NetworkManager.RemoveSecurityGroupRule(groupID, ruleID)
	if err != nil {
		return cli.NewExitError(T("Failed to remove rule {{.RuleId}} in security group {{.GroupID}}.\n",
			map[string]interface{}{"RuleId": ruleID, "GroupID": groupID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Rule {{.RuleId}} in security group {{.GroupID}} is removed.", map[string]interface{}{"RuleId": ruleID, "GroupID": groupID}))
	return nil
}

func SecurityGroupRuleRemoveMetaData() cli.Command {
	return cli.Command{
		Category:    "securitygroup",
		Name:        "rule-remove",
		Description: T("Remove a rule from a security group"),
		Usage:       "${COMMAND_NAME} sl securitygroup rule-remove SECURITYGROUP_ID RULE_ID [OPTIONS]",
		Flags: []cli.Flag{
			metadata.ForceFlag(),
		},
	}
}

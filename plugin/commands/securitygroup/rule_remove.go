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

type RuleRemoveCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	ForceFlag      bool
}

func NewRuleRemoveCommand(sl *metadata.SoftlayerCommand) (cmd *RuleRemoveCommand) {
	thisCmd := &RuleRemoveCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "rule-remove " + T("SECURITYGROUP_ID") + " " + T("RULE_ID"),
		Short: T("Remove a rule from a security group"),
		Args:  metadata.TwoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVarP(&thisCmd.ForceFlag, "force", "f", false, T("Force operation without confirmation"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *RuleRemoveCommand) Run(args []string) error {
	groupID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group ID")
	}
	ruleID, err := strconv.Atoi(args[1])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group rule ID")
	}
	if !cmd.ForceFlag {
		confirm, err := cmd.UI.Confirm(T("This will remove rule {{.RuleId}} in security group {{.GroupId}} and cannot be undone. Continue?",
			map[string]interface{}{"RuleId": ruleID, "GroupId": groupID}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.NetworkManager.RemoveSecurityGroupRule(groupID, ruleID)
	if err != nil {
		return errors.NewAPIError(T("Failed to remove rule {{.RuleId}} in security group {{.GroupID}}.\n",
			map[string]interface{}{"RuleId": ruleID, "GroupID": groupID}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Rule {{.RuleId}} in security group {{.GroupID}} is removed.", map[string]interface{}{"RuleId": ruleID, "GroupID": groupID}))
	return nil
}

package loadbal

import (
	"github.com/spf13/cobra"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type L7RuleDelCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	PolicyUuid          string
	RuleUuid            string
	Force               bool
}

func NewL7RuleDelCommand(sl *metadata.SoftlayerCommand) *L7RuleDelCommand {
	thisCmd := &L7RuleDelCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "l7rule-delete",
		Short: T("Delete a L7 rule."),
		Long:  T("${COMMAND_NAME} sl loadbal l7rule-delete (--policy-uuid L7POLICY_UUID) (--rule-uuid L7RULE_UUID) [-f, --force]"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.PolicyUuid, "policy-uuid", "", T("UUID for the load balancer policy [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.RuleUuid, "rule-uuid", "", T("UUID for the load balancer rule [required]"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *L7RuleDelCommand) Run(args []string) error {
	l7PolicyID := cmd.PolicyUuid
	if l7PolicyID == "" {
		return errors.NewMissingInputError("--policy-uuid")
	}

	l7RuleID := cmd.RuleUuid
	if l7RuleID == "" {
		return errors.NewMissingInputError("--rule-uuid")
	}

	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will delete the load balancer L7 rule: {{.RuleID}} and cannot be undone. Continue?", map[string]interface{}{"RuleID": l7RuleID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Say(T("Aborted."))
			return nil
		}
	}

	_, err := cmd.LoadBalancerManager.DeleteL7Rule(&l7PolicyID, l7RuleID)
	if err != nil {
		return cli.NewExitError(T("Failed to delete L7Rule {{.L7RuleID}}: {{.Error}}.\n",
			map[string]interface{}{"L7RuleID": l7RuleID, "Error": err.Error()}), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("L7Rule {{.L7RuleID}} removed", map[string]interface{}{"L7RuleID": l7RuleID}))
	return nil
}

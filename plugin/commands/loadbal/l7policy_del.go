package loadbal

import (
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type L7PolicyDeleteCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	PolicyId            int
	Force               bool
}

func NewL7PolicyDeleteCommand(sl *metadata.SoftlayerCommand) *L7PolicyDeleteCommand {
	thisCmd := &L7PolicyDeleteCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "l7policy-delete",
		Short: T("Delete a L7 policy"),
		Long:  T("${COMMAND_NAME} sl loadbal l7policy-delete (--policy-id POLICY_ID) [-f, --force]"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.PolicyId, "policy-id", 0, T("ID for the load balancer policy [required]"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *L7PolicyDeleteCommand) Run(args []string) error {
	policyID := cmd.PolicyId
	if policyID == 0 {
		return errors.NewMissingInputError("--policy-id")
	}

	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will cancel the load balancer policy: {{.PolicyID}} and cannot be undone. Continue?", map[string]interface{}{"PolicyID": policyID}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Say(T("Aborted."))
			return nil
		}
	}

	_, err := cmd.LoadBalancerManager.DeleteL7Policy(policyID)
	if err != nil {
		return errors.New(T("Failed to delete l7 policy: {{.Error}}.\n", map[string]interface{}{"Error": err.Error()}))
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("L7 policy deleted"))
	return nil
}

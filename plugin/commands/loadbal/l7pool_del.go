package loadbal

import (
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type L7PoolDelCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	PoolId              int
	Force               bool
}

func NewL7PoolDelCommand(sl *metadata.SoftlayerCommand) *L7PoolDelCommand {
	thisCmd := &L7PoolDelCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "l7pool-delete",
		Short: T("Delete a L7 pool"),
		Long:  T("${COMMAND_NAME} sl loadbal l7pool-delete (--pool-id L7POOL_ID)"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.PoolId, "pool-id", 0, T("ID for the load balancer pool [required]"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *L7PoolDelCommand) Run(args []string) error {
	l7PoolID := cmd.PoolId
	if l7PoolID == 0 {
		return errors.NewMissingInputError("--pool-id")
	}

	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will delete the load balancer L7 pool: {{.PoolID}} and cannot be undone. Continue?", map[string]interface{}{"PoolID": l7PoolID}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Say(T("Aborted."))
			return nil
		}
	}

	_, err := cmd.LoadBalancerManager.DeleteLoadBalancerL7Pool(l7PoolID)
	if err != nil {
		return errors.New(T("Failed to delete L7Pool {{.L7PoolID}}: {{.Error}}.\n",
			map[string]interface{}{"L7PoolID": l7PoolID, "Error": err.Error()}))
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("L7Pool {{.L7PoolID}} removed", map[string]interface{}{"L7PoolID": l7PoolID}))
	return nil
}

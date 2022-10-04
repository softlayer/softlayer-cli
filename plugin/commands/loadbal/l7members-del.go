package loadbal

import (
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type L7MembersDelCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	PoolUuid            string
	MemberUuid          string
	Force               bool
}

func NewL7MembersDelCommand(sl *metadata.SoftlayerCommand) *L7MembersDelCommand {
	thisCmd := &L7MembersDelCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "l7member-delete",
		Short: T("Remove a load balancer member"),
		Long:  T("${COMMAND_NAME} sl loadbal l7member-del (--pool-uuid L7POOL_UUID) (--member-uuid L7MEMBER_UUID)"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.PoolUuid, "pool-uuid", "", T("UUID for the load balancer pool [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.MemberUuid, "member-uuid", "", T("UUID for the load balancer member [required]"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *L7MembersDelCommand) Run(args []string) error {
	L7POOL_UUID := cmd.PoolUuid
	if L7POOL_UUID == "" {
		return errors.NewMissingInputError("--pool-uuid")
	}

	memberUUID := cmd.MemberUuid
	if memberUUID == "" {
		return errors.NewMissingInputError("--member-uuid")
	}

	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will delete the load balancer L7 member: {{.MemberID}} and cannot be undone. Continue?", map[string]interface{}{"MemberID": memberUUID}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Say(T("Aborted."))
			return nil
		}
	}

	_, err := cmd.LoadBalancerManager.DeleteL7Member(&L7POOL_UUID, memberUUID)
	if err != nil {
		return errors.New(T("Failed to delete L7member {{.Member}}: {{.Error}}.\n",
			map[string]interface{}{"Member": memberUUID, "Error": err.Error()}))
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("Member {{.MemberID}} removed from {{.L7POOL}}", map[string]interface{}{"MemberID": memberUUID, "L7POOL": L7POOL_UUID}))
	return nil
}

package loadbal

import (
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type MembersDelCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	LbId                int
	MemberUuid          string
	Force               bool
}

func NewMembersDelCommand(sl *metadata.SoftlayerCommand) *MembersDelCommand {
	thisCmd := &MembersDelCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "member-delete",
		Short: T("Remove a load balancer member"),
		Long:  T("${COMMAND_NAME} sl loadbal member-del (--lb-id LOADBAL_ID) (-m, --member-uuid MEMBER_UUID)"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.LbId, "lb-id", 0, T("ID for the load balancer [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.MemberUuid, "member-uuid", "m", "", T("Member UUID [required]"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *MembersDelCommand) Run(args []string) error {
	loadbalID := cmd.LbId
	if loadbalID == 0 {
		return errors.NewMissingInputError("--lb-id")
	}

	memberUUID := cmd.MemberUuid
	if memberUUID == "" {
		return errors.NewMissingInputError("-m, --member-uuid")
	}

	loadbalancerUUID, err := cmd.LoadBalancerManager.GetLoadBalancerUUID(loadbalID)
	if err != nil {
		return errors.New(T("Failed to get load balancer: {{.ERR}}.", map[string]interface{}{"ERR": err.Error()}))
	}

	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will delete the load balancer member: {{.MemberID}} and cannot be undone. Continue?", map[string]interface{}{"MemberID": memberUUID}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Say(T("Aborted."))
			return nil
		}
	}

	_, err = cmd.LoadBalancerManager.DeleteLoadBalancerMember(&loadbalancerUUID, []string{memberUUID})
	if err != nil {
		return errors.New(T("Failed to delete load balancer member {{.Member}}: {{.Error}}.\n",
			map[string]interface{}{"Member": memberUUID, "Error": err.Error()}))
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("Member {{.MemberID}} removed", map[string]interface{}{"MemberID": memberUUID}))
	return nil
}

package loadbal

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type MembersAddCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	Id                  int
	Ip                  string
}

func NewMembersAddCommand(sl *metadata.SoftlayerCommand) *MembersAddCommand {
	thisCmd := &MembersAddCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "member-add",
		Short: T("Add a new load balancer member"),
		Long:  T("${COMMAND_NAME} sl loadbal member-add (--id LOADBAL_ID) (--ip PRIVATE_IP)"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.Id, "id", 0, T("ID for the load balancer [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.Ip, "ip", "", T("Private IP of the new member [required]"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *MembersAddCommand) Run(args []string) error {
	loadbalID := cmd.Id
	if loadbalID == 0 {
		return errors.NewMissingInputError("--id")
	}

	if cmd.Ip == "" {
		return errors.NewMissingInputError("--ip")
	}

	loadbalancerUUID, err := cmd.LoadBalancerManager.GetLoadBalancerUUID(loadbalID)
	if err != nil {
		return cli.NewExitError(T("Failed to get load balancer: {{.ERR}}.", map[string]interface{}{"ERR": err.Error()}), 2)
	}

	ip := cmd.Ip
	toAdd := datatypes.Network_LBaaS_LoadBalancerServerInstanceInfo{
		PrivateIpAddress: &ip,
	}

	_, err = cmd.LoadBalancerManager.AddLoadBalancerMember(&loadbalancerUUID, []datatypes.Network_LBaaS_LoadBalancerServerInstanceInfo{toAdd})
	if err != nil {
		return cli.NewExitError(T("Failed to add load balancer member: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("Member {{.MemberID}} added", map[string]interface{}{"MemberID": ip}))
	return nil
}

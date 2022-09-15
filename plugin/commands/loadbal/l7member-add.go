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

type L7MembersAddCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	PoolUuid            string
	Address             string
	Port                int
}

func NewL7MembersAddCommand(sl *metadata.SoftlayerCommand) *L7MembersAddCommand {
	thisCmd := &L7MembersAddCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "l7member-add",
		Short: T("Add a new L7 pool member."),
		Long:  T("${COMMAND_NAME} sl loadbal member-add (--pool-uuid L7POOL_UUID) (--address IP_ADDRESS) (--port PORT)"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.PoolUuid, "pool-uuid", "", T("UUID for the load balancer pool [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.Address, "address", "", T("Backend servers IP address. [required]"))
	cobraCmd.Flags().IntVarP(&thisCmd.Port, "port", "t", 0, T("Backend servers port. [required]"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *L7MembersAddCommand) Run(args []string) error {
	poolUUID := cmd.PoolUuid
	if poolUUID == "" {
		return errors.NewMissingInputError("--pool-uuid")
	}

	toAdd := datatypes.Network_LBaaS_L7Member{}

	ip := cmd.Address
	if ip == "" {
		return errors.NewMissingInputError("--address")
	}
	toAdd.Address = &ip

	port := cmd.Port
	if port == 0 {
		return errors.NewMissingInputError("--port")
	}
	toAdd.Port = &port

	_, err := cmd.LoadBalancerManager.AddL7Member(&poolUUID, []datatypes.Network_LBaaS_L7Member{toAdd})
	if err != nil {
		return cli.NewExitError(T("Failed to add L7 member: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("L7 Member {{.MemberID}} added in pool {{.Pool}}", map[string]interface{}{"MemberID": ip, "Pool": poolUUID}))
	return nil
}

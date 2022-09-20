package loadbal

import (
	"github.com/spf13/cobra"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type ProtocolDeleteCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	LbId                int
	ProtocolUuid        string
	Force               bool
}

func NewProtocolDeleteCommand(sl *metadata.SoftlayerCommand) *ProtocolDeleteCommand {
	thisCmd := &ProtocolDeleteCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "protocol-delete",
		Short: T("Delete a protocol."),
		Long:  T("${COMMAND_NAME} sl loadbal protocol-delete (--lb-id LOADBAL_ID) (--protocol-uuid PROTOCOL_UUID)"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.LbId, "lb-id", 0, T("ID for the load balancer [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.ProtocolUuid, "protocol-uuid", "", T("UUID for the protocol [required]"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ProtocolDeleteCommand) Run(args []string) error {
	loadbalID := cmd.LbId
	if loadbalID == 0 {
		return errors.NewMissingInputError("--lb-id")
	}

	protocolUUID := cmd.ProtocolUuid
	if protocolUUID == "" {
		return errors.NewMissingInputError("--protocol-uuid")
	}

	loadbalancerUUID, err := cmd.LoadBalancerManager.GetLoadBalancerUUID(loadbalID)
	if err != nil {
		return cli.NewExitError(T("Failed to get load balancer: {{.ERR}}.", map[string]interface{}{"ERR": err.Error()}), 2)
	}

	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will delete the load balancer protocol: {{.ProtocolID}} and cannot be undone. Continue?", map[string]interface{}{"ProtocolID": protocolUUID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Say(T("Aborted."))
			return nil
		}
	}

	_, err = cmd.LoadBalancerManager.DeleteLoadBalancerListener(&loadbalancerUUID, []string{protocolUUID})
	if err != nil {
		return cli.NewExitError(T("Failed to delete protocol {{.ProtocolID}}: {{.Error}}.\n",
			map[string]interface{}{"ProtocolID": protocolUUID, "Error": err.Error()}), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("Protocol {{.ProtocolID}} removed", map[string]interface{}{"ProtocolID": protocolUUID}))
	return nil
}

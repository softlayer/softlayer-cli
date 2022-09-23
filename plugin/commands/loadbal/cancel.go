package loadbal

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CancelCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	Force               bool
}

func NewCancelCommand(sl *metadata.SoftlayerCommand) *CancelCommand {
	thisCmd := &CancelCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "cancel " + T("IDENTIFIER"),
		Short: T("Cancel an existing load balancer"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CancelCommand) Run(args []string) error {
	loadbalID, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("LoadBalancer ID")
	}

	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will cancel the load balancer: {{.LBID}} and cannot be undone. Continue?", map[string]interface{}{"LBID": loadbalID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	loadbalUUID, err := cmd.LoadBalancerManager.GetLoadBalancerUUID(loadbalID)
	if err != nil {
		return cli.NewExitError(T("Failed to get load balancer: {{.ERR}}.", map[string]interface{}{"ERR": err.Error()}), 2)
	}

	_, err = cmd.LoadBalancerManager.CancelLoadBalancer(&loadbalUUID)
	if err != nil {
		return errors.NewAPIError(T("Failed to cancel load balancer {{.LBID}}.\n", map[string]interface{}{"LBID": loadbalID}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Load balancer {{.LBID}} is cancelled.", map[string]interface{}{"LBID": loadbalID}))
	return nil
}

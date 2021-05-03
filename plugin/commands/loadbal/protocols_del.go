package loadbal

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
)

type ProtocolDeleteCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewProtocolDeleteCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *ProtocolDeleteCommand) {
	return &ProtocolDeleteCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *ProtocolDeleteCommand) Run(c *cli.Context) error {
	loadbalID := c.Int("lb-id")
	if loadbalID == 0 {
		return errors.NewMissingInputError("--lb-id")
	}

	protocolUUID := c.String("protocol-uuid")
	if protocolUUID == "" {
		return errors.NewMissingInputError("--protocol-uuid")
	}

	loadbalancerUUID, err := cmd.LoadBalancerManager.GetLoadBalancerUUID(loadbalID)
	if err != nil {
		return cli.NewExitError(T("Failed to get load balancer: {{.ERR}}.", map[string]interface{}{"ERR": err.Error()}), 2)
	}

	if !c.IsSet("f") {
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

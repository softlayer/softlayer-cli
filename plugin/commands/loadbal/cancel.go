package loadbal

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CancelCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewCancelCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *CancelCommand) {
	return &CancelCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *CancelCommand) Run(c *cli.Context) error {
	loadbalID := c.Int("id")
	if loadbalID == 0 {
		return errors.NewMissingInputError("--id")
	}

	if !c.IsSet("f") {
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
		return cli.NewExitError(T("Failed to cancel load balancer {{.LBID}}.\n", map[string]interface{}{"LBID": loadbalID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Load balancer {{.LBID}} is cancelled.", map[string]interface{}{"LBID": loadbalID}))
	return nil
}

func LoadbalCancelMetadata() cli.Command {
	return cli.Command{
		Category:    "loadbal",
		Name:        "cancel",
		Description: T("Cancel an existing load balancer"),
		Usage:       "${COMMAND_NAME} sl loadbal cancel (--id LOADBAL_ID)",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "id",
				Usage: T("ID for the load balancer [required]"),
			},
			metadata.ForceFlag(),
		},
	}
}

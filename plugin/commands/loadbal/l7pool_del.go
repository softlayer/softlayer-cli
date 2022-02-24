package loadbal

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type L7PoolDelCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewL7PoolDelCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *L7PoolDelCommand) {
	return &L7PoolDelCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *L7PoolDelCommand) Run(c *cli.Context) error {

	l7PoolID := c.Int("pool-id")
	if l7PoolID == 0 {
		return errors.NewMissingInputError("--pool-id")
	}

	if !c.IsSet("f") {
		confirm, err := cmd.UI.Confirm(T("This will delete the load balancer L7 pool: {{.PoolID}} and cannot be undone. Continue?", map[string]interface{}{"PoolID": l7PoolID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Say(T("Aborted."))
			return nil
		}
	}

	_, err := cmd.LoadBalancerManager.DeleteLoadBalancerL7Pool(l7PoolID)
	if err != nil {
		return cli.NewExitError(T("Failed to delete L7Pool {{.L7PoolID}}: {{.Error}}.\n",
			map[string]interface{}{"L7PoolID": l7PoolID, "Error": err.Error()}), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("L7Pool {{.L7PoolID}} removed", map[string]interface{}{"L7PoolID": l7PoolID}))
	return nil
}

func LoadbalL7PoolDelMetadata() cli.Command {
	return cli.Command{
		Category:    "loadbal",
		Name:        "l7pool-delete",
		Description: T("Delete a L7 pool"),
		Usage:       "${COMMAND_NAME} sl loadbal l7pool-delete (--pool-id L7POOL_ID)",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "pool-id",
				Usage: T("ID for the load balancer pool [required]"),
			},
		},
	}
}

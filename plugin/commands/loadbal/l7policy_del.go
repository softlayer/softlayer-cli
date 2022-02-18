package loadbal

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type L7PolicyDeleteCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewL7PolicyDeleteCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *L7PolicyDeleteCommand) {
	return &L7PolicyDeleteCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *L7PolicyDeleteCommand) Run(c *cli.Context) error {
	policyID := c.Int("policy-id")
	if policyID == 0 {
		return errors.NewMissingInputError("--policy-id")
	}

	if !c.IsSet("f") {
		confirm, err := cmd.UI.Confirm(T("This will cancel the load balancer policy: {{.PolicyID}} and cannot be undone. Continue?", map[string]interface{}{"PolicyID": policyID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Say(T("Aborted."))
			return nil
		}
	}

	_, err := cmd.LoadBalancerManager.DeleteL7Policy(policyID)
	if err != nil {
		return cli.NewExitError(T("Failed to delete l7 policy: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("L7 policy deleted"))
	return nil
}

func LoadbalL7PolicyDeleteMetadata() cli.Command {
	return cli.Command{
		Category:    "loadbal",
		Name:        "l7policy-delete",
		Description: T("Delete a L7 policy"),
		Usage:       "${COMMAND_NAME} sl loadbal l7policy-delete (--policy-id POLICY_ID) [-f, --force]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "policy-id",
				Usage: T("ID for the load balancer policy [required]"),
			},
			metadata.ForceFlag(),
		},
	}
}

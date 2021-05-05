package loadbal

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type L7MembersDelCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewL7MembersDelCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *L7MembersDelCommand) {
	return &L7MembersDelCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *L7MembersDelCommand) Run(c *cli.Context) error {
	L7POOL_UUID := c.String("pool-uuid")
	if L7POOL_UUID == "" {
		return errors.NewMissingInputError("--pool-uuid")
	}

	memberUUID := c.String("member-uuid")
	if memberUUID == "" {
		return errors.NewMissingInputError("--member-uuid")
	}

	if !c.IsSet("f") {
		confirm, err := cmd.UI.Confirm(T("This will delete the load balancer L7 member: {{.MemberID}} and cannot be undone. Continue?", map[string]interface{}{"MemberID": memberUUID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Say(T("Aborted."))
			return nil
		}
	}

	_, err := cmd.LoadBalancerManager.DeleteL7Member(&L7POOL_UUID, memberUUID)
	if err != nil {
		return cli.NewExitError(T("Failed to delete L7member {{.Member}}: {{.Error}}.\n",
			map[string]interface{}{"Member": memberUUID, "Error": err.Error()}), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("Member {{.MemberID}} removed from {{.L7POOL}}", map[string]interface{}{"MemberID": memberUUID, "L7POOL": L7POOL_UUID}))
	return nil
}

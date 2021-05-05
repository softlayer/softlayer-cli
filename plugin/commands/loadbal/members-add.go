package loadbal

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type MembersAddCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewMembersAddCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *MembersAddCommand) {
	return &MembersAddCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *MembersAddCommand) Run(c *cli.Context) error {
	loadbalID := c.Int("id")
	if loadbalID == 0 {
		return errors.NewMissingInputError("--id")
	}

	if !c.IsSet("ip") {
		return errors.NewMissingInputError("--ip")
	}

	loadbalancerUUID, err := cmd.LoadBalancerManager.GetLoadBalancerUUID(loadbalID)
	if err != nil {
		return cli.NewExitError(T("Failed to get load balancer: {{.ERR}}.", map[string]interface{}{"ERR": err.Error()}), 2)
	}

	ip := c.String("ip")
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

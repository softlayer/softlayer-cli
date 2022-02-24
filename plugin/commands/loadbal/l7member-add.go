package loadbal

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type L7MembersAddCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewL7MembersAddCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *L7MembersAddCommand) {
	return &L7MembersAddCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *L7MembersAddCommand) Run(c *cli.Context) error {
	poolUUID := c.String("pool-uuid")
	if poolUUID == "" {
		return errors.NewMissingInputError("--pool-uuid")
	}

	toAdd := datatypes.Network_LBaaS_L7Member{}

	ip := c.String("address")
	if ip == "" {
		return errors.NewMissingInputError("--address")
	}
	toAdd.Address = &ip

	port := c.Int("port")
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

func LoadbalMemberAddMetadata() cli.Command {
	return cli.Command{
		Category:    "loadbal",
		Name:        "member-add",
		Description: T("Add a new load balancer member"),
		Usage:       "${COMMAND_NAME} sl loadbal member-add (--id LOADBAL_ID) (--ip PRIVATE_IP)",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "id",
				Usage: T("ID for the load balancer [required]"),
			},
			cli.StringFlag{
				Name:  "ip",
				Usage: T("Private IP of the new member [required]"),
			},
		},
	}
}

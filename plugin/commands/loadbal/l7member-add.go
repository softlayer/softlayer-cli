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

func LoadbalL7MemberAddMetadata() cli.Command {
	return cli.Command{
		Category:    "loadbal",
		Name:        "l7member-add",
		Description: T("Add a new L7 pool member"),
		Usage:       "${COMMAND_NAME} sl loadbal member-add (--pool-uuid L7POOL_UUID) (--address IP_ADDRESS) (--port PORT)",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "pool-uuid",
				Usage: T("UUID for the load balancer pool [required]"),
			},
			cli.StringFlag{
				Name:  "address",
				Usage: T("Backend servers IP address. [required]"),
			},
			cli.IntFlag{
				Name:  "port",
				Usage: T("Backend servers port. [required]"),
			},
		},
	}
}

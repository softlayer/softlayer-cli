package loadbal


import (
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)


func GetCommandAcionBindings(ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	lbManager := managers.NewLoadBalancerManager(session)
	networkManager := managers.NewNetworkManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{

		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_CANCLE_NAME: func(c *cli.Context) error {
			return NewCancelCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_ORDER_NAME: func(c *cli.Context) error {
			return NewCreateCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_ORDER_OPTIONS_NAME: func(c *cli.Context) error {
			return NewOptionsCommand(ui, lbManager, networkManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_DETAIL_NAME: func(c *cli.Context) error {
			return NewDetailCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_LIST_NAME: func(c *cli.Context) error {
			return NewListCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_HEALTH_NAME: func(c *cli.Context) error {
			return NewHealthChecksCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_PROTOCOL_ADD_NAME: func(c *cli.Context) error {
			return NewProtocolAddCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_PROTOCOL_DELETE_NAME: func(c *cli.Context) error {
			return NewProtocolDeleteCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-protocol-edit": func(c *cli.Context) error {
			return NewProtocolEditCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_MEMBER_ADD_NAME: func(c *cli.Context) error {
			return NewMembersAddCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_MEMBER_DEL_NAME: func(c *cli.Context) error {
			return NewMembersDelCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_L7POOL_ADD_NAME: func(c *cli.Context) error {
			return NewL7PoolAddCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_L7POOL_DELETE_NAME: func(c *cli.Context) error {
			return NewL7PoolDelCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_L7POOL_DETAIL_NAME: func(c *cli.Context) error {
			return NewL7PoolDetailCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_L7POOL_EDIT_NAME: func(c *cli.Context) error {
			return NewL7PoolEditCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_L7MEMBER_ADD_NAME: func(c *cli.Context) error {
			return NewL7MembersAddCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_L7MEMBER_DELETE_NAME: func(c *cli.Context) error {
			return NewL7MembersDelCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_L7POLICY_ADD_NAME: func(c *cli.Context) error {
			return NewL7PolicyAddCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_L7POLICY_DELETE_NAME: func(c *cli.Context) error {
			return NewL7PolicyDeleteCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_L7POLICY_LIST_NAME: func(c *cli.Context) error {
			return NewL7PolicyListCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_L7RULE_ADD_NAME: func(c *cli.Context) error {
			return NewL7RuleAddCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_L7RULE_DELETE_NAME: func(c *cli.Context) error {
			return NewL7RuleDelCommand(ui, lbManager).Run(c)
		},
		NS_LOADBAL_NAME + "-" + CMD_LOADBAL_L7RULE_LIST_NAME: func(c *cli.Context) error {
			return NewL7RuleListCommand(ui, lbManager).Run(c)
		},
	}

	return CommandActionBindings
}
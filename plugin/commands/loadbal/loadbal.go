package loadbal

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {

	loadbalManager := managers.NewLoadBalancerManager(session)
	networkManager := managers.NewNetworkManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"loadbal-cancel": func(c *cli.Context) error {
			return NewCancelCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-detail": func(c *cli.Context) error {
			return NewDetailCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-health-edit": func(c *cli.Context) error {
			return NewHealthChecksCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-l7member-add": func(c *cli.Context) error {
			return NewMembersAddCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-l7member-delete": func(c *cli.Context) error {
			return NewMembersDelCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-l7policies": func(c *cli.Context) error {
			return NewL7PolicyListCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-l7policy-add": func(c *cli.Context) error {
			return NewL7PolicyAddCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-l7policy-delete": func(c *cli.Context) error {
			return NewL7PolicyDeleteCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-l7policy-edit": func(c *cli.Context) error {
			return NewL7PolicyEditCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-l7pool-add": func(c *cli.Context) error {
			return NewL7PoolAddCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-l7pool-delete": func(c *cli.Context) error {
			return NewL7PoolDelCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-l7pool-detail": func(c *cli.Context) error {
			return NewL7PoolDetailCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-l7pool-edit": func(c *cli.Context) error {
			return NewL7PoolEditCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-l7rule-add": func(c *cli.Context) error {
			return NewL7RuleAddCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-l7rule-delete": func(c *cli.Context) error {
			return NewL7RuleDelCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-l7rules": func(c *cli.Context) error {
			return NewL7RuleListCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-list": func(c *cli.Context) error {
			return NewListCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-member-add": func(c *cli.Context) error {
			return NewMembersAddCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-member-delete": func(c *cli.Context) error {
			return NewMembersDelCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-order": func(c *cli.Context) error {
			return NewCreateCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-order-options": func(c *cli.Context) error {
			return NewOptionsCommand(ui, loadbalManager, networkManager).Run(c)
		},
		"loadbal-protocol-add": func(c *cli.Context) error {
			return NewProtocolAddCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-protocol-delete": func(c *cli.Context) error {
			return NewProtocolDeleteCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-protocol-edit": func(c *cli.Context) error {
			return NewProtocolEditCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-ns-detail": func(c *cli.Context) error {
			return NewNetscalerDetailCommand(ui, loadbalManager).Run(c)
		},
		"loadbal-ns-list": func(c *cli.Context) error {
			return NewNetscalerListCommand(ui, loadbalManager).Run(c)
		},
	}

	return CommandActionBindings
}

func LoadbalNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "loadbal",
		Description: T("Classic infrastructure Load Balancers"),
	}
}

func LoadbalMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "loadbal",
		Description: T("Classic infrastructure Load Balancers"),
		Usage:       "${COMMAND_NAME} sl loadbal",
		Subcommands: []cli.Command{
			LoadbalCancelMetadata(),
			LoadbalOrderMetadata(),
			LoadbalOrderOptionsMetadata(),
			LoadbalDetailMetadata(),
			LoadbalListMetadata(),
			LoadbalHealthMetadata(),
			LoadbalMemberAddMetadata(),
			LoadbalMemberDelMetadata(),
			LoadbalProtocolAddMetadata(),
			LoadbalProtocolDelMetadata(),
			LoadbalProtocolEditMetadata(),
			LoadbalL7PoolDelMetadata(),
			LoadbalL7PoolAddMetadata(),
			LoadbalL7PoolDetailMetadata(),
			LoadbalL7PoolEditMetadata(),
			LoadbalL7MemberAddMetadata(),
			LoadbalL7MemberDeleteMetadata(),
			LoadbalL7PolicyAddMetadata(),
			LoadbalL7PolicyEditMetadata(),
			LoadbalL7PolicyDeleteMetadata(),
			LoadbalL7PolicyListMetadata(),
			LoadbalL7RuleAddMetadata(),
			LoadbalL7RuleDelMetadata(),
			LoadbalL7RuleListMetadata(),
			LoadbalNetscalerDetailMetadata(),
			LoadbalNsListMetadata(),
		},
	}
}

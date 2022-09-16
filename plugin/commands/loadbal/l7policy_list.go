package loadbal

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type L7PolicyListCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	ProtocolId          int
}

func NewL7PolicyListCommand(sl *metadata.SoftlayerCommand) *L7PolicyListCommand {
	thisCmd := &L7PolicyListCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "l7policies",
		Short: T("List L7 policies."),
		Long:  T("${COMMAND_NAME} sl loadbal l7policies (--protocol-id PROTOCOL_ID)"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.ProtocolId, "protocol-id", 0, T("ID for the load balancer protocol [required]"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *L7PolicyListCommand) Run(args []string) error {
	protocolID := cmd.ProtocolId
	if protocolID == 0 {
		return errors.NewMissingInputError("--protocol-id")
	}

	l7Policies, err := cmd.LoadBalancerManager.GetL7Policies(protocolID)
	if err != nil {
		return cli.NewExitError(T("Failed to get l7 policies: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}
	PrintPolicies(l7Policies, cmd.UI)

	return nil
}

func PrintPolicies(l7Policies []datatypes.Network_LBaaS_L7Policy, cmdUI terminal.UI) {

	if len(l7Policies) == 0 {
		cmdUI.Say(T("No l7 policies was found."))
	} else {
		table := cmdUI.Table([]string{T("ID"), T("UUID"), T("Name"), T("Action"), T("Redirect"), T("Priority"), T("Create Date")})
		for _, l7Policy := range l7Policies {
			if l7Policy.Action != nil && (*l7Policy.Action == REDIRECT_URL || *l7Policy.Action == REDIRECT_HTTPS) {
				table.Add(utils.FormatIntPointer(l7Policy.Id),
					utils.FormatStringPointer(l7Policy.Uuid),
					utils.FormatStringPointer(l7Policy.Name),
					utils.FormatStringPointer(l7Policy.Action),
					utils.FormatStringPointer(l7Policy.RedirectUrl),
					utils.FormatIntPointer(l7Policy.Priority),
					utils.FormatSLTimePointer(l7Policy.CreateDate),
				)
			}
			if l7Policy.Action != nil && *l7Policy.Action == REDIRECT_POOL {
				table.Add(utils.FormatIntPointer(l7Policy.Id),
					utils.FormatStringPointer(l7Policy.Uuid),
					utils.FormatStringPointer(l7Policy.Name),
					utils.FormatStringPointer(l7Policy.Action),
					utils.FormatIntPointer(l7Policy.RedirectL7PoolId),
					utils.FormatIntPointer(l7Policy.Priority),
					utils.FormatSLTimePointer(l7Policy.CreateDate),
				)
			}
			if l7Policy.Action != nil && *l7Policy.Action == REJECT {
				table.Add(utils.FormatIntPointer(l7Policy.Id),
					utils.FormatStringPointer(l7Policy.Uuid),
					utils.FormatStringPointer(l7Policy.Name),
					utils.FormatStringPointer(l7Policy.Action),
					"-",
					utils.FormatIntPointer(l7Policy.Priority),
					utils.FormatSLTimePointer(l7Policy.CreateDate),
				)
			}
		}
		table.Print()
	}
}

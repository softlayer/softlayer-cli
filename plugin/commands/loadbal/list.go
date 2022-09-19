package loadbal

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
}

func NewListCommand(sl *metadata.SoftlayerCommand) *ListCommand {
	thisCmd := &ListCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List active load balancers."),
		Long:  T("${COMMAND_NAME} sl loadbal list"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ListCommand) Run(args []string) error {
	loadbalancers, err := cmd.LoadBalancerManager.GetLoadBalancers()
	if err != nil {
		return errors.NewAPIError(T("Failed to get load balancers on your account."), err.Error(), 2)
	}
	if len(loadbalancers) == 0 {
		cmd.UI.Print(T("No load balancer was found."))
	} else {
		table := cmd.UI.Table([]string{T("ID"), T("Name"), T("Address"), T("Type"), T("Location"), T("Create Date"), T("Status")})
		for _, lb := range loadbalancers {
			var location, lbType string
			if lb.Datacenter != nil {
				location = utils.FormatStringPointer(lb.Datacenter.LongName)
			}
			if utils.FormatIntPointer(lb.Type) == "1" {
				lbType = "Public to Private"
			} else if utils.FormatIntPointer(lb.Type) == "0" {
				lbType = "Private to Private"
			} else if utils.FormatIntPointer(lb.Type) == "2" {
				lbType = "Public to Public"
			}

			table.Add(utils.FormatIntPointer(lb.Id),
				utils.FormatStringPointer(lb.Name),
				utils.FormatStringPointer(lb.Address),
				lbType,
				location,
				utils.FormatSLTimePointer(lb.CreateDate),
				fmt.Sprintf("%s/%s", utils.FormatStringPointer(lb.ProvisioningStatus), utils.FormatStringPointer(lb.OperatingStatus)),
			)
		}
		table.Print()
	}
	return nil
}

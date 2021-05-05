package loadbal

import (
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewListCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *ListCommand) {
	return &ListCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *ListCommand) Run(c *cli.Context) error {
	loadbalancers, err := cmd.LoadBalancerManager.GetLoadBalancers()
	if err != nil {
		return cli.NewExitError(T("Failed to get load balancers on your account.")+err.Error(), 2)
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

package loadbal

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type NetscalerListCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewNetscalerListCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *NetscalerListCommand) {
	return &NetscalerListCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *NetscalerListCommand) Run(c *cli.Context) error {
	netscalers, err := cmd.LoadBalancerManager.GetADCs()
	if err != nil {
		return cli.NewExitError(T("Failed to get netscalers on your account.")+err.Error(), 2)
	}
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}
	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, netscalers)
	}
	if len(netscalers) == 0 {
		cmd.UI.Say(T("No netscalers was found."))
	} else {
		table := cmd.UI.Table([]string{T("ID"), T("Location"), T("Name"), T("Description"), T("IP Address"), T("Management IP"), T("Bandwidth"), T("Create Date")})
		for _, ns := range netscalers {
			var location string
			if ns.Datacenter != nil {
				location = utils.FormatStringPointer(ns.Datacenter.LongName)
			}

			table.Add(utils.FormatIntPointer(ns.Id),
				location,
				utils.FormatStringPointer(ns.Name),
				utils.FormatStringPointer(ns.Description),
				utils.FormatStringPointer(ns.PrimaryIpAddress),
				utils.FormatStringPointer(ns.ManagementIpAddress),
				utils.FormatSLFloatPointerToFloat(ns.OutboundPublicBandwidthUsage),
				utils.FormatSLTimePointer(ns.CreateDate),
			)
		}
		table.Print()
	}
	return nil
}

func LoadbalNsListMetadata() cli.Command {
	return cli.Command{
		Category:    "loadbal",
		Name:        "ns-list",
		Description: T("List netscalers"),
		Usage:       "${COMMAND_NAME} sl loadbal netscalers",
		Flags:       []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

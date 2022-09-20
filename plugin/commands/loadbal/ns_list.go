package loadbal

import (
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type NetscalerListCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
}

func NewNetscalerListCommand(sl *metadata.SoftlayerCommand) *NetscalerListCommand {
	thisCmd := &NetscalerListCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "ns-list",
		Short: T("List netscalers."),
		Long:  T("${COMMAND_NAME} sl loadbal netscalers"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *NetscalerListCommand) Run(args []string) error {
	netscalers, err := cmd.LoadBalancerManager.GetADCs()
	if err != nil {
		return errors.NewAPIError(T("Failed to get netscalers on your account."), err.Error(), 2)
	}

	outputFormat := cmd.GetOutputFlag()

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
		utils.PrintTable(cmd.UI, table, outputFormat)
	}
	return nil
}

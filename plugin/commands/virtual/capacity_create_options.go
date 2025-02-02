package virtual

import (
	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/spf13/cobra"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CapacityCreateOptionsCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
}

func NewCapacityCreateOptionsCommand(sl *metadata.SoftlayerCommand) (cmd *CapacityCreateOptionsCommand) {
	thisCmd := &CapacityCreateOptionsCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "capacity-create-options",
		Short: T("List options for creating Reserved Capacity Group instance"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CapacityCreateOptionsCommand) Run(args []string) error {
	datacenters, _ := cmd.VirtualServerManager.GetRouters("RESERVED_CAPACITY")
	pods, err := cmd.VirtualServerManager.GetPods()
	if err != nil {
		return slErrors.NewAPIError("Failed to get Pods.", err.Error(), 2)
	}
	outputFormat := cmd.GetOutputFlag()
	tableRegion := cmd.UI.Table([]string{T("Location"), T("POD"), T("BackendRouterId")})
	for _, datacenter := range datacenters {
		for _, pod := range pods {
			if utils.FormatStringPointer(datacenter.Location.Location.Name) == utils.FormatStringPointer(pod.DatacenterName) {
				tableRegion.Add(utils.FormatStringPointer(datacenter.Keyname),
					utils.FormatStringPointer(pod.BackendRouterName),
					utils.FormatIntPointer(pod.BackendRouterId))
			}
		}
	}

	tableItems := cmd.UI.Table([]string{T("KeyName"), T("Description"), T("term"), T("Default Hourly Price Per Instance")})
	items, err := cmd.VirtualServerManager.GetCapacityCreateOptions("RESERVED_CAPACITY")
	if err != nil {
		return slErrors.NewInvalidUsageError("Internal error.")
	}
	for _, item := range items {
		tableItems.Add(utils.FormatStringPointer(item.KeyName),
			utils.FormatStringPointer(item.Description),
			utils.FormatSLFloatPointerToFloat(item.Capacity), getPrices(item.Prices))
	}

	utils.PrintTable(cmd.UI, tableItems, outputFormat)
	utils.PrintTable(cmd.UI, tableRegion, outputFormat)
	return nil
}

// Finds the price with the default locationGroupId
func getPrices(prices []datatypes.Product_Item_Price) string {
	itemPrices := ""
	for _, price := range prices {
		if price.LocationGroupId == nil {
			itemPrices = utils.FormatSLFloatPointerToFloat(price.HourlyRecurringFee)
		}
	}
	return itemPrices
}

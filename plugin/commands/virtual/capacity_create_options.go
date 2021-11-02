package virtual

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CapacityCreateOptiosCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewCapacityCreateOptiosCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *CapacityCreateOptiosCommand) {
	return &CapacityCreateOptiosCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *CapacityCreateOptiosCommand) Run(c *cli.Context) error {
	datacenters, _ := cmd.VirtualServerManager.GetRouters("RESERVED_CAPACITY")
	pods, err := cmd.VirtualServerManager.GetPods()
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	tableRegion := cmd.UI.Table([]string{T("Location"), T("POD"), T("BackendRouterId")})
	for _, datacenter := range datacenters{
		for _, pod := range pods{
			if (utils.FormatStringPointer(datacenter.Location.Location.Name) == utils.FormatStringPointer(pod.DatacenterName)){
				tableRegion.Add(utils.FormatStringPointer(datacenter.Keyname),
					utils.FormatStringPointer(pod.BackendRouterName),
					utils.FormatIntPointer(pod.BackendRouterId))
			}
		}
	}
	var capacityCreateOptions []interface{}
	capacityCreateOptions = append(capacityCreateOptions, pods)
	capacityCreateOptions = append(capacityCreateOptions, datacenters)
	if outputFormat == "JSON" {
		return utils.PrintPrettyJSONList(cmd.UI, capacityCreateOptions)
	}
	tableItems := cmd.UI.Table([]string{T("KeyName"), T("Description"),
		T("term"),T("Default Hourly Price Per Instance")})
	items,err := cmd.VirtualServerManager.GetCapacityCreateOptions("RESERVED_CAPACITY")
	if err != nil {
		return slErrors.NewInvalidUsageError("Internal error.")
	}
	for _, item:=range items{
		tableItems.Add(utils.FormatStringPointer(item.KeyName),
			utils.FormatStringPointer(item.Description),
			utils.FormatSLFloatPointerToFloat(item.Capacity),getPrices(item.Prices))
	}

	tableItems.Print()
	tableRegion.Print()
	return nil
}

//Finds the price with the default locationGroupId
func getPrices(prices []datatypes.Product_Item_Price) string {
	itemPrices:= ""
	for _, price := range prices{
		if price.LocationGroupId == nil{
			itemPrices = utils.FormatSLFloatPointerToFloat(price.HourlyRecurringFee)
		}
	}
	return itemPrices
}

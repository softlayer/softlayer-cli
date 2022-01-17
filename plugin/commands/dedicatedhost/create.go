package dedicatedhost

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateCommand struct {
	UI                   terminal.UI
	DedicatedHostManager managers.DedicatedHostManager
	NetworkManager       managers.NetworkManager
	Context              plugin.PluginContext
}

func NewCreateCommand(ui terminal.UI, dedicatedHostManager managers.DedicatedHostManager, networkManager managers.NetworkManager, context plugin.PluginContext) (cmd *CreateCommand) {
	return &CreateCommand{
		UI:                   ui,
		DedicatedHostManager: dedicatedHostManager,
		NetworkManager:       networkManager,
		Context:              context,
	}
}

func (cmd *CreateCommand) Run(c *cli.Context) error {
	size := managers.HOST_DEFAULT_SIZE
	if c.IsSet("size") {
		size = c.String("size")
	}
	hostname := c.String("H")
	if hostname == "" {
		return errors.NewMissingInputError("-H|--hostname")
	}
	domain := c.String("D")
	if domain == "" {
		return errors.NewMissingInputError("-D|--domain")
	}
	datacenter := c.String("d")
	if datacenter == "" {
		return errors.NewMissingInputError("-d|--datacenter")
	}
	billing := "hourly"
	if c.IsSet("b") {
		billing = c.String("b")
		if billing != "hourly" && billing != "monthly" {
			return errors.NewInvalidUsageError(T("[-b|--billing] has to be either hourly or monthly."))
		}
	}
	vlanId := c.Int("v")
	if vlanId == 0 {
		return errors.NewMissingInputError("-v|--vlan-private")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	vlan, err := cmd.NetworkManager.GetVlan(vlanId, "mask[id,primaryRouter[id]]")
	if err != nil {
		return cli.NewExitError(T("Failed to get vlan {{.VlanId}}.\n", map[string]interface{}{"VlanId": vlanId})+err.Error(), 2)
	}
	if !c.IsSet("f") && !c.IsSet("force") && outputFormat != "JSON" && !c.IsSet("test") {
		confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	if vlan.PrimaryRouter == nil || vlan.PrimaryRouter.Id == nil {
		return cli.NewExitError(T("Failed to get vlan primary router ID."), 2)
	}

	orderTemplate, err := cmd.DedicatedHostManager.GenerateOrderTemplate(size, hostname, domain, datacenter, billing, *vlan.PrimaryRouter.Id)
	if err != nil {
		return cli.NewExitError(T("Failed to generate the order template.\n")+err.Error(), 2)
	}

	var orderReceipt = datatypes.Container_Product_Order_Receipt{}
	if c.IsSet("test") {
		_, err := cmd.DedicatedHostManager.VerifyInstanceCreation(orderTemplate)
		if err != nil {
			return cli.NewExitError(T("Failed to verify virtual server creation.\n")+err.Error(), 2)
		}
		cmd.UI.Ok()
		cmd.UI.Print(T("The order is correct."))
	} else {
		orderReceipt, err := cmd.DedicatedHostManager.OrderInstance(orderTemplate)
		if err != nil {
			return cli.NewExitError(T("Failed to Order the dedicatedhost.\n")+err.Error(), 2)
		}
		cmd.UI.Ok()
		cmd.UI.Print(T("The order {{.OrderID}} was placed.", map[string]interface{}{"OrderID": *orderReceipt.OrderId}))
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
	}

	return nil
}

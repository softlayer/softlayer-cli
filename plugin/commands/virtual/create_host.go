package virtual

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type CreateHostCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
	NetworkManager       managers.NetworkManager
	Context              plugin.PluginContext
}

func NewCreateHostCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager, networkManager managers.NetworkManager, context plugin.PluginContext) (cmd *CreateHostCommand) {
	return &CreateHostCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
		NetworkManager:       networkManager,
		Context:              context,
	}
}

func (cmd *CreateHostCommand) Run(c *cli.Context) error {
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
	if !c.IsSet("f") && !c.IsSet("force") && outputFormat != "JSON" {
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
	orderReceipt, err := cmd.VirtualServerManager.CreateDedicatedHost(size, hostname, domain, datacenter, billing, *vlan.PrimaryRouter.Id)
	if err != nil {
		return cli.NewExitError(T("Failed to create dedicated host.\n")+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("The order {{.OrderID}} was placed.", map[string]interface{}{"OrderID": *orderReceipt.OrderId}))
	cmd.UI.Print(T("You may run '{{.CommandName}} sl vs host-list --order {{.OrderID}}' to find this dedicated host after it is ready.",
		map[string]interface{}{"OrderID": *orderReceipt.OrderId, "CommandName": cmd.Context.CLIName()}))
	return nil
}

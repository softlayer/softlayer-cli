package vlan

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	bmxErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type CreateCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
	Context        plugin.PluginContext
}

func NewCreateCommand(ui terminal.UI, networkManager managers.NetworkManager, context plugin.PluginContext) (cmd *CreateCommand) {
	return &CreateCommand{
		UI:             ui,
		NetworkManager: networkManager,
		Context:        context,
	}
}

func (cmd *CreateCommand) Run(c *cli.Context) error {
	if c.IsSet("r") {
		//set routers, then no need to set vlan-type or datacenter
		if c.IsSet("d") || c.IsSet("t") {
			return bmxErr.NewInvalidUsageError(T("[-r|--router] is not allowed with [-d|--datacenter] or [-t|--vlan-type].\nRun '{{.CommandName}} sl vlan options' to check available options.",
				map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
		}
	} else {
		//not set router, then need to set both vlan-type and datacenter
		if !c.IsSet("d") || !c.IsSet("t") {
			return bmxErr.NewInvalidUsageError(T("[-d|--datacenter] and [-t|--vlan-type] are required.\nRun '{{.CommandName}} sl vlan options' to check available options.",
				map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
		}
		vlanType := c.String("t")
		if vlanType != "public" && vlanType != "private" {
			return bmxErr.NewInvalidUsageError(T("[-t|--vlan-type] is required, must be either public or private.\nRun '{{.CommandName}} sl vlan options' to check available options.",
				map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
		}

	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
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

	orderReceipt, err := cmd.NetworkManager.AddVlan(c.String("t"), c.String("d"), c.String("r"), c.String("n"))
	if err != nil {
		return cli.NewExitError(T("Failed to add VLAN.\n")+err.Error(), 2)
	}
	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("The order {{.OrderID}} was placed.", map[string]interface{}{"OrderID": *orderReceipt.OrderId}))
	return nil
}

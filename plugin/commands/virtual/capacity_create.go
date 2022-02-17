package virtual

import (
"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
"github.com/softlayer/softlayer-go/datatypes"
"github.com/urfave/cli"
. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CapacityCreateCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
	ImageManager         managers.ImageManager
	Context              plugin.PluginContext
}

func NewCapacityCreateCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager, context plugin.PluginContext) (cmd *CapacityCreateCommand) {
	return &CapacityCreateCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
		Context:              context,
	}
}

func (cmd *CapacityCreateCommand) Run(c *cli.Context) error {
	capacity_create := datatypes.Container_Product_Order_Virtual_ReservedCapacity{}
	if c.NumFlags() == 0 {
		confirm, err := cmd.UI.Confirm(T("Please make sure you know all the creation options by running command: '{{.CommandName}} sl vs options'. Continue?",
			map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
		params := make(map[string]interface{})
		params["backendRouterId"], _ = cmd.UI.Ask(T("backendRouterId: "))
		params["flavor"], _ = cmd.UI.Ask(T("flavor: "))
		params["quantity"], _ = cmd.UI.Ask(T("instances: "))
		//params["complexType"] = T("SoftLayer_Container_Product_Order_Virtual_ReservedCapacity")
		params["hourly"] = true

		orderReceipt, err := cmd.VirtualServerManager.GenerateInstanceCapacityCreationTemplate(&capacity_create, params)
		if err != nil {
			return err
		}
		createTable(cmd, orderReceipt, c.IsSet("test"))
	} else {
		//create virtual reservedCapacity server with customized parameters
		params, err := verifyCapacityParams(c)
		if err != nil {
			return err
		}
		orderReceipt, err := cmd.VirtualServerManager.GenerateInstanceCapacityCreationTemplate(&capacity_create, params)
		if err != nil {
			return err
		}
		createTable(cmd, orderReceipt, c.IsSet("test"))
	}

	return nil
}

func createTable(cmd *CapacityCreateCommand, receipt interface{}, set bool) {
	if set {
		val := receipt.(datatypes.Container_Product_Order)
		table := cmd.UI.Table([]string{T("name"), T("value")})
		table.Add("name", utils.FormatStringPointer(val.QuoteName))
		table.Add("location", utils.FormatStringPointer(val.LocationObject.LongName))
		if val.Prices != nil {
			for _, price := range val.Prices {
				table.Add("Contract", utils.FormatStringPointer(price.Item.Description))
			}
		}
		table.Add("Hourly Total", utils.FormatSLFloatPointerToInt(val.PostTaxRecurringHourly))
		table.Print()
	} else {
		val := receipt.(datatypes.Container_Product_Order_Receipt)
		table := cmd.UI.Table([]string{T("name"), T("value")})
		table.Add("Order Date", utils.FormatSLTimePointer(val.OrderDate))
		table.Add("Order Id", utils.FormatIntPointer(val.OrderId))
		table.Add("Status", *val.PlacedOrder.Status)
		table.Add("Hourly Total", utils.FormatSLFloatPointerToInt(val.OrderDetails.PostTaxRecurringHourly))
		table.Print()
	}
}

func verifyCapacityParams(c *cli.Context) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	if c.IsSet("flavor") || (c.IsSet("fl")) {
		params["flavor"] = c.String("flavor")
	}
	if c.IsSet("b") {
		params["backendRouterId"] = c.Int("b")
	}
	if c.IsSet("i") {
		params["quantity"] = c.Int("i")
	}
	if c.IsSet("n") {
		params["name"] = c.String("n")
	}
	if c.IsSet("test") {
		params["test"] = c.Bool("test")
	}

	return params, nil
}

func VSCapacityCreateMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "capacity-create",
		Description: T("Create a Reserved Capacity instance."),
		Usage: T(`${COMMAND_NAME} sl vs capacity-create [OPTIONS]
EXAMPLE:
${COMMAND_NAME} sl vs capacity-create -n myvsi -b 1234567 -fl C1_2X2_1_YEAR_TERM -i 2
This command orders a Reserved Capacity instance with name is myvsi, backendRouterId 1234567, flavor C1_2X2_1_YEAR_TERM and 2 instances,
${COMMAND_NAME} sl vs capacity-create --name myvsi --backendRouterId 1234567 --flavor C1_2X2_1_YEAR_TERM --instances 2 --test
This command tests whether the order is valid with above options before the order is actually placed.

WARNING: Reserved Capacity is on a yearly contract and not cancelable until the contract is expired.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("Name for your new reserved capacity  [required]"),
			},
			cli.IntFlag{
				Name:  "b,backendRouterId",
				Usage: T("BackendRouterId, create-options has a list of valid ids to use. [required]"),
			},
			cli.IntFlag{
				Name:  "i,instances",
				Usage: T("Number of VSI instances this capacity reservation can support. [required]"),
			},
			cli.StringFlag{
				Name:  "fl,flavor",
				Usage: T(" Capacity keyname (C1_2X2_1_YEAR_TERM for example). [required]"),
			},
			cli.BoolFlag{
				Name:  "test",
				Usage: T(" Do not actually create the reserved capacity"),
			},
			metadata.ForceFlag(),
		},
	}
}
package subnet

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewCreateCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *CreateCommand) {
	return &CreateCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *CreateCommand) Run(c *cli.Context) error {
	if c.NArg() != 3 {
		return errors.NewInvalidUsageError(T("This command requires three arguments."))
	}
	network := c.Args()[0]
	if network != "public" && network != "private" {
		return errors.NewInvalidUsageError(T("NETWORK has to be either public or private."))
	}
	quantity, err := strconv.Atoi(c.Args()[1])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("QUANTITY")
	}
	vlanID, err := utils.ResolveVlanId(c.Args()[2])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("VLAN ID")
	}
	version := 4
	if c.IsSet("v6") {
		version = 6
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	testOrder := false
	if c.IsSet("test") {
		testOrder = true
	}
	if testOrder == false {
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
	}

	orderReceipt, err := cmd.NetworkManager.AddSubnet(network, quantity, vlanID, version, testOrder)
	if err != nil {
		return cli.NewExitError(T("Failed to add subnet.\n")+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
	}

	if testOrder {
		cmd.UI.Print(T("The order is correct."))
		return nil
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Order {{.OrderID}} was placed.", map[string]interface{}{"OrderID": *orderReceipt.OrderId}))
	cmd.UI.Print("")
	table := cmd.UI.Table([]string{T("item"), T("cost")})
	total := 0.0
	if orderReceipt.OrderDetails != nil && orderReceipt.OrderDetails.Prices != nil && len(orderReceipt.OrderDetails.Prices) > 0 {
		for _, price := range orderReceipt.OrderDetails.Prices {
			rate := 0.0
			if price.RecurringFee != nil {
				rate = float64(*price.RecurringFee)
			}
			if price.Item != nil && price.Item.Description != nil {
				table.Add(*price.Item.Description, strconv.FormatFloat(rate, 'f', 2, 64))
			}
			total += rate
		}
		table.Add(T("Total monthly cost"), strconv.FormatFloat(total, 'f', 2, 64))
	}
	table.Print()
	return nil
}

func SubnetCreateMetaData() cli.Command {
	return cli.Command{
		Category:    "subnet",
		Name:        "create",
		Description: T("Add a new subnet to your account"),
		Usage: T(`${COMMAND_NAME} sl subnet create NETWORK QUANTITY VLAN_ID [OPTIONS]
	
	Add a new subnet to your account. Valid quantities vary by type.
	
	Type    - Valid Quantities (IPv4)
  	public  - 4, 8, 16, 32
  	private - 4, 8, 16, 32, 64

  	Type    - Valid Quantities (IPv6)
	public  - 64

EXAMPLE:
   ${COMMAND_NAME} sl subnet create public 16 567 
   This command creates a public subnet with 16 IPv4 addresses and places it on vlan with ID 567.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "v6,ipv6",
				Usage: T("Order IPv6 Addresses"),
			},
			cli.BoolFlag{
				Name:  "test",
				Usage: T("Do not order the subnet; just get a quote"),
			},
			metadata.ForceFlag(),
			metadata.OutputFlag(),
		},
	}
}

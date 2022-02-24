package globalip

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
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
	version := 4
	if c.IsSet("v6") {
		version = 6
	}
	testOrder := false
	if c.IsSet("test") {
		testOrder = true
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
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

	orderReceipt, err := cmd.NetworkManager.AddGlobalIP(version, testOrder)
	if err != nil {
		return cli.NewExitError(T("Failed to add global IP.\n")+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
	}

	if testOrder {
		cmd.UI.Ok()
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
			var description string
			if price.Item != nil && price.Item.Description != nil {
				description = *price.Item.Description
			}
			table.Add(description, strconv.FormatFloat(rate, 'f', 2, 64))
			total += rate
		}
		table.Add(T("Total monthly cost"), strconv.FormatFloat(total, 'f', 2, 64))
	}
	table.Print()
	return nil
}

func GlobalIpCreateMetaData() cli.Command {
	return cli.Command{
		Category:    "globalip",
		Name:        "create",
		Description: T("Create a global IP"),
		Usage: T(`${COMMAND_NAME} sl globalip create [OPTIONS]

EXAMPLE:
    ${COMMAND_NAME} sl globalip create --v6 
	This command creates an IPv6 address.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "v6",
				Usage: T("Order an IPv6 IP address"),
			},
			cli.BoolFlag{
				Name:  "test",
				Usage: T("Test order"),
			},
			metadata.ForceFlag(),
			metadata.OutputFlag(),
		},
	}
}

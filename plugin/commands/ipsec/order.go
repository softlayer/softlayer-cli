package ipsec

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type OrderCommand struct {
	UI           terminal.UI
	IPSECManager managers.IPSECManager
	Context      plugin.PluginContext
}

func NewOrderCommand(ui terminal.UI, ipsecManager managers.IPSECManager, context plugin.PluginContext) (cmd *OrderCommand) {
	return &OrderCommand{
		UI:           ui,
		IPSECManager: ipsecManager,
		Context:      context,
	}
}

func (cmd *OrderCommand) Run(c *cli.Context) error {
	location := c.String("d")
	if location == "" {
		return errors.NewMissingInputError("-d|--datacenter")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	orderReceipt, err := cmd.IPSECManager.OrderTunnelContext(location)
	if err != nil {
		return cli.NewExitError(T("Failed to order IPSec.Please try again later.\n")+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Order {{.OrderID}} was placed.", map[string]interface{}{"OrderID": *orderReceipt.OrderId}))
	cmd.UI.Print(T("You may run '{{.CommandName}} sl ipsec list --order {{.OrderID}}' to find this IPSec VPN after it is ready.",
		map[string]interface{}{"OrderID": *orderReceipt.OrderId, "CommandName": cmd.Context.CLIName()}))
	return nil
}

func IpsecOrderMetaData() cli.Command {
	return cli.Command{
		Category:    "ipsec",
		Name:        "order",
		Description: T("Order a IPSec VPN tunnel"),
		Usage:       T(`${COMMAND_NAME} sl ipsec order [OPTIONS]`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Short name of the datacenter for the IPSec. For example, dal09[required]"),
			},
			metadata.OutputFlag(),
		},
	}
}

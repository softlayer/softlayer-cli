package ipsec

import (
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type OrderCommand struct {
	*metadata.SoftlayerCommand
	IPSECManager managers.IPSECManager
	Command      *cobra.Command
	Datacenter   string
}

func NewOrderCommand(sl *metadata.SoftlayerCommand) (cmd *OrderCommand) {
	thisCmd := &OrderCommand{
		SoftlayerCommand: sl,
		IPSECManager:     managers.NewIPSECManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "order",
		Short: T("Order a IPSec VPN tunnel"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVarP(&thisCmd.Datacenter, "datacenter", "d", "", T("Short name of the datacenter for the IPSec. For example, dal09[required]"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *OrderCommand) Run(args []string) error {
	location := cmd.Datacenter
	if location == "" {
		return errors.NewMissingInputError("-d|--datacenter")
	}

	outputFormat := cmd.GetOutputFlag()

	orderReceipt, err := cmd.IPSECManager.OrderTunnelContext(location)
	if err != nil {
		return errors.NewAPIError(T("Failed to order IPSec.Please try again later.\n"), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Order {{.OrderID}} was placed.", map[string]interface{}{"OrderID": *orderReceipt.OrderId}))
	cmd.UI.Print(T("You may run '{{.CommandName}} sl ipsec list --order {{.OrderID}}' to find this IPSec VPN after it is ready.",
		map[string]interface{}{"OrderID": *orderReceipt.OrderId, "CommandName": "ibmcloud"}))
	return nil
}

package subnet

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
	"strconv"
	"strings"
)

type CreateCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	Ipv6           bool
	Test           bool
	Force          bool
}

func NewCreateCommand(sl *metadata.SoftlayerCommand) *CreateCommand {
	thisCmd := &CreateCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "create" + strings.ToUpper(fmt.Sprintf(" %s %s %s", T("Network"), T("Quantity"), T("VLAN"))),
		Short: T("Add a new subnet to your account"),
		Long: T(`Valid quantities vary by type.
	- public IPv4: 4, 8, 16, 32
	- private IPv4: 4, 8, 16, 32, 64
	- public IPv6: 64

EXAMPLE:
	${COMMAND_NAME} sl subnet create public 16 567
	This command creates a public subnet with 16 IPv4 addresses and places it on vlan with ID 567.`),
		Args: metadata.ThreeArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVarP(&thisCmd.Ipv6, "ipv6", "6", false, T("Order IPv6 Addresses"))
	cobraCmd.Flags().BoolVar(&thisCmd.Test, "test", false, T("Do not order the subnet; just get a quote"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CreateCommand) Run(args []string) error {
	network := args[0]
	if network != "public" && network != "private" {
		return errors.NewInvalidUsageError(T("NETWORK has to be either public or private."))
	}
	quantity, err := strconv.Atoi(args[1])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("QUANTITY")
	}
	vlanID, err := utils.ResolveVlanId(args[2])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("VLAN ID")
	}
	version := 4
	if cmd.Ipv6 {
		version = 6
	}

	outputFormat := cmd.GetOutputFlag()

	testOrder := false
	if cmd.Test {
		testOrder = true
	}
	if testOrder == false {
		if !cmd.Force {
			confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
			if err != nil {
				return err
			}
			if !confirm {
				cmd.UI.Print(T("Aborted."))
				return nil
			}
		}
	}

	orderReceipt, err := cmd.NetworkManager.AddSubnet(network, quantity, vlanID, version, testOrder)
	if err != nil {
		return errors.NewAPIError(T("Failed to add subnet.\n"), err.Error(), 2)
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

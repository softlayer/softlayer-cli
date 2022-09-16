package globalip

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	V6             bool
	Test           bool
	Force          bool
}

func NewCreateCommand(sl *metadata.SoftlayerCommand) *CreateCommand {
	thisCmd := &CreateCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "create",
		Short: T("Create a global IP"),
		Long: T(`${COMMAND_NAME} sl globalip create [OPTIONS]

EXAMPLE:
    ${COMMAND_NAME} sl globalip create --v6 
	This command creates an IPv6 address.`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVar(&thisCmd.V6, "v6", false, T("Order an IPv6 IP address"))
	cobraCmd.Flags().BoolVar(&thisCmd.Test, "test", false, T("Test order"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CreateCommand) Run(args []string) error {
	version := 4
	if cmd.V6 {
		version = 6
	}
	testOrder := false
	if cmd.Test {
		testOrder = true
	}

	outputFormat := cmd.GetOutputFlag()

	if !testOrder {
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

	orderReceipt, err := cmd.NetworkManager.AddGlobalIP(version, testOrder)
	if err != nil {
		return errors.NewAPIError(T("Failed to add global IP.\n"), err.Error(), 2)
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

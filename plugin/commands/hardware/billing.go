package hardware

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type BillingCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
}

func NewBillingCommand(sl *metadata.SoftlayerCommand) (cmd *BillingCommand) {
	thisCmd := &BillingCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "billing " + T("IDENTIFIER"),
		Short: T("Get billing for a hardware device."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *BillingCommand) Run(args []string) error {
	hardwareID, err := strconv.Atoi(args[0])

	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat := cmd.GetOutputFlag()

	mask := "mask[id,billingItem[id,recurringFee,nextInvoiceTotalRecurringAmount,provisionTransaction[createDate],nextInvoiceChildren[description,categoryCode,nextInvoiceTotalRecurringAmount]]]"
	hardware, err := cmd.HardwareManager.GetHardware(hardwareID, mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get hardware server {{.ID}}.\n", map[string]interface{}{"ID": hardwareID}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, hardware.BillingItem)
	}

	table := cmd.UI.Table([]string{T("name"), T("value")})
	table.Add("Id", utils.FormatIntPointer(&hardwareID))
	table.Add("Billing Item Id", utils.FormatIntPointer(hardware.BillingItem.Id))
	table.Add("Recurring Fee", fmt.Sprintf("%.2f", *hardware.BillingItem.RecurringFee))
	table.Add("Total", fmt.Sprintf("%.2f", *hardware.BillingItem.NextInvoiceTotalRecurringAmount))
	table.Add("Provision Date", utils.FormatSLTimePointer(hardware.BillingItem.ProvisionTransaction.CreateDate))

	if hardware.BillingItem != nil && hardware.BillingItem.NextInvoiceTotalRecurringAmount != nil {
		buf := new(bytes.Buffer)
		priceTable := terminal.NewTable(buf, []string{T("Item"), T("CategoryCode"), T("Recurring Price")})
		for _, item := range hardware.BillingItem.NextInvoiceChildren {
			if item.NextInvoiceTotalRecurringAmount != nil {
				priceTable.Add(*item.Description, *item.CategoryCode, fmt.Sprintf("%.2f", *item.NextInvoiceTotalRecurringAmount))
			}
		}
		priceTable.Print()
		table.Add("Prices", buf.String())
	}
	table.Print()
	return nil
}

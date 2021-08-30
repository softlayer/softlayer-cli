package hardware

import (
	"bytes"
	"fmt"
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

type BillingCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewBillingCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *BillingCommand) {
	return &BillingCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *BillingCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	hardwareID, err := strconv.Atoi(c.Args()[0])

	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	hardware, err := cmd.HardwareManager.GetHardware(hardwareID, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get hardware server {{.ID}}.\n", map[string]interface{}{"ID": hardwareID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, hardware.BillingItem)
	}

	table := cmd.UI.Table([]string{T("name"), T("value")})
	table.Add("Id", utils.FormatIntPointer(&hardwareID))
	table.Add("Billing Item Id", utils.FormatIntPointer(hardware.BillingItem.Id))
	table.Add("Recurring Fee", utils.FormatSLFloatPointerToFloat(hardware.BillingItem.RecurringFee))
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

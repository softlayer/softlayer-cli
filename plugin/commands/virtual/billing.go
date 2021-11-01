package virtual

import (
	"bytes"
	"fmt"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type BillingCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewBillingCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *BillingCommand) {
	return &BillingCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *BillingCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError("This command requires one argument.")
	}
	vsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	virtualGuest, err := cmd.VirtualServerManager.GetInstance(vsID, managers.INSTANCE_DETAIL_MASK)
	if err != nil {
		return cli.NewExitError(T("Failed to get virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": vsID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, virtualGuest)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(virtualGuest.Id))
	table.Add(T("Billing Item Id"), utils.FormatIntPointer(virtualGuest.BillingItem.Id))
	table.Add(T("Recurring Fee"), utils.FormatSLFloatPointerToFloat(virtualGuest.BillingItem.RecurringFee))
	table.Add(T("Total"), utils.FormatSLFloatPointerToFloat(virtualGuest.BillingItem.NextInvoiceTotalRecurringAmount))
	table.Add(T("Provisioning Date"), utils.FormatSLTimePointer(virtualGuest.ProvisionDate))

	buf := new(bytes.Buffer)
	tablePrices := terminal.NewTable(buf, []string{T("Recurring Fee")})
	for _, item := range virtualGuest.BillingItem.Children {
		tablePrices.Add(fmt.Sprintf("%.2f", *item.NextInvoiceTotalRecurringAmount))
	}
	tablePrices.Print()
	table.Add(T("Prices"), buf.String())
	table.Print()
	return nil
}

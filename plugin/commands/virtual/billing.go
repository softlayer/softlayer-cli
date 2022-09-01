package virtual

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type BillingCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
}

func NewBillingCommand(sl *metadata.SoftlayerCommand) (cmd *BillingCommand) {
	thisCmd := &BillingCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "billing " + T("IDENTIFIER"),
		Short: T("Get billing details for a virtual server instance"),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *BillingCommand) Run(args []string) error {

	vsID, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	outputFormat := cmd.GetOutputFlag()

	virtualGuest, err := cmd.VirtualServerManager.GetInstance(vsID, managers.INSTANCE_DETAIL_MASK)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to get virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": vsID}), err.Error(), 2)
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
	tablePrices := terminal.NewTable(buf, []string{T("Item"), T("Recurring Fee")})
	for _, item := range virtualGuest.BillingItem.Children {
		tablePrices.Add(utils.FormatStringPointer(item.Description), fmt.Sprintf("%.2f", *item.NextInvoiceTotalRecurringAmount))
	}
	tablePrices.Print()
	table.Add(T("Prices"), buf.String())
	table.Print()
	return nil
}

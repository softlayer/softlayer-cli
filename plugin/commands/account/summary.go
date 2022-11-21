package account

import (
	"bytes"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SummaryCommand struct {
	*metadata.SoftlayerCommand
	AccountManager managers.AccountManager
	Command        *cobra.Command
}

func NewSummaryCommand(sl *metadata.SoftlayerCommand) *SummaryCommand {
	thisCmd := &SummaryCommand{
		SoftlayerCommand: sl,
		AccountManager:   managers.NewAccountManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "summary",
		Short: T("Prints some various bits of information about an account."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SummaryCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()
	mask := "mask[nextInvoiceTotalAmount,pendingInvoice[invoiceTotalAmount],blockDeviceTemplateGroupCount,dedicatedHostCount,domainCount,hardwareCount,networkStorageCount,openTicketCount,networkVlanCount,subnetCount,userCount,virtualGuestCount]"
	account, err := cmd.AccountManager.GetSummary(mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get summary."), err.Error(), 2)
	}
	PrintSummary(account, cmd.UI, outputFormat)

	return nil
}

func PrintSummary(account datatypes.Account, ui terminal.UI, outputFormat string) {
	bufEvent := new(bytes.Buffer)
	table := terminal.NewTable(bufEvent, []string{
		T("Name"),
		T("Value"),
	})

	table.Add("Company Name", utils.FormatStringPointer(account.CompanyName))
	// table.Add("Balance", utils.FormatSLFloatPointerToFloat(account.PendingInvoice.StartingBalance))
	// table.Add("Upcoming Invoice", utils.FormatSLFloatPointerToFloat(account.PendingInvoice.InvoiceTotalAmount))
	Balance := "-"
	UpcomingInvoice := "-"
	if account.PendingInvoice != nil {
		Balance = utils.FormatSLFloatPointerToFloat(account.PendingInvoice.StartingBalance)
		UpcomingInvoice = utils.FormatSLFloatPointerToFloat(account.PendingInvoice.InvoiceTotalAmount)
	}
	table.Add("Balance", Balance)
	table.Add("Upcoming Invoice", UpcomingInvoice)
	table.Add("Image Templates", utils.FormatUIntPointer(account.BlockDeviceTemplateGroupCount))
	table.Add("Dedicated Hosts", utils.FormatUIntPointer(account.DedicatedHostCount))
	table.Add("Hardware", utils.FormatUIntPointer(account.HardwareCount))
	table.Add("Virtual Guests", utils.FormatUIntPointer(account.VirtualGuestCount))
	table.Add("Domains", utils.FormatUIntPointer(account.DomainCount))
	table.Add("Network Storage Volumes", utils.FormatUIntPointer(account.NetworkStorageCount))
	table.Add("Open Tickets", utils.FormatUIntPointer(account.OpenTicketCount))
	table.Add("Network Vlans", utils.FormatUIntPointer(account.NetworkVlanCount))
	table.Add("Subnets", utils.FormatUIntPointer(account.SubnetCount))
	table.Add("Users", utils.FormatUIntPointer(account.UserCount))

	utils.PrintTableWithTitle(ui, table, bufEvent, "Account Snapshot", outputFormat)
}

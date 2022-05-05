package account

import (
	"bytes"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SummaryCommand struct {
	UI             terminal.UI
	AccountManager managers.AccountManager
}

func NewSummaryCommand(ui terminal.UI, accountManager managers.AccountManager) (cmd *SummaryCommand) {
	return &SummaryCommand{
		UI:             ui,
		AccountManager: accountManager,
	}
}

func SummaryMetaData() cli.Command {
	return cli.Command{
		Category:    "account",
		Name:        "summary",
		Description: T("Prints some various bits of information about an account."),
		Usage:       T(`${COMMAND_NAME} sl account summary [OPTIONS]`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

func (cmd *SummaryCommand) Run(c *cli.Context) error {
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	mask := "mask[nextInvoiceTotalAmount,pendingInvoice[invoiceTotalAmount],blockDeviceTemplateGroupCount,dedicatedHostCount,domainCount,hardwareCount,networkStorageCount,openTicketCount,networkVlanCount,subnetCount,userCount,virtualGuestCount]"
	account, err := cmd.AccountManager.GetSummary(mask)
	if err != nil {
		return cli.NewExitError(T("Failed to get summary.")+err.Error(), 2)
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
	table.Add("Balance", utils.FormatSLFloatPointerToFloat(account.PendingInvoice.StartingBalance))
	table.Add("Upcoming Invoice", utils.FormatSLFloatPointerToFloat(account.PendingInvoice.InvoiceTotalAmount))
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

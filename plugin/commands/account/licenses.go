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

type LicensesCommand struct {
	*metadata.SoftlayerCommand
	AccountManager managers.AccountManager
	Command        *cobra.Command
}

func NewLicensesCommand(sl *metadata.SoftlayerCommand) *LicensesCommand {
	thisCmd := &LicensesCommand{
		SoftlayerCommand: sl,
		AccountManager:   managers.NewAccountManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "licenses",
		Short: T("Show all licenses."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *LicensesCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	mask := "mask[billingItem[categoryCode,createDate,description],key,id,ipAddress,softwareDescription[longDescription,name,manufacturer],subnet]"
	virtualLicenses, err := cmd.AccountManager.GetActiveVirtualLicenses(mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get virtual licenses."), err.Error(), 2)
	}
	PrintVirtualLicenses(virtualLicenses, cmd.UI, outputFormat)

	mask = "mask[billingItem,softwareDescription]"
	vmwares, err := cmd.AccountManager.GetActiveAccountLicenses(mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get account licenses."), err.Error(), 2)
	}
	PrintVmwaresLicenses(vmwares, cmd.UI, outputFormat)
	return nil
}

func PrintVirtualLicenses(virtualLicenses []datatypes.Software_VirtualLicense, ui terminal.UI, outputFormat string) {
	bufEvent := new(bytes.Buffer)
	table := terminal.NewTable(bufEvent, []string{
		T("Id"),
		T("Ip_address"),
		T("Manufacturer"),
		T("Software"),
		T("Key"),
		T("Subnet"),
		T("Subnet notes"),
	})

	for _, virtualLicense := range virtualLicenses {
		SoftwareDescriptionManufacturer := "-"
		SoftwareDescriptionLongDescription := "-"
		if virtualLicense.SoftwareDescription != nil {
			SoftwareDescriptionManufacturer = utils.FormatStringPointer(virtualLicense.SoftwareDescription.Manufacturer)
			SoftwareDescriptionLongDescription = utils.ShortenString(utils.FormatStringPointer(virtualLicense.SoftwareDescription.LongDescription))
		}

		SubnetBroadcastAddress := "-"
		SubnetNote := "-"
		if virtualLicense.Subnet != nil {
			SubnetBroadcastAddress = utils.FormatStringPointer(virtualLicense.Subnet.BroadcastAddress)
			SubnetNote = utils.FormatStringPointer(virtualLicense.Subnet.Note)
		}

		table.Add(
			utils.FormatIntPointer(virtualLicense.Id),
			utils.FormatStringPointer(virtualLicense.IpAddress),
			SoftwareDescriptionManufacturer,
			SoftwareDescriptionLongDescription,
			utils.FormatStringPointer(virtualLicense.Key),
			SubnetBroadcastAddress,
			SubnetNote,
		)
	}

	utils.PrintTableWithTitle(ui, table, bufEvent, "Control Panel Licenses", outputFormat)
}

func PrintVmwaresLicenses(vmwares []datatypes.Software_AccountLicense, ui terminal.UI, outputFormat string) {
	bufEvent := new(bytes.Buffer)
	table := terminal.NewTable(bufEvent, []string{
		T("Name"),
		T("License key"),
		T("CPUs"),
		T("Description"),
		T("Manufacturer"),
		T("Required User"),
	})

	for _, vmware := range vmwares {
		SoftwareDescriptionName := "-"
		SoftwareDescriptionManufacturer := "-"
		SoftwareDescriptionRequiredUser := "-"
		if vmware.SoftwareDescription != nil {
			SoftwareDescriptionName = utils.FormatStringPointer(vmware.SoftwareDescription.Name)
			SoftwareDescriptionManufacturer = utils.FormatStringPointer(vmware.SoftwareDescription.Manufacturer)
			SoftwareDescriptionRequiredUser = utils.FormatStringPointer(vmware.SoftwareDescription.RequiredUser)
		}

		BillingItemDescription := "-"
		if vmware.BillingItem != nil {
			BillingItemDescription = utils.FormatStringPointer(vmware.BillingItem.Description)
		}

		table.Add(
			SoftwareDescriptionName,
			utils.FormatStringPointer(vmware.Key),
			utils.FormatStringPointer(vmware.Capacity),
			BillingItemDescription,
			SoftwareDescriptionManufacturer,
			SoftwareDescriptionRequiredUser,
		)
	}

	utils.PrintTableWithTitle(ui, table, bufEvent, "VMware Licenses", outputFormat)
}

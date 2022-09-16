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

type DetailCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
	Passwords       bool
	Price           bool
	Components      bool
}

func NewDetailCommand(sl *metadata.SoftlayerCommand) (cmd *DetailCommand) {
	thisCmd := &DetailCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}


	cobraCmd := &cobra.Command{
		Use:   "detail " + T("IDENTIFIER"),
		Short: T("Get details for a hardware server"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVarP(&thisCmd.Passwords, "passwords", "p", false, T("Show passwords (check over your shoulder!)"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Price, "price", "c", false, T("Show associated prices"))
	cobraCmd.Flags().BoolVar(&thisCmd.Components, "components", false, T("Show associated hardware components"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DetailCommand) Run(args []string) error {
	hardwareId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat := cmd.GetOutputFlag()

	hardware, err := cmd.HardwareManager.GetHardware(hardwareId, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, hardware)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(hardware.Id))
	table.Add(T("GUID"), utils.FormatStringPointer(hardware.GlobalIdentifier))
	table.Add(T("Hostname"), utils.FormatStringPointer(hardware.Hostname))
	table.Add(T("Domain"), utils.FormatStringPointer(hardware.Domain))
	table.Add(T("FQDN"), utils.FormatStringPointer(hardware.FullyQualifiedDomainName))
	if hardware.HardwareStatus != nil {
		table.Add(T("Status"), utils.FormatStringPointer(hardware.HardwareStatus.Status))
	}
	if hardware.Datacenter != nil {
		table.Add(T("Datacenter"), utils.FormatStringPointer(hardware.Datacenter.Name))
	}
	table.Add(T("CPU cores"), utils.FormatUIntPointer(hardware.ProcessorPhysicalCoreAmount))
	table.Add(T("Memory"), utils.FormatUIntPointer(hardware.MemoryCapacity)+"G")
	table.Add(T("Public IP"), utils.FormatStringPointer(hardware.PrimaryIpAddress))
	table.Add(T("Private IP"), utils.FormatStringPointer(hardware.PrimaryBackendIpAddress))
	table.Add(T("IPMI IP"), utils.FormatStringPointer(hardware.NetworkManagementIpAddress))
	if hardware.OperatingSystem != nil &&
		hardware.OperatingSystem.SoftwareLicense != nil &&
		hardware.OperatingSystem.SoftwareLicense.SoftwareDescription != nil {
		table.Add(T("OS"), utils.FormatStringPointer(hardware.OperatingSystem.SoftwareLicense.SoftwareDescription.Name))
		table.Add(T("OS version"), utils.FormatStringPointer(hardware.OperatingSystem.SoftwareLicense.SoftwareDescription.Version))
	}
	table.Add(T("Created"), utils.FormatSLTimePointer(hardware.ProvisionDate))
	if hardware.BillingItem != nil &&
		hardware.BillingItem.OrderItem != nil &&
		hardware.BillingItem.OrderItem.Order != nil &&
		hardware.BillingItem.OrderItem.Order.UserRecord != nil {
		table.Add(T("Owner"), utils.FormatStringPointer(hardware.BillingItem.OrderItem.Order.UserRecord.Username))
	}
	if hardware.Notes != nil && *hardware.Notes != "" {
		table.Add(T("Note"), utils.FormatStringPointer(hardware.Notes))
	}
	if tags := hardware.TagReferences; len(tags) > 0 {
		table.Add(T("Tag"), utils.TagRefsToString(tags))
	}

	if vlans := hardware.NetworkVlans; len(vlans) > 0 {
		buf := new(bytes.Buffer)
		vlanTable := terminal.NewTable(buf, []string{T("Type"), T("Number"), T("ID")})
		for _, vlan := range vlans {
			vlanTable.Add(utils.FormatStringPointer(vlan.NetworkSpace),
				utils.FormatIntPointer(vlan.VlanNumber),
				utils.FormatIntPointer(vlan.Id))
		}
		vlanTable.Print()
		table.Add("Vlans", buf.String())
	}

	if cmd.Price {
		if hardware.BillingItem != nil && hardware.BillingItem.NextInvoiceTotalRecurringAmount != nil {
			buf := new(bytes.Buffer)
			priceTable := terminal.NewTable(buf, []string{T("Item"), T("CategoryCode"), T("Recurring Price")})

			totalPrice := hardware.BillingItem.NextInvoiceTotalRecurringAmount
			priceTable.Add("Total", "-", fmt.Sprintf("%.2f", *totalPrice))
			sum := *hardware.BillingItem.NextInvoiceTotalRecurringAmount
			for _, item := range hardware.BillingItem.NextInvoiceChildren {
				if item.NextInvoiceTotalRecurringAmount != nil {
					sum += *item.NextInvoiceTotalRecurringAmount
					priceTable.Add(*item.Description, *item.CategoryCode, fmt.Sprintf("%.2f", *item.NextInvoiceTotalRecurringAmount))
				}
			}
			priceTable.Print()
			table.Add("Prices", buf.String())
			table.Add(T("Price rate"), fmt.Sprintf("%.2f", sum))
		}
	}

	if cmd.Passwords {
		if hardware.OperatingSystem != nil && hardware.OperatingSystem.Passwords != nil {
			buf := new(bytes.Buffer)
			userTable := terminal.NewTable(buf, []string{T("Username"), T("Password")})
			for _, pwd := range hardware.OperatingSystem.Passwords {
				userTable.Add(utils.FormatStringPointer(pwd.Username), utils.FormatStringPointer(pwd.Password))
			}
			userTable.Print()
			table.Add("Users", buf.String())
		}

		if hardware.RemoteManagementAccounts != nil {
			buf := new(bytes.Buffer)
			userTable := terminal.NewTable(buf, []string{T("IPMI_username"), T("Password")})
			for _, pwd := range hardware.RemoteManagementAccounts {
				userTable.Add(utils.FormatStringPointer(pwd.Username), utils.FormatStringPointer(pwd.Password))
			}
			userTable.Print()
			table.Add("Remote users", buf.String())
		}
	}

	if cmd.Components {
		components, err := cmd.HardwareManager.GetHardwareComponents(hardwareId)
		componentIds := []int{}
		if err != nil {
			return errors.NewAPIError(T("Failed to get components\n"), err.Error(), 2)
		}
		buf := new(bytes.Buffer)
		componentTable := terminal.NewTable(buf, []string{T("Name"), T("Firmware version"), T("Firmware build date"), T("Type")})
		for _, component := range components {
			if utils.IntInSlice(*component.Id, componentIds) == -1 {
				componentTable.Add(
					utils.FormatStringPointer(component.HardwareComponentModel.LongDescription),
					utils.FormatStringPointer(component.HardwareComponentModel.Firmwares[0].Version),
					utils.FormatSLTimePointer(component.HardwareComponentModel.Firmwares[0].CreateDate),
					utils.FormatStringPointer(component.HardwareComponentModel.HardwareGenericComponentModel.HardwareComponentType.KeyName),
				)
				componentIds = append(componentIds, *component.Id)
			}
		}
		componentTable.Print()
		table.Add("Components", buf.String())
	}

	table.Print()
	return nil
}

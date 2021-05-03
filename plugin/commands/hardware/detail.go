package hardware

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewDetailCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *DetailCommand) {
	return &DetailCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *DetailCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	hardwareId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	hardware, err := cmd.HardwareManager.GetHardware(hardwareId, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId})+err.Error(), 2)
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

	if c.IsSet("price") {
		if hardware.BillingItem != nil && hardware.BillingItem.NextInvoiceTotalRecurringAmount != nil {
			sum := *hardware.BillingItem.NextInvoiceTotalRecurringAmount
			for _, item := range hardware.BillingItem.Children {
				if item.NextInvoiceTotalRecurringAmount != nil {
					sum += *item.NextInvoiceTotalRecurringAmount
				}
			}
			table.Add(T("Price rate"), fmt.Sprintf("%.2f", sum))
		}
	}

	if c.IsSet("passwords") {
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
	table.Print()
	return nil
}

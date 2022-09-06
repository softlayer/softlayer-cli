package subnet

import (
	"bytes"
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type LookupCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
}

func NewLookupCommand(sl *metadata.SoftlayerCommand) *LookupCommand {
	thisCmd := &LookupCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "lookup " + T("IDENTIFIER"),
		Short: T("Find an IP address and display its subnet and device information."),
		Long: T(`${COMMAND_NAME} sl subnet lookup IP_ADDRESS [OPTIONS]
	
EXAMPLE:
	${COMMAND_NAME} sl subnet lookup 9.125.235.255
	This command finds the IP address record with IP address 9.125.235.255 and displays its subnet and device information.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *LookupCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	ipAddress := args[0]
	ipAddressRecord, err := cmd.NetworkManager.IPLookup(ipAddress)
	if err != nil {
		return errors.NewAPIError(T("Failed to lookup IP address: {{.IPAddress}}.\n", map[string]interface{}{"IPAddress": ipAddress}), err.Error(), 2)
	}

	if ipAddressRecord.Id == nil {
		cmd.UI.Print(T("IP address {{.IPAddress}} is not found.", map[string]interface{}{"IPAddress": ipAddress}))
		return nil
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(ipAddressRecord.Id))
	table.Add(T("ipAddress"), utils.FormatStringPointer(ipAddressRecord.IpAddress))
	if ipAddressRecord.Subnet == nil {
		table.Add(T("subnet"), T("none"))
	} else {
		buf := new(bytes.Buffer)
		subnetTable := terminal.NewTable(buf, []string{T("name"), T("value")})
		subnetTable.Add(T("ID"), utils.FormatIntPointer(ipAddressRecord.Subnet.Id))
		subnetTable.Add(T("identifier"), fmt.Sprintf("%s/%s", utils.FormatStringPointer(ipAddressRecord.Subnet.NetworkIdentifier), utils.FormatIntPointer(ipAddressRecord.Subnet.Cidr)))
		subnetTable.Add(T("netmask"), utils.FormatStringPointer(ipAddressRecord.Subnet.Netmask))
		subnetTable.Add(T("gateway"), utils.FormatStringPointer(ipAddressRecord.Subnet.Gateway))
		subnetTable.Add(T("type"), utils.FormatStringPointer(ipAddressRecord.Subnet.SubnetType))
		subnetTable.Print()
		table.Add(T("subnet"), buf.String())
	}

	if ipAddressRecord.VirtualGuest == nil && ipAddressRecord.Hardware == nil {
		table.Add(T("device"), T("none"))
	} else {
		buf := new(bytes.Buffer)
		deviceTable := terminal.NewTable(buf, []string{T("ID"), T("FQDN"), T("type")})
		var deviceID, deviceType, FQDN string
		if ipAddressRecord.VirtualGuest != nil {
			deviceID = utils.FormatIntPointer(ipAddressRecord.VirtualGuest.Id)
			FQDN = utils.FormatStringPointer(ipAddressRecord.VirtualGuest.FullyQualifiedDomainName)
			deviceType = T("virtual guest")
		} else if ipAddressRecord.Hardware != nil {
			deviceID = utils.FormatIntPointer(ipAddressRecord.Hardware.Id)
			FQDN = utils.FormatStringPointer(ipAddressRecord.Hardware.FullyQualifiedDomainName)
			deviceType = T("hardware")
		}
		deviceTable.Add(deviceID, FQDN, deviceType)
		deviceTable.Print()
		table.Add(T("device"), buf.String())
	}
	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

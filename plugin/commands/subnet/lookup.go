package subnet

import (
	"bytes"
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type LookupCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewLookupCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *LookupCommand) {
	return &LookupCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *LookupCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	ipAddress := c.Args()[0]
	ipAddressRecord, err := cmd.NetworkManager.IPLookup(ipAddress)
	if err != nil {
		return cli.NewExitError(T("Failed to lookup IP address: {{.IPAddress}}.\n", map[string]interface{}{"IPAddress": ipAddress})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, ipAddressRecord)
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
	table.Print()
	return nil
}

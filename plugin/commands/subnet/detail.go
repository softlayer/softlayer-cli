package subnet

import (
	"bytes"
	"fmt"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewDetailCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *DetailCommand) {
	return &DetailCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *DetailCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	subnetID, err := utils.ResolveSubnetId(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Subnet ID")
	}

	subnet, err := cmd.NetworkManager.GetSubnet(subnetID, "ipAddresses[id, ipAddress,note], datacenter, virtualGuests, hardware,networkVlan[networkSpace], tagReferences")
	if err != nil {
		return cli.NewExitError(T("Failed to get subnet: {{.ID}}.\n", map[string]interface{}{"ID": subnetID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, subnet)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(subnet.Id))
	table.Add(T("identifier"), fmt.Sprintf("%s/%s", utils.FormatStringPointer(subnet.NetworkIdentifier), utils.FormatIntPointer(subnet.Cidr)))
	if subnet.SubnetType != nil {
		table.Add(T("subnet type"), utils.FormatStringPointer(subnet.SubnetType))
	}
	if subnet.NetworkVlan != nil {
		table.Add(T("network space"), utils.FormatStringPointer(subnet.NetworkVlan.NetworkSpace))
	}
	table.Add(T("gateway"), utils.FormatStringPointer(subnet.Gateway))
	table.Add(T("broadcast"), utils.FormatStringPointer(subnet.BroadcastAddress))

	if subnet.Datacenter != nil {
		table.Add(T("datacenter"), utils.FormatStringPointer(subnet.Datacenter.Name))
	}
	table.Add(T("usable ips"), utils.FormatSLFloatPointerToInt(subnet.UsableIpAddressCount))

	if !c.IsSet("no-ipAddress") {
		if subnet.IpAddresses == nil || len(subnet.IpAddresses) == 0 {
			table.Add(T("ip address"), T("none"))
		} else {
			buf := new(bytes.Buffer)
			ipTable := terminal.NewTable(buf, []string{T("ip"), T("ipAddress")})
			for _, ip := range subnet.IpAddresses {
				ipTable.Add(utils.FormatIntPointer(ip.Id),
					utils.FormatStringPointer(ip.IpAddress))
			}
			ipTable.Print()
			table.Add(T("ip address"), buf.String())
		}
	}

	if !c.IsSet("no-vs") {
		if subnet.VirtualGuests == nil || len(subnet.VirtualGuests) == 0 {
			table.Add(T("virtual guests"), T("none"))
		} else {
			buf := new(bytes.Buffer)
			vsTable := terminal.NewTable(buf, []string{T("hostname"), T("domain"), T("public_ip"), T("private_ip")})
			for _, vs := range subnet.VirtualGuests {
				vsTable.Add(utils.FormatStringPointer(vs.Hostname),
					utils.FormatStringPointer(vs.Domain),
					utils.FormatStringPointer(vs.PrimaryIpAddress),
					utils.FormatStringPointer(vs.PrimaryBackendIpAddress))
			}
			vsTable.Print()
			table.Add(T("virtual guests"), buf.String())
		}
	}
	if !c.IsSet("no-hardware") {
		if subnet.Hardware == nil || len(subnet.Hardware) == 0 {
			table.Add(T("hardware"), T("none"))
		} else {
			buf := new(bytes.Buffer)
			hwTable := terminal.NewTable(buf, []string{T("hostname"), T("domain"), T("public_ip"), T("private_ip")})
			for _, hw := range subnet.Hardware {
				hwTable.Add(utils.FormatStringPointer(hw.Hostname),
					utils.FormatStringPointer(hw.Domain),
					utils.FormatStringPointer(hw.PrimaryIpAddress),
					utils.FormatStringPointer(hw.PrimaryBackendIpAddress))
			}
			hwTable.Print()
			table.Add(T("hardware"), buf.String())
		}
	}

	if !c.IsSet("no-tags") {
		if subnet.TagReferences == nil || len(subnet.TagReferences) == 0 {
			table.Add(T("tags"), T("none"))
		} else {
			buf := new(bytes.Buffer)
			vsTable := terminal.NewTable(buf, []string{T("id")})
			for _, tag := range subnet.TagReferences {
				vsTable.Add(utils.FormatIntPointer(tag.TagId))
			}
			vsTable.Print()
			table.Add(T("tags"), buf.String())
		}
	}
	table.Print()
	return nil
}

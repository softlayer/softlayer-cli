package vlan

import (
	"bytes"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewDetailCommand(ui terminal.UI, networkManager managers.NetworkManager) *DetailCommand {
	return &DetailCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *DetailCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	id, err := utils.ResolveVlanId(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("VLAN ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	vlan, err := cmd.NetworkManager.GetVlan(id, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get VLAN: {{.VLANID}}.\n", map[string]interface{}{"VLANID": id})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, vlan)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("id"), utils.FormatIntPointer(vlan.Id))
	table.Add(T("number"), utils.FormatIntPointer(vlan.VlanNumber))
	if vlan.PrimaryRouter != nil {
		table.Add(T("datacenter"), utils.FormatStringPointer(vlan.PrimaryRouter.DatacenterName))
		table.Add(T("primary_router"), utils.FormatStringPointer(vlan.PrimaryRouter.FullyQualifiedDomainName))
	}
	firewall := T("Yes")
	if len(vlan.FirewallInterfaces) <= 0 {
		firewall = T("No")
	}
	table.Add(T("firewall"), firewall)

	subnets := vlan.Subnets
	if len(subnets) == 0 {
		table.Add(T("subnets"), T("none"))
	} else {
		buf := new(bytes.Buffer)
		snTable := terminal.NewTable(buf, []string{T("ID"), T("identifier"), T("netmask"), T("gateway"), T("type"), T("usable_ips")})
		for _, subnet := range subnets {
			snTable.Add(utils.FormatIntPointer(subnet.Id),
				utils.FormatStringPointer(subnet.NetworkIdentifier),
				utils.FormatStringPointer(subnet.Netmask),
				utils.FormatStringPointer(subnet.Gateway),
				utils.FormatStringPointer(subnet.SubnetType),
				utils.FormatSLFloatPointerToInt(subnet.UsableIpAddressCount))
		}
		snTable.Print()
		table.Add(T("subnets"), buf.String())
	}

	if !c.IsSet("no-vs") {
		vs := vlan.VirtualGuests
		if len(vs) == 0 {
			table.Add(T("virtual servers"), T("none"))
		} else {
			buf := new(bytes.Buffer)
			vsTable := terminal.NewTable(buf, []string{T("hostname"), T("domain"), T("public_ip"), T("private_ip")})
			for _, v := range vs {
				vsTable.Add(utils.FormatStringPointer(v.Hostname),
					utils.FormatStringPointer(v.Domain),
					utils.FormatStringPointer(v.PrimaryIpAddress),
					utils.FormatStringPointer(v.PrimaryBackendIpAddress))
			}
			vsTable.Print()
			table.Add(T("virtual servers"), buf.String())
		}
	}

	if !c.IsSet("no-hardware") {
		hw := vlan.Hardware
		if len(hw) == 0 {
			table.Add(T("hardware"), T("none"))
		} else {
			buf := new(bytes.Buffer)
			hwTable := terminal.NewTable(buf, []string{T("hostname"), T("domain"), T("public_ip"), T("private_ip")})
			for _, h := range hw {
				hwTable.Add(utils.FormatStringPointer(h.Hostname),
					utils.FormatStringPointer(h.Domain),
					utils.FormatStringPointer(h.PrimaryIpAddress),
					utils.FormatStringPointer(h.PrimaryBackendIpAddress))
			}
			hwTable.Print()
			table.Add(T("hardware"), buf.String())
		}
	}

	table.Print()
	return nil
}

package vlan

import (
	"bytes"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	Vs             bool
	Hardware       bool
}

func NewDetailCommand(sl *metadata.SoftlayerCommand) *DetailCommand {
	thisCmd := &DetailCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "detail " + T("IDENTIFIER"),
		Short: T("Get details about a VLAN."),
		Long: T(`${COMMAND_NAME} sl vlan detail IDENTIFIER [OPTIONS]

EXAMPLE:
	${COMMAND_NAME} sl vlan detail 12345678	--no-vs --no-hardware
	This command shows details of vlan with ID 12345678, and not list virtual server or hardware server.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVar(&thisCmd.Vs, "no-vs", false, T("Hide virtual server listing"))
	cobraCmd.Flags().BoolVar(&thisCmd.Hardware, "no-hardware", false, T("Hide hardware listing"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DetailCommand) Run(args []string) error {
	id, err := utils.ResolveVlanId(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("VLAN ID")
	}

	outputFormat := cmd.GetOutputFlag()

	vlan, err := cmd.NetworkManager.GetVlan(id, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get VLAN: {{.VLANID}}.\n", map[string]interface{}{"VLANID": id}), err.Error(), 2)
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

	if !cmd.Vs {
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

	if !cmd.Hardware {
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
	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

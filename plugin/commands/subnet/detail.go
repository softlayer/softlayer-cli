package subnet

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
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
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	NoVs           bool
	NoHardware     bool
	NoIp           bool
	NoTag          bool
}

func NewDetailCommand(sl *metadata.SoftlayerCommand) *DetailCommand {
	thisCmd := &DetailCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "detail " + T("IDENTIFIER"),
		Short: T("Get details of a subnet"),
		Long: T(`${COMMAND_NAME} sl subnet detail IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl subnet detail 12345678 
   This command shows detailed information about subnet with ID 12345678, including virtual servers and hardware servers information.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVar(&thisCmd.NoVs, "no-vs", false, T("Hide virtual server listing"))
	cobraCmd.Flags().BoolVar(&thisCmd.NoHardware, "no-hardware", false, T("Hide hardware listing"))
	cobraCmd.Flags().BoolVar(&thisCmd.NoIp, "no-ip", false, T("Hide IP address listing"))
	cobraCmd.Flags().BoolVar(&thisCmd.NoTag, "no-Tag", false, T("Hide Tag listing"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DetailCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	subnetID, err := utils.ResolveSubnetId(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Subnet ID")
	}

	mask := "mask[endPointIpAddress[virtualGuest,hardware],ipAddresses[id, ipAddress,note,hardware,virtualGuest], datacenter, virtualGuests, hardware,networkVlan[networkSpace], tagReferences]"
	subnet, err := cmd.NetworkManager.GetSubnet(subnetID, mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get subnet: {{.ID}}.\n", map[string]interface{}{"ID": subnetID}), err.Error(), 2)
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
	table.Add(T("usable ips"), strconv.FormatFloat(float64(sl.Get(subnet.UsableIpAddressCount).(datatypes.Float64)), 'f', 0, 64))
	if !cmd.NoIp {
		if subnet.IpAddresses == nil || len(subnet.IpAddresses) == 0 {
			table.Add(T("IP address"), T("none"))
		} else {
			buf := new(bytes.Buffer)
			ipTable := terminal.NewTable(buf, []string{T("ID"), T("IP address"), T("Description"), T("Note")})
			endPointIpAddressDescription := "-"
			if subnet.EndPointIpAddress != nil {
				routedSubnet := fmt.Sprintf(
					"%s/%d",
					*subnet.EndPointIpAddress.Subnet.NetworkIdentifier,
					*subnet.EndPointIpAddress.Subnet.Cidr,
				)
				if subnet.EndPointIpAddress.Hardware != nil {
					routedSubnet = *subnet.EndPointIpAddress.Hardware.FullyQualifiedDomainName
				}
				if subnet.EndPointIpAddress.VirtualGuest != nil {
					routedSubnet = *subnet.EndPointIpAddress.VirtualGuest.FullyQualifiedDomainName
				}

				endPointIpAddressDescription = fmt.Sprintf(
					"Routed to %s → %s",
					*subnet.EndPointIpAddress.IpAddress,
					routedSubnet,
				)
			}
			for _, ip := range subnet.IpAddresses {
				description := "-"
				if subnet.EndPointIpAddress != nil {
					description = endPointIpAddressDescription
				}
				if ip.Hardware != nil {
					description = *ip.Hardware.FullyQualifiedDomainName
				}
				if ip.VirtualGuest != nil {
					description = *ip.VirtualGuest.FullyQualifiedDomainName
				}
				ipTable.Add(
					utils.FormatIntPointer(ip.Id),
					utils.FormatStringPointer(ip.IpAddress),
					description,
					utils.FormatStringPointer(ip.Note),
				)
			}
			ipTable.Print()
			table.Add(T("IP address"), buf.String())
		}
	}

	if !cmd.NoVs {
		if subnet.VirtualGuests == nil || len(subnet.VirtualGuests) == 0 {
			table.Add(T("virtual guests"), T("none"))
		} else {
			buf := new(bytes.Buffer)
			vsTable := terminal.NewTable(buf, []string{T("Hostname"), T("domain"), T("public_ip"), T("private_ip")})
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
	if !cmd.NoHardware {
		if subnet.Hardware == nil || len(subnet.Hardware) == 0 {
			table.Add(T("hardware"), T("none"))
		} else {
			buf := new(bytes.Buffer)
			hwTable := terminal.NewTable(buf, []string{T("Hostname"), T("domain"), T("public_ip"), T("private_ip")})
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

	if !cmd.NoTag {
		if subnet.TagReferences == nil || len(subnet.TagReferences) == 0 {
			table.Add(T("Tag"), T("none"))
		} else {
			buf := new(bytes.Buffer)
			vsTable := terminal.NewTable(buf, []string{T("ID")})
			for _, tag := range subnet.TagReferences {
				vsTable.Add(utils.FormatIntPointer(tag.TagId))
			}
			vsTable.Print()
			table.Add(T("Tag"), buf.String())
		}
	}
	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

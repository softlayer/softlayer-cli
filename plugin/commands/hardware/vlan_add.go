package hardware

import (
	"strconv"

	// "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"

	"github.com/softlayer/softlayer-go/datatypes"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type VlanAddCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	NetworkManager  managers.NetworkManager
	Command         *cobra.Command
}

func NewVlanAddCommand(sl *metadata.SoftlayerCommand) (cmd *VlanAddCommand) {
	thisCmd := &VlanAddCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "vlan-add " + T("IDENTIFIER") + " " + T("VLAN") + "...",
		Short: T("Trunks a VLAN to the this server."),
		Long: T(`IDENTIFIER is the id of the server
VLANS is the ID of the VLANs. Multiple vlans can be added at the same time.`),
		Args: metadata.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VlanAddCommand) Run(args []string) error {
	hardwareId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat := cmd.GetOutputFlag()
	h_mask := `mask[frontendNetworkComponents[id, name, port, macAddress, primaryIpAddress],
				    backendNetworkComponents[id, name, port, macAddress, primaryIpAddress],
				    id, hostname, domain]`
	v_mask := "mask[id, vlanNumber, networkSpace]"
	hardware, err := cmd.HardwareManager.GetHardware(hardwareId, h_mask)
	// I18N Subs go in here.
	subs := map[string]string{
		"ID":     args[0],
		"VLANID": "none",
	}
	if err != nil {
		return slErr.NewAPIError(T("Failed to get hardware server: {{.ID}}.\n", subs), err.Error(), 2)
	}
	// API will return an error if you try to add a public vlan to a private network component
	pub_vlans := []datatypes.Network_Vlan{}
	pri_vlans := []datatypes.Network_Vlan{}
	table := cmd.UI.Table([]string{T("Id"), T("VLAN"), T("Network")})
	for i := 1; i < len(args); i++ {
		vlan_id, err := strconv.Atoi(args[i])
		if err != nil {
			return slErr.NewInvalidSoftlayerIdInputError("VLAN ID")
		}
		i_vlan, err := cmd.NetworkManager.GetVlan(vlan_id, v_mask)
		if err != nil {
			subs["VLANID"] = args[i]
			return slErr.NewAPIError(T("Failed to get VLAN: {{.VLANID}}.\n", subs), err.Error(), 2)
		}
		if *i_vlan.NetworkSpace == "PUBLIC" {
			pub_vlans = append(pub_vlans, i_vlan)
		} else {
			pri_vlans = append(pri_vlans, i_vlan)
		}
	}
	// If we need to add vlans, find the first Frontend/Backend Network Component with a primary IP,
	// and add the appropriate vlans there.
	if len(pub_vlans) > 0 {
		for _, component := range hardware.FrontendNetworkComponents {
			if component.PrimaryIpAddress != nil {
				added_vlans, err := cmd.HardwareManager.TrunkVlans(*component.Id, pub_vlans)
				if err != nil {
					return err
				}
				for _, v := range added_vlans {
					table.Add(
						utils.FormatIntPointer(v.Id),
						utils.FormatIntPointer(v.VlanNumber),
						utils.FormatStringPointer(v.Name),
					)
				}
				break
			}
		}
	}
	if len(pri_vlans) > 0 {
		for _, component := range hardware.BackendNetworkComponents {
			if component.PrimaryIpAddress != nil {
				added_vlans, err := cmd.HardwareManager.TrunkVlans(*component.Id, pri_vlans)
				if err != nil {
					return err
				}
				for _, v := range added_vlans {
					table.Add(
						utils.FormatIntPointer(v.Id),
						utils.FormatIntPointer(v.VlanNumber),
						utils.FormatStringPointer(v.Name),
					)
				}
				break
			}
		}
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

package hardware

import (
	"strconv"

	// "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type VlanTrunkableCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	NetworkManager  managers.NetworkManager
	Command         *cobra.Command
}

func NewVlanTrunkableCommand(sl *metadata.SoftlayerCommand) (cmd *VlanTrunkableCommand) {
	thisCmd := &VlanTrunkableCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "vlan-trunkable " + T("IDENTIFIER"),
		Short: T("Lists VLANs this server can be given access to."),
		Long:  T("This command will only show VLANs not yet trunked to this server."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VlanTrunkableCommand) Run(args []string) error {
	hardwareId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat := cmd.GetOutputFlag()
	h_mask := `mask[networkComponents[id, name, port, macAddress, primaryIpAddress, 
						networkVlansTrunkable[id, name, vlanNumber, fullyQualifiedName,networkSpace]
			   ]]`
	hardware, err := cmd.HardwareManager.GetHardware(hardwareId, h_mask)
	// I18N Subs go in here.
	subs := map[string]string{
		"ID": args[0],
	}
	if err != nil {
		return slErr.NewAPIError(T("Failed to get hardware server: {{.ID}}.\n", subs), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Id"), T("Fully qualified name"), T("Name"), T("Network")})
	for _, component := range hardware.NetworkComponents {
		if component.PrimaryIpAddress != nil {
			for _, vlan := range component.NetworkVlansTrunkable {
				table.Add(
					utils.FormatIntPointer(vlan.Id),
					utils.FormatStringPointer(vlan.FullyQualifiedName),
					utils.FormatStringPointer(vlan.Name),
					utils.FormatStringPointer(vlan.NetworkSpace),
				)
			}
		}
	}
	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

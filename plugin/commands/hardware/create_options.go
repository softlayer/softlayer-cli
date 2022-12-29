package hardware

import (
	"sort"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateOptionsCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	NetworkManager  managers.NetworkManager
	Command         *cobra.Command
}

func NewCreateOptionsCommand(sl *metadata.SoftlayerCommand) (cmd *CreateOptionsCommand) {
	thisCmd := &CreateOptionsCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "create-options",
		Short: T("Server order options for a given chassis"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CreateOptionsCommand) Run(args []string) error {
	productPackage, err := cmd.HardwareManager.GetPackage()
	if err != nil {
		return errors.NewAPIError(T("Failed to get product package for hardware server."), err.Error(), 2)
	}
	options := cmd.HardwareManager.GetCreateOptions(productPackage)
	pods, err := cmd.NetworkManager.GetPods("", true)
	if err != nil {
		return errors.NewAPIError(T("Failed to get Pods."), err.Error(), 2)
	}
	//datacenters
	dcTable := cmd.UI.Table([]string{T("Data center"), T("Value"), T("Note")})
	locations := options[managers.KEY_LOCATIONS]
	var sortedLocations []string
	for key, _ := range locations {
		sortedLocations = append(sortedLocations, key)
	}
	sort.Strings(sortedLocations)
	for _, key := range sortedLocations {
		note := getPodWithClosedAnnouncement(locations[key], pods)
		dcTable.Add(locations[key], key, note)
	}
	dcTable.Print()
	cmd.UI.Print("")

	//preset
	presetTable := cmd.UI.Table([]string{T("Size"), T("Value")})
	presets := options[managers.KEY_SIZES]
	var sortedPresets []string
	for key, _ := range presets {
		sortedPresets = append(sortedPresets, key)
	}
	sort.Strings(sortedPresets)
	for _, key := range sortedPresets {
		presetTable.Add(presets[key], key)
	}
	presetTable.Print()
	cmd.UI.Print("")

	//operating system
	osTable := cmd.UI.Table([]string{T("Operating system"), T("Key"), T("Reference Code")})
	oses := options[managers.KEY_OS]
	nameOses := options[managers.KEY_NAME_OS]
	var sortedoses []string
	for key, _ := range oses {
		sortedoses = append(sortedoses, key)
	}
	sort.Strings(sortedoses)
	for _, key := range sortedoses {
		osTable.Add(oses[key], nameOses[key], key)
	}
	osTable.Print()
	cmd.UI.Print("")

	//port speed
	portTable := cmd.UI.Table([]string{T("Port speed"), T("Speed"), T("Key")})
	ports := options[managers.KEY_PORT_SPEED]
	portsDescription := options[managers.KEY_PORT_SPEED_DESCRIPTION]
	var sortedPorts []string
	for key, _ := range ports {
		sortedPorts = append(sortedPorts, key)
	}
	sort.Strings(sortedPorts)
	for _, key := range sortedPorts {
		portTable.Add(portsDescription[key], ports[key], key)
	}
	portTable.Print()
	cmd.UI.Print("")

	//extras
	extraTable := cmd.UI.Table([]string{T("Extras"), T("Value")})
	extras := options[managers.KEY_EXTRAS]
	var sortedExtras []string
	for key, _ := range extras {
		sortedExtras = append(sortedExtras, key)
	}
	sort.Strings(sortedExtras)
	for _, key := range sortedExtras {
		extraTable.Add(extras[key], key)
	}
	extraTable.Print()
	cmd.UI.Print("")

	//routers
	routerTable := cmd.UI.Table([]string{T("Routers"), T("Hostname"), T("Name")})
	routers, err := cmd.NetworkManager.GetRouters("")
	if err != nil {
		return errors.NewAPIError(T("Failed to get Routers."), err.Error(), 2)
	}
	for _, router := range routers {
		routerTable.Add(
			utils.FormatIntPointer(router.Id),
			utils.FormatStringPointer(router.Hostname),
			utils.FormatStringPointer(router.TopLevelLocation.LongName),
		)
	}
	routerTable.Print()
	return nil
}

func getPodWithClosedAnnouncement(key string, pods []datatypes.Network_Pod) string {
	for _, pod := range pods {
		if key == *pod.DatacenterLongName {
			return "closed soon: " + *pod.Name
		}
	}
	return "-"
}

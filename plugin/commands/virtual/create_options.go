package virtual

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
	"sort"
)

type CreateOptionsCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewCreateOptionsCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *CreateOptionsCommand) {
	return &CreateOptionsCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *CreateOptionsCommand) Run(c *cli.Context) error {
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	//createOptions, err := cmd.VirtualServerManager.GetCreateOptions("PUBLIC_CLOUD_SERVER", "dal13")
	createOptions, err := cmd.VirtualServerManager.GetCreateOptions("PUBLIC_CLOUD_SERVER", "")
	if err != nil {
		return cli.NewExitError(T("Failed to get virtual server creation options.\n")+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, createOptions)
	}

	table := cmd.UI.Table([]string{T("datacenter"), T("value")})
	locations := createOptions[managers.KEY_LOCATIONS]
	var sortedLocations []string
	for key, _ := range locations {
		sortedLocations = append(sortedLocations, key)
	}
	sort.Strings(sortedLocations)
	for _, key := range sortedLocations {
		table.Add(locations[key], key)
	}
	table.Print()

	//preset
	presetTable := cmd.UI.Table([]string{T("Size"), T("Value")})
	presets := createOptions[managers.KEY_SIZES]
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
	osTable := cmd.UI.Table([]string{T("Operating system"), T("Value")})
	oses := createOptions[managers.KEY_OS]
	var sortedoses []string
	for key, _ := range oses {
		sortedoses = append(sortedoses, key)
	}
	sort.Strings(sortedoses)
	for _, key := range sortedoses {
		osTable.Add(oses[key], key)
	}
	osTable.Print()
	cmd.UI.Print("")

	//port speed
	portTable := cmd.UI.Table([]string{T("Port speed"), T("Value")})
	ports := createOptions[managers.KEY_PORT_SPEED]
	var sortedPorts []string
	for key, _ := range ports {
		sortedPorts = append(sortedPorts, key)
	}
	sort.Strings(sortedPorts)
	for _, key := range sortedPorts {
		portTable.Add(ports[key], key)
	}
	portTable.Print()
	cmd.UI.Print("")

	//Disk
	diskTable := cmd.UI.Table([]string{T("disk_guest"), T("Value")})
	disks := createOptions[managers.KEY_GUEST]
	var sortedDisks []string
	for key, _ := range disks {
		sortedDisks = append(sortedDisks, key)
	}
	sort.Strings(sortedDisks)
	for _, key := range sortedDisks {
		diskTable.Add(disks[key],key)
	}
	diskTable.Print()
	cmd.UI.Print("")

	//extras
	extraTable := cmd.UI.Table([]string{T("Extras"), T("Value")})
	extras := createOptions[managers.KEY_EXTRAS]
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

	return nil
}

package hardware

import (
	"sort"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type CreateOptionsCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewCreateOptionsCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *CreateOptionsCommand) {
	return &CreateOptionsCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *CreateOptionsCommand) Run(c *cli.Context) error {
	productPackage, err := cmd.HardwareManager.GetPackage()
	if err != nil {
		return cli.NewExitError(T("Failed to get product package for hardware server.")+err.Error(), 2)
	}
	options := cmd.HardwareManager.GetCreateOptions(productPackage)
	//datacenters
	dcTable := cmd.UI.Table([]string{T("Data center"), T("Value")})
	locations := options[managers.KEY_LOCATIONS]
	var sortedLocations []string
	for key, _ := range locations {
		sortedLocations = append(sortedLocations, key)
	}
	sort.Strings(sortedLocations)
	for _, key := range sortedLocations {
		dcTable.Add(locations[key], key)
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
	osTable := cmd.UI.Table([]string{T("Operating system"), T("Value")})
	oses := options[managers.KEY_OS]
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
	ports := options[managers.KEY_PORT_SPEED]
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
	return nil
}

func HardwareCreateOptionsMetaData() cli.Command {
	return cli.Command{
		Category:    "hardware",
		Name:        "create-options",
		Description: T("Server order options for a given chassis"),
		Usage:       "${COMMAND_NAME} sl hardware create-options",
	}
}

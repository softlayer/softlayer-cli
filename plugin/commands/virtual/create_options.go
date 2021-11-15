package virtual

import (
	"sort"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
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
	datacenters := createOptions.Locations
	for _, location := range datacenters {
		table.Add(
			utils.FormatStringPointer(location.LongName),
			utils.FormatStringPointer(location.Name),
		)
	}
	table.Print()

	table = cmd.UI.Table([]string{T("Size"), T("value")})
	sizes := createOptions.Sizes
	for _, size := range sizes {
		table.Add(
			utils.FormatStringPointer(size.Description),
			utils.FormatStringPointer(size.KeyName),
		)
	}
	table.Print()

	table = cmd.UI.Table([]string{T("OS"), T("key"), T("Reference Code")})
	osList := createOptions.OperatingSystems
	sort.SliceStable(osList, func(i, j int) bool {
		return utils.FormatStringPointer(osList[i].Description) < utils.FormatStringPointer(osList[j].Description)
	})

	for _, os := range osList {
		table.Add(
			utils.FormatStringPointer(os.Description),
			utils.FormatStringPointer(os.KeyName),
			utils.FormatStringPointer(os.SoftwareDescription.ReferenceCode),
		)
	}
	table.Print()

	table = cmd.UI.Table([]string{T("network"), T("key")})
	networkList := createOptions.PortSpeed

	for _, network := range networkList {
		table.Add(
			utils.FormatStringPointer(network.Description),
			utils.FormatStringPointer(network.KeyName),
		)
	}
	table.Print()

	table = cmd.UI.Table([]string{T("database"), T("key")})
	databaseList := createOptions.PortSpeed
	sort.SliceStable(databaseList, func(i, j int) bool {
		return utils.FormatStringPointer(databaseList[i].Description) < utils.FormatStringPointer(databaseList[j].Description)
	})

	for _, database := range databaseList {
		table.Add(
			utils.FormatStringPointer(database.Description),
			utils.FormatStringPointer(database.KeyName),
		)
	}
	table.Print()

	table = cmd.UI.Table([]string{T("guest disk"), T("key"), T("capacity"), T("disk")})
	guestDisks := createOptions.GuestDisk

	for _, disk := range guestDisks {
		table.Add(
			utils.FormatStringPointer(disk.Description),
			utils.FormatStringPointer(disk.KeyName),
			utils.FormatSLFloatPointerToInt(disk.Capacity),
			utils.FormatStringPointer(disk.LongDescription),
		)
	}

	table.Print()

	return nil
}

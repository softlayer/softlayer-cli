package order

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type PackageLocationCommand struct {
	UI           terminal.UI
	OrderManager managers.OrderManager
}

func NewPackageLocationCommand(ui terminal.UI, orderManager managers.OrderManager) (cmd *PackageLocationCommand) {
	return &PackageLocationCommand{
		UI:           ui,
		OrderManager: orderManager,
	}
}

func (cmd *PackageLocationCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	packageKeyname := c.Args()[0]

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	locations, err := cmd.OrderManager.PackageLocation(packageKeyname)
	if err != nil {
		return cli.NewExitError(T("Failed to list package locations.\n")+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, locations)
	}

	cmd.Print(locations)
	return nil
}

func (cmd *PackageLocationCommand) Print(locations []datatypes.Location_Region) {
	table := cmd.UI.Table([]string{T("id"), T("dc"), T("description"), T("keyName")})

	for _, region := range locations {
		for _, datacenter := range region.Locations {
			table.Add(utils.FormatIntPointer(datacenter.Location.Id),
				utils.FormatStringPointer(datacenter.Location.Name),
				utils.FormatStringPointer(region.Description),
				utils.FormatStringPointer(region.Keyname))
		}
	}
	table.Print()
}

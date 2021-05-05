package order

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PackageListCommand struct {
	UI           terminal.UI
	OrderManager managers.OrderManager
}

func NewPackageListCommand(ui terminal.UI, orderManager managers.OrderManager) (cmd *PackageListCommand) {
	return &PackageListCommand{
		UI:           ui,
		OrderManager: orderManager,
	}
}

func (cmd *PackageListCommand) Run(c *cli.Context) error {

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	keyword := c.String("keyword")
	packageType := c.String("package-type")

	packages, err := cmd.OrderManager.ListPackage(keyword, packageType)
	if err != nil {
		return cli.NewExitError(T("Failed to list packages.\n")+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, packages)
	}

	cmd.Print(packages)
	return nil
}

func (cmd *PackageListCommand) Print(packages []datatypes.Product_Package) {
	table := cmd.UI.Table([]string{T("id"), T("name"), T("keyName"), T("type")})

	for _, pac := range packages {
		table.Add(utils.FormatIntPointer(pac.Id),
			utils.FormatStringPointer(pac.Name),
			utils.FormatStringPointer(pac.KeyName),
			utils.FormatStringPointer(pac.Type.KeyName))
	}
	table.Print()
}

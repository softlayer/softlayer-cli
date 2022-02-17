package order

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CategoryListCommand struct {
	UI           terminal.UI
	OrderManager managers.OrderManager
}

func NewCategoryListCommand(ui terminal.UI, orderManager managers.OrderManager) (cmd *CategoryListCommand) {
	return &CategoryListCommand{
		UI:           ui,
		OrderManager: orderManager,
	}
}

func (cmd *CategoryListCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	packageKeyname := c.Args()[0]

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	categories, err := cmd.OrderManager.ListCategories(packageKeyname)
	if err != nil {
		return cli.NewExitError(T("Failed to list categories.\n")+err.Error(), 2)
	}
	var CategoriesRequired []datatypes.Product_Package_Order_Configuration
	if c.Bool("required") {
		for _, cat := range categories {
			if *cat.IsRequired != 0 {
				CategoriesRequired = append(CategoriesRequired, cat)
			}
		}
		categories = CategoriesRequired
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, categories)
	}

	cmd.Print(categories)
	return nil
}

func (cmd *CategoryListCommand) Print(categories []datatypes.Product_Package_Order_Configuration) {
	table := cmd.UI.Table([]string{T("Name"), T("Category Code"), T("Is Required")})

	for _, cat := range categories {
		var isRequired string
		if *cat.IsRequired != 0 {
			isRequired = "Y"
		} else {
			isRequired = "N"
		}
		table.Add(utils.FormatStringPointer(cat.ItemCategory.Name),
			utils.FormatStringPointer(cat.ItemCategory.CategoryCode),
			isRequired)
	}
	table.Print()
}

func OrderCategoryListMetaData() cli.Command {
	return cli.Command{
		Category:    "order",
		Name:        "category-list",
		Description: T("List the categories of a package"),
		Usage: T(`${COMMAND_NAME} sl order category-list [OPTIONS] PACKAGE_KEYNAME
	
EXAMPLE: 
   ${COMMAND_NAME} sl order category-list BARE_METAL_SERVER
   This command lists the categories of Bare Metal servers.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "required",
				Usage: T("List only the required categories for the package"),
			},
			metadata.OutputFlag(),
		},
	}
}

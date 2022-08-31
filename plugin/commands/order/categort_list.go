package order

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CategoryListCommand struct {
	*metadata.SoftlayerCommand
	OrderManager managers.OrderManager
	Command      *cobra.Command
	Required     bool
}

func NewCategoryListCommand(sl *metadata.SoftlayerCommand) (cmd *CategoryListCommand) {
	thisCmd := &CategoryListCommand{
		SoftlayerCommand: sl,
		OrderManager:     managers.NewOrderManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "category-list " + T("PACKAGE_KEYNAME"),
		Short: T("List the categories of a package"),
		Long: T(`
EXAMPLE: 
	${COMMAND_NAME} sl order category-list BARE_METAL_SERVER --required`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVar(&thisCmd.Required, "required", false, T("List only the required categories for the package"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CategoryListCommand) Run(args []string) error {
	packageKeyname := args[0]

	outputFormat := cmd.GetOutputFlag()

	categories, err := cmd.OrderManager.ListCategories(packageKeyname)
	if err != nil {
		return errors.NewAPIError(T("Failed to list categories.\n"), err.Error(), 2)
	}
	var CategoriesRequired []datatypes.Product_Package_Order_Configuration
	if cmd.Required {
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

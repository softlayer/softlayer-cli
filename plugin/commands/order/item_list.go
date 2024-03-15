package order

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

type ItemListCommand struct {
	*metadata.SoftlayerCommand
	OrderManager managers.OrderManager
	Command      *cobra.Command
	Keyword      string
	Category     string
}

func NewItemListCommand(sl *metadata.SoftlayerCommand) (cmd *ItemListCommand) {
	thisCmd := &ItemListCommand{
		SoftlayerCommand: sl,
		OrderManager:     managers.NewOrderManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "item-list " + T("PACKAGE_KEYNAME"),
		Short: T("List package items that are used for ordering"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Keyword, "keyword", "", T("A word (or string) that is used to filter item names"))
	cobraCmd.Flags().StringVar(&thisCmd.Category, "category", "", T("Category code that is used to filter items"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ItemListCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	packageKeyname := args[0]

	keyword := cmd.Keyword
	category := cmd.Category

	items, err := cmd.OrderManager.ListItems(packageKeyname, keyword, category)
	if err != nil {
		return errors.NewAPIError(T("Failed to list items.\n"), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, items)
	}

	cmd.Print(items)
	return nil
}

// """sorts the items into a dictionary of categories, with a list of items"""
func sortItems(items []datatypes.Product_Item) map[string][]datatypes.Product_Item {

	sortedItems := make(map[string][]datatypes.Product_Item)

	for _, item := range items {
		category := item.ItemCategory.CategoryCode
		if _, ok := sortedItems[*category]; !ok {
			sortedItems[*category] = nil
		}
		sortedItems[*category] = append(sortedItems[*category], item)
	}
	return sortedItems
}

func (cmd *ItemListCommand) Print(items []datatypes.Product_Item) {
	table := cmd.UI.Table([]string{T("category"), T("Key Name"), T("Description")})

	sortedItems := sortItems(items)

	var keys []string
	for k := range sortedItems {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		for _, item := range sortedItems[k] {
			table.Add(k,
				utils.FormatStringPointer(item.KeyName),
				utils.FormatStringPointer(item.Description))
		}
	}
	table.Print()
}

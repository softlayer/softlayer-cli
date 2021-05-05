package order

import (
	"sort"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ItemListCommand struct {
	UI           terminal.UI
	OrderManager managers.OrderManager
}

func NewItemListCommand(ui terminal.UI, orderManager managers.OrderManager) (cmd *ItemListCommand) {
	return &ItemListCommand{
		UI:           ui,
		OrderManager: orderManager,
	}
}

func (cmd *ItemListCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	packageKeyname := c.Args()[0]

	keyword := c.String("keyword")
	category := c.String("category")

	items, err := cmd.OrderManager.ListItems(packageKeyname, keyword, category)
	if err != nil {
		return cli.NewExitError(T("Failed to list items.\n")+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, items)
	}

	cmd.Print(items)
	return nil
}

//"""sorts the items into a dictionary of categories, with a list of items"""
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

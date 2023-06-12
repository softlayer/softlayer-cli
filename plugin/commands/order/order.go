package order

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "order",
		Short: T("Classic infrastructure Orders"),
		RunE:  nil,
	}

	cobraCmd.AddCommand(NewCategoryListCommand(sl).Command)
	cobraCmd.AddCommand(NewItemListCommand(sl).Command)
	cobraCmd.AddCommand(NewPackageListCommand(sl).Command)
	cobraCmd.AddCommand(NewPackageLocationCommand(sl).Command)
	cobraCmd.AddCommand(NewPlaceCommand(sl).Command)
	cobraCmd.AddCommand(NewPlaceQuoteCommand(sl).Command)
	cobraCmd.AddCommand(NewPresetListCommand(sl).Command)
	cobraCmd.AddCommand(NewQuoteListCommand(sl).Command)
	cobraCmd.AddCommand(NewQuoteDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewQuoteSaveCommand(sl).Command)
	cobraCmd.AddCommand(NewQuoteCommand(sl).Command)
	cobraCmd.AddCommand(NewQuoteDeleteCommand(sl).Command)
	cobraCmd.AddCommand(NewLookupCommand(sl).Command)
	return cobraCmd
}

func OrderNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "order",
		Description: T("Classic infrastructure Orders"),
	}
}

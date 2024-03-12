package search

import (

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SearchCommand struct {
	*metadata.SoftlayerCommand
	SearchManager managers.SearchManager
	Command       *cobra.Command
	Query		  string
}

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	thisCmd := &SearchCommand{
		SoftlayerCommand: sl,
		SearchManager:    managers.NewSearchManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "search",
		Short: T("Perform a query against the SoftLayer search database."),
		Long: T(`Read More: https://sldn.softlayer.com/reference/services/SoftLayer_Search/search/
Examples::

    sl search --query 'test.com'
    sl search --query '_objectType:SoftLayer_Virtual_Guest test.com'
`),
		Args:  metadata.NoArgs,
		RunE:  func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.Query, "query", "q", "", T("The search query you want to use."))
	cobraCmd.AddCommand(NewSearchTypesCommand(sl).Command)
	return cobraCmd
}

func SearchNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "search",
		Description: T("Perform a query against the SoftLayer search database."),
	}
}

func (cmd *SearchCommand) Run(args []string) error {

	results, err := cmd.SearchManager.AdvancedSearch("", cmd.Query)
	if err != nil { return err}

	return utils.PrintPrettyJSON(cmd.UI, results)
}
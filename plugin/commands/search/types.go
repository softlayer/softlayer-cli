package search

import (
	"bytes"
	"github.com/spf13/cobra"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SearchTypesCommand struct {
	*metadata.SoftlayerCommand
	SearchManager managers.SearchManager
	Command       *cobra.Command
}

func NewSearchTypesCommand(sl *metadata.SoftlayerCommand) *SearchTypesCommand {
	thisCmd := &SearchTypesCommand{
		SoftlayerCommand: sl,
		SearchManager:    managers.NewSearchManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "types",
		Short: T("Display searchable types."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd

}

func (cmd *SearchTypesCommand) Run(args []string) error {

	type_results, err := cmd.SearchManager.GetTypes()
	if err != nil {
		return err
	}
	outputFormat := cmd.GetOutputFlag()
	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, type_results)
	}
	table := cmd.UI.Table([]string{T("Name"), T("Properties")})
	for _, search_type := range type_results {
		sub_buf := new(bytes.Buffer)
		sub_table := terminal.NewTable(sub_buf, []string{T("Property"), T("Sortable"), T("Type")})
		for _, t_prop := range search_type.Properties {
			sub_table.Add(*t_prop.Name, strconv.FormatBool(*t_prop.SortableFlag), *t_prop.Type)
		}
		sub_table.Print()
		table.Add(*search_type.Name, sub_buf.String())
	}
	table.Print()
	return nil
}

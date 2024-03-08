package search

import (
	"fmt"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

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

	fmt.Printf("Search Types would go here\n")
	return nil
}

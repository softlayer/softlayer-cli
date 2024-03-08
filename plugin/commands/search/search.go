package search

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "search",
		Short: T("Perform a query against the SoftLayer search database."),
		RunE:  nil,
	}

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

package tags

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "tags",
		Short: T("Classic infrastructure Tag management"),
		RunE:  nil,
	}

	cobraCmd.AddCommand(NewDeleteCommand(sl).Command)
	cobraCmd.AddCommand(NewDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewListCommand(sl).Command)
	cobraCmd.AddCommand(NewSetCommand(sl).Command)
	cobraCmd.AddCommand(NewCleanupCommand(sl).Command)
	return cobraCmd
}

func TagsNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "tags",
		Description: T("Classic infrastructure Tag management"),
	}
}

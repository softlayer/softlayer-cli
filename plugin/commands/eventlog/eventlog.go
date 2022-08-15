package eventlog

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "event-log",
		Short: T("Classic infrastructure Event Log Group"),
		RunE:  nil,
	}

	cobraCmd.AddCommand(NewGetCommand(sl).Command)
	cobraCmd.AddCommand(NewTypesCommand(sl).Command)
	return cobraCmd
}

func EventLogNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "event-log",
		Description: T("Classic infrastructure Event Log Group"),
	}
}

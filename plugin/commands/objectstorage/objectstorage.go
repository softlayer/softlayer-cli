package objectstorage

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "object-storage",
		Short: T("Classic infrastructure Object Storage commands"),
		RunE:  nil,
	}
	cobraCmd.AddCommand(NewAccountsCommand(sl).Command)
	cobraCmd.AddCommand(NewCredentialCreateCommand(sl).Command)
	cobraCmd.AddCommand(NewCredentialDeleteCommand(sl).Command)
	cobraCmd.AddCommand(NewCredentialLimitCommand(sl).Command)
	cobraCmd.AddCommand(NewCredentialListCommand(sl).Command)
	cobraCmd.AddCommand(NewEndpointsCommand(sl).Command)
	return cobraCmd
}

func ObjectStorageNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "object-storage",
		Description: T("Classic infrastructure Object Storage commands"),
	}
}

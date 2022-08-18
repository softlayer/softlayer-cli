package licenses

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "licenses",
		Short: T("Classic infrastructure Licenses"),
		RunE:  nil,
	}
	cobraCmd.AddCommand(NewLicensesOptionsCommand(sl).Command)
	cobraCmd.AddCommand(NewCreateCommand(sl).Command)
	cobraCmd.AddCommand(NewCancelItemCommand(sl).Command)
	return cobraCmd
}

func LicensesNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "licenses",
		Description: T("Classic infrastructure Licenses"),
	}
}

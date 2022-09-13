package security

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "security",
		Short: T("Classic infrastructure SSH Keys and SSL Certificates"),
		RunE:  nil,
	}
	cobraCmd.AddCommand(NewCertAddCommand(sl).Command)
	cobraCmd.AddCommand(NewCertDownloadCommand(sl).Command)
	cobraCmd.AddCommand(NewCertEditCommand(sl).Command)
	cobraCmd.AddCommand(NewCertListCommand(sl).Command)
	cobraCmd.AddCommand(NewCertRemoveCommand(sl).Command)
	cobraCmd.AddCommand(NewKeyAddCommand(sl).Command)
	cobraCmd.AddCommand(NewKeyEditCommand(sl).Command)
	cobraCmd.AddCommand(NewKeyListCommand(sl).Command)
	cobraCmd.AddCommand(NewKeyPrintCommand(sl).Command)
	cobraCmd.AddCommand(NewKeyRemoveCommand(sl).Command)

	return cobraCmd
}

func SecurityNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "security",
		Aliases:     []string{"ssl", "sshkey"},
		Description: T("Classic infrastructure SSH Keys and SSL Certificates"),
	}
}

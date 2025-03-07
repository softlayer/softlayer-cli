package cdn

import (

	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "cdn",
		Short: T("Classic infrastructure CDN commands") + " " + T("Deprecated"),
		RunE:  nil,
		Deprecated: "https://cloud.ibm.com/docs/CDN?topic=CDN-cdn-deprecation",
	}

	return cobraCmd
}

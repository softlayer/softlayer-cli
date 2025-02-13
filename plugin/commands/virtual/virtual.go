package virtual

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func VSNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "vs",
		Description: T("Classic infrastructure Virtual Servers"),
	}
}

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "vs",
		Short: T("Classic infrastructure Virtual Servers"),
		RunE:  nil,
	}
	cobraCmd.AddCommand(NewAuthorizeStorageCommand(sl).Command)
	cobraCmd.AddCommand(NewBandwidthCommand(sl).Command)
	cobraCmd.AddCommand(NewBillingCommand(sl).Command)
	cobraCmd.AddCommand(NewCancelCommand(sl).Command)
	cobraCmd.AddCommand(NewCapacityCreateCommand(sl).Command)
	cobraCmd.AddCommand(NewCapacityCreateOptionsCommand(sl).Command)
	cobraCmd.AddCommand(NewCapacityDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewCapacityListCommand(sl).Command)
	cobraCmd.AddCommand(NewCaptureCommand(sl).Command)
	cobraCmd.AddCommand(NewCreateCommand(sl).Command)
	cobraCmd.AddCommand(NewCreateHostCommand(sl).Command)
	cobraCmd.AddCommand(NewCreateOptionsCommand(sl).Command)
	cobraCmd.AddCommand(NewCredentialsCommand(sl).Command)
	cobraCmd.AddCommand(NewDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewDnsSyncCommand(sl).Command)
	cobraCmd.AddCommand(NewEditCommand(sl).Command)
	cobraCmd.AddCommand(NewListCommand(sl).Command)
	cobraCmd.AddCommand(NewListHostCommand(sl).Command)
	cobraCmd.AddCommand(NewMigrateCommand(sl).Command)
	cobraCmd.AddCommand(NewMonitoringListCommand(sl).Command)
	cobraCmd.AddCommand(NewPauseCommand(sl).Command)
	cobraCmd.AddCommand(NewPowerOffCommand(sl).Command)
	cobraCmd.AddCommand(NewPowerOnCommand(sl).Command)
	cobraCmd.AddCommand(NewReadyCommand(sl).Command)
	cobraCmd.AddCommand(NewRebootCommand(sl).Command)
	cobraCmd.AddCommand(NewReloadCommand(sl).Command)
	cobraCmd.AddCommand(NewRescueCommand(sl).Command)
	cobraCmd.AddCommand(NewResumeCommand(sl).Command)
	cobraCmd.AddCommand(NewStorageCommand(sl).Command)
	cobraCmd.AddCommand(NewUpgradeCommand(sl).Command)
	cobraCmd.AddCommand(NewUsageCommand(sl).Command)
	cobraCmd.AddCommand(NewNotifiactionsCommand(sl).Command)
	cobraCmd.AddCommand(NewNotificationsDeleteCommand(sl).Command)
	cobraCmd.AddCommand(NewNotificationsAddCommand(sl).Command)
	cobraCmd.AddCommand(NewOsAvailableCommand(sl).Command)

	return cobraCmd
}

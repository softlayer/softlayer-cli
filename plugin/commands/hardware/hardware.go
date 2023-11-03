package hardware

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "hardware",
		Short: T("Classic infrastructure hardware servers"),
		RunE:  nil,
	}

	cobraCmd.AddCommand(NewAuthorizeStorageCommand(sl).Command)
	cobraCmd.AddCommand(NewBillingCommand(sl).Command)
	cobraCmd.AddCommand(NewCancelCommand(sl).Command)
	cobraCmd.AddCommand(NewCancelReasonsCommand(sl).Command)
	cobraCmd.AddCommand(NewCreateCommand(sl).Command)
	cobraCmd.AddCommand(NewCreateOptionsCommand(sl).Command)
	cobraCmd.AddCommand(NewCredentialsCommand(sl).Command)
	cobraCmd.AddCommand(NewDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewEditCommand(sl).Command)
	cobraCmd.AddCommand(NewListCommand(sl).Command)
	cobraCmd.AddCommand(NewPowerCycleCommand(sl).Command)
	cobraCmd.AddCommand(NewPowerOffCommand(sl).Command)
	cobraCmd.AddCommand(NewPowerOnCommand(sl).Command)
	cobraCmd.AddCommand(NewRebootCommand(sl).Command)
	cobraCmd.AddCommand(NewReloadCommand(sl).Command)
	cobraCmd.AddCommand(NewRescueCommand(sl).Command)
	cobraCmd.AddCommand(NewUpdateFirmwareCommand(sl).Command)
	cobraCmd.AddCommand(NewToggleIPMICommand(sl).Command)
	cobraCmd.AddCommand(NewBandwidthCommand(sl).Command)
	cobraCmd.AddCommand(NewStorageCommand(sl).Command)
	cobraCmd.AddCommand(NewMonitoringListCommand(sl).Command)
	cobraCmd.AddCommand(NewSensorCommand(sl).Command)
	cobraCmd.AddCommand(NewReflashFirmwareCommand(sl).Command)
	cobraCmd.AddCommand(NewNotificationsCommand(sl).Command)
	cobraCmd.AddCommand(NewNotificationsDeleteCommand(sl).Command)
	cobraCmd.AddCommand(NewNotificationsAddCommand(sl).Command)
	cobraCmd.AddCommand(NewCreateCredentialCommand(sl).Command)
	cobraCmd.AddCommand(NewVlanAddCommand(sl).Command)
	cobraCmd.AddCommand(NewVlanRemoveCommand(sl).Command)
	return cobraCmd
}

func HardwareNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "hardware",
		Description: T("Classic infrastructure hardware servers"),
	}
}

package hardware

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	hardwareManager := managers.NewHardwareServerManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"hardware-authorize-storage": func(c *cli.Context) error {
			return NewAuthorizeStorageCommand(ui, hardwareManager).Run(c)
		},
		"hardware-billing": func(c *cli.Context) error {
			return NewBillingCommand(ui, hardwareManager).Run(c)
		},
		"hardware-cancel": func(c *cli.Context) error {
			return NewCancelCommand(ui, hardwareManager).Run(c)
		},
		"hardware-cancel-reasons": func(c *cli.Context) error {
			return NewCancelReasonsCommand(ui, hardwareManager).Run(c)
		},
		"hardware-create": func(c *cli.Context) error {
			return NewCreateCommand(ui, hardwareManager, context).Run(c)
		},
		"hardware-create-options": func(c *cli.Context) error {
			return NewCreateOptionsCommand(ui, hardwareManager).Run(c)
		},
		"hardware-credentials": func(c *cli.Context) error {
			return NewCredentialsCommand(ui, hardwareManager).Run(c)
		},
		"hardware-detail": func(c *cli.Context) error {
			return NewDetailCommand(ui, hardwareManager).Run(c)
		},
		"hardware-edit": func(c *cli.Context) error {
			return NewEditCommand(ui, hardwareManager).Run(c)
		},
		"hardware-list": func(c *cli.Context) error {
			return NewListCommand(ui, hardwareManager).Run(c)
		},
		"hardware-power-cycle": func(c *cli.Context) error {
			return NewPowerCycleCommand(ui, hardwareManager).Run(c)
		},
		"hardware-power-off": func(c *cli.Context) error {
			return NewPowerOffCommand(ui, hardwareManager).Run(c)
		},
		"hardware-power-on": func(c *cli.Context) error {
			return NewPowerOnCommand(ui, hardwareManager).Run(c)
		},
		"hardware-reboot": func(c *cli.Context) error {
			return NewRebootCommand(ui, hardwareManager).Run(c)
		},
		"hardware-reload": func(c *cli.Context) error {
			return NewReloadCommand(ui, hardwareManager).Run(c)
		},
		"hardware-rescue": func(c *cli.Context) error {
			return NewRescueCommand(ui, hardwareManager).Run(c)
		},
		"hardware-update-firmware": func(c *cli.Context) error {
			return NewUpdateFirmwareCommand(ui, hardwareManager).Run(c)
		},
		"hardware-toggle-ipmi": func(c *cli.Context) error {
			return NewToggleIPMICommand(ui, hardwareManager).Run(c)
		},
		"hardware-bandwidth": func(c *cli.Context) error {
			return NewBandwidthCommand(ui, hardwareManager).Run(c)
		},
		"hardware-storage": func(c *cli.Context) error {
			return NewStorageCommand(ui, hardwareManager).Run(c)
		},
		"hardware-guests": func(c *cli.Context) error {
			return NewGuestsCommand(ui, hardwareManager).Run(c)
		},
		"hardware-monitoring-list": func(c *cli.Context) error {
			return NewMonitoringListCommand(ui, hardwareManager).Run(c)
		},
		"hardware-sensor": func(c *cli.Context) error {
			return NewSensorCommand(ui, hardwareManager).Run(c)
		},
	}
	return CommandActionBindings
}

func HardwareNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "hardware",
		Description: T("Classic infrastructure hardware servers"),
	}
}

func HardwareMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "hardware",
		Description: T("Classic infrastructure hardware servers"),
		Usage:       "${COMMAND_NAME} sl hardware",
		Subcommands: []cli.Command{
			HardwareAuthorizeStorageMetaData(),
			HardwareBillingMetaData(),
			HardwareCancelMetaData(),
			HardwareCancelReasonsMetaData(),
			HardwareCreateMetaData(),
			HardwareCreateOptionsMetaData(),
			HardwareCredentialsMetaData(),
			HardwareDetailMetaData(),
			HardwareEditMetaData(),
			HardwareListMetaData(),
			HardwarePowerCycleMetaData(),
			HardwarePowerOffMetaData(),
			HardwarePowerOnMetaData(),
			HardwarePowerRebootMetaData(),
			HardwareReloadMetaData(),
			HardwareRescueMetaData(),
			HardwareUpdateFirmwareMetaData(),
			HardwareToggleIPMIMetaData(),
			HardwareBandwidthMetaData(),
			HardwareStorageMetaData(),
			HardwareGuestsMetaData(),
			HardwareMonitoringListMetaData(),
			HardwareSensorMetaData(),
		},
	}
}

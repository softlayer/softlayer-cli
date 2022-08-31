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
		Long:  "${COMMAND_NAME} sl vs",
		RunE:  nil,
	}
	cobraCmd.AddCommand(NewAuthorizeStorageCommand(sl).Command)
	return cobraCmd
}

/*
func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	virtualServerManager := managers.NewVirtualServerManager(session)
	imageManager := managers.NewImageManager(session)
	networkManager := managers.NewNetworkManager(session)
	dnsManager := managers.NewDNSManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{

		"vs-cancel": func(c *cli.Context) error {
			return NewCancelCommand(ui, virtualServerManager).Run(c)
		},
		"vs-capture": func(c *cli.Context) error {
			return NewCaptureCommand(ui, virtualServerManager).Run(c)
		},
		"vs-create": func(c *cli.Context) error {
			return NewCreateCommand(ui, virtualServerManager, imageManager, context).Run(c)
		},
		"vs-host-create": func(c *cli.Context) error {
			return NewCreateHostCommand(ui, virtualServerManager, networkManager, context).Run(c)
		},
		"vs-options": func(c *cli.Context) error {
			return NewCreateOptionsCommand(ui, virtualServerManager).Run(c)
		},
		"vs-credentials": func(c *cli.Context) error {
			return NewCredentialsCommand(ui, virtualServerManager).Run(c)
		},
		"vs-detail": func(c *cli.Context) error {
			return NewDetailCommand(ui, virtualServerManager).Run(c)
		},
		"vs-dns-sync": func(c *cli.Context) error {
			return NewDnsSyncCommand(ui, virtualServerManager, dnsManager).Run(c)
		},
		"vs-edit": func(c *cli.Context) error {
			return NewEditCommand(ui, virtualServerManager).Run(c)
		},
		"vs-list": func(c *cli.Context) error {
			return NewListCommand(ui, virtualServerManager).Run(c)
		},
		"vs-host-list": func(c *cli.Context) error {
			return NewListHostCommand(ui, virtualServerManager).Run(c)
		},
		"vs-migrate": func(c *cli.Context) error {
			return NewMigrageCommand(ui, virtualServerManager).Run(c)
		},
		"vs-pause": func(c *cli.Context) error {
			return NewPauseCommand(ui, virtualServerManager).Run(c)
		},
		"vs-power-off": func(c *cli.Context) error {
			return NewPowerOffCommand(ui, virtualServerManager).Run(c)
		},
		"vs-power-on": func(c *cli.Context) error {
			return NewPowerOnCommand(ui, virtualServerManager).Run(c)
		},
		"vs-ready": func(c *cli.Context) error {
			return NewReadyCommand(ui, virtualServerManager).Run(c)
		},
		"vs-billing": func(c *cli.Context) error {
			return NewBillingCommand(ui, virtualServerManager).Run(c)
		},
		"vs-reboot": func(c *cli.Context) error {
			return NewRebootCommand(ui, virtualServerManager).Run(c)
		},
		"vs-reload": func(c *cli.Context) error {
			return NewReloadCommand(ui, virtualServerManager, context).Run(c)
		},
		"vs-rescue": func(c *cli.Context) error {
			return NewRescueCommand(ui, virtualServerManager).Run(c)
		},
		"vs-resume": func(c *cli.Context) error {
			return NewResumeCommand(ui, virtualServerManager).Run(c)
		},
		"vs-upgrade": func(c *cli.Context) error {
			return NewUpgradeCommand(ui, virtualServerManager).Run(c)
		},
		"vs-capacity-create-options": func(c *cli.Context) error {
			return NewCapacityCreateOptiosCommand(ui, virtualServerManager).Run(c)
		},
		"vs-capacity-detail": func(c *cli.Context) error {
			return NewCapacityDetailCommand(ui, virtualServerManager).Run(c)
		},
		"vs-bandwidth": func(c *cli.Context) error {
			return NewBandwidthCommand(ui, virtualServerManager).Run(c)
		},
		"vs-storage": func(c *cli.Context) error {
			return NewStorageCommand(ui, virtualServerManager).Run(c)
		},
		"vs-placementgroup-list": func(c *cli.Context) error {
			return NewPlacementGroupListCommand(ui, virtualServerManager).Run(c)
		},
		"vs-placementgroup-create-options": func(c *cli.Context) error {
			return NewPlacementGruopCreateOptionsCommand(ui, virtualServerManager).Run(c)
		},
		"vs-placementgroup-create": func(c *cli.Context) error {
			return NewVSPlacementGroupCreateCommand(ui, virtualServerManager, context).Run(c)
		},
		"vs-capacity-list": func(c *cli.Context) error {
			return NewCapacityListCommand(ui, virtualServerManager).Run(c)
		},
		"vs-capacity-create": func(c *cli.Context) error {
			return NewCapacityCreateCommand(ui, virtualServerManager, context).Run(c)
		},
		"vs-usage": func(c *cli.Context) error {
			return NewUsageCommand(ui, virtualServerManager).Run(c)
		},
		"vs-placementgroup-details": func(c *cli.Context) error {
			return NewPlacementGroupDetailsCommand(ui, virtualServerManager).Run(c)
		},
		"vs-monitoring-list": func(c *cli.Context) error {
			return NewMonitoringListCommand(ui, virtualServerManager).Run(c)
		},
	}

	return CommandActionBindings
}

*/

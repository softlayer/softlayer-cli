package user

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "user",
		Short: T("Classic infrastructure Manage Users"),
		RunE:  nil,
	}

	cobraCmd.AddCommand(NewCreateCommand(sl).Command)
	cobraCmd.AddCommand(NewListCommand(sl).Command)
	cobraCmd.AddCommand(NewDeleteCommand(sl).Command)
	cobraCmd.AddCommand(NewDetailsCommand(sl).Command)
	cobraCmd.AddCommand(NewPermissionsCommand(sl).Command)
	cobraCmd.AddCommand(NewEditCommand(sl).Command)
	cobraCmd.AddCommand(NewEditPermissionCommand(sl).Command)
	cobraCmd.AddCommand(NewNotificationsCommand(sl).Command)
	cobraCmd.AddCommand(NewEditNotificationsCommand(sl).Command)
	cobraCmd.AddCommand(NewGrantAccessCommand(sl).Command)
	cobraCmd.AddCommand(NewRemoveAccessCommand(sl).Command)
	cobraCmd.AddCommand(NewDeviceAccessCommand(sl).Command)
	cobraCmd.AddCommand(NewVpnSubnetCommand(sl).Command)
	cobraCmd.AddCommand(NewVpnManualCommand(sl).Command)
	cobraCmd.AddCommand(NewVpnPasswordCommand(sl).Command)
	cobraCmd.AddCommand(NewVpnDisableCommand(sl).Command)
	cobraCmd.AddCommand(NewApikeyCommand(sl).Command)
	return cobraCmd
}

func UserNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "user",
		Description: T("Classic infrastructure Manage Users"),
	}
}

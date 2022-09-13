package virtual

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type MigrateCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Guest                int
	Host                 int
	All                  bool
}

func NewMigrateCommand(sl *metadata.SoftlayerCommand) (cmd *MigrateCommand) {
	thisCmd := &MigrateCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "migrate",
		Short: T("Manage VSIs that require migration"),
		Long: T(`${COMMAND_NAME} sl vs migrate [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs migrate --guest 1234567
   Manage VSIs that require migration. Can migrate Dedicated Instance from one dedicated host to another dedicated host as well.`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().IntVarP(&thisCmd.Guest, "guest", "g", 0, T("Guest ID to immediately migrate."))
	cobraCmd.Flags().IntVarP(&thisCmd.Host, "host", "H", 0, T("Dedicated Host ID to migrate to. Only works on guests that are already on a dedicated host."))
	cobraCmd.Flags().BoolVarP(&thisCmd.All, "all", "a", false, T("Migrate ALL guests that require migration immediately."))
	return thisCmd
}

func (cmd *MigrateCommand) Run(args []string) error {
	filters := filter.New()
	vsList := []datatypes.Virtual_Guest{}
	objMask := "mask[id, hostname, domain, datacenter, pendingMigrationFlag, powerState, primaryIpAddress,primaryBackendIpAddress, dedicatedHost]"

	outputFormat := cmd.GetOutputFlag()

	// No options, just show what is going to be migrated
	if cmd.Guest == 0 && !cmd.All && cmd.Host == 0 {
		vsPendignMigateList := getMigrationServerList(objMask, nil, cmd)
		for _, pendingMigrationsVs := range vsPendignMigateList {
			if *pendingMigrationsVs.PendingMigrationFlag {
				vsList = append(vsList, pendingMigrationsVs)
			}
		}

		dedicatedFilter := append(filters, utils.QueryFilter("not null", "virtualGuests.dedicatedHost.id"))
		dedicatedMigrateList := getMigrationServerList(objMask, dedicatedFilter, cmd)

		var migrationList []interface{}
		migrationList = append(migrationList, vsList)
		migrationList = append(migrationList, dedicatedMigrateList)
		if outputFormat == "JSON" {
			return utils.PrintPrettyJSONList(cmd.UI, migrationList)
		}

		showsServerPendingMigration(vsList, cmd, "vs")
		showsServerPendingMigration(dedicatedMigrateList, cmd, "dedicated")
	} else {
		if cmd.All {
			guestMigration := getMigrationServerList(objMask, nil, cmd)
			if len(guestMigration) == 0 {
				return slErrors.New(T("No guests require migration at this time.\n"))
			}
			for _, guest := range guestMigration {
				if *guest.PendingMigrationFlag {
					result, err := cmd.VirtualServerManager.MigrateInstance(*guest.Id)
					if err != nil {
						return slErrors.NewAPIError(T("Failed to migrate the virtual server instance.\n"), err.Error(), 2)
					}

					cmd.UI.Ok()
					cmd.UI.Print(T("The virtual server is migrating: {{.VsId}}.", map[string]interface{}{"VsId": result.Id}))
				}
			}
		}
		if cmd.Host != 0 {
			if cmd.Guest == 0 {
				return slErrors.New(T("Please add the '--guest' id too.\n"))
			}
			err := cmd.VirtualServerManager.MigrateDedicatedHost(cmd.Guest, cmd.Host)
			if err != nil {
				return slErrors.NewAPIError(T("Failed to migrate the dedicated host instance.\n"), err.Error(), 2)
			}

			cmd.UI.Print(T("The dedicated host is migrating: {{.HostId}}.", map[string]interface{}{"HostId": cmd.Host}))
		}
		if cmd.Guest != 0 {
			result, err := cmd.VirtualServerManager.MigrateInstance(cmd.Guest)
			if err != nil {
				return slErrors.NewAPIError(T("Failed to migrate the virtual server instance.\n"), err.Error(), 2)
			}
			if outputFormat == "JSON" {
				return utils.PrintPrettyJSON(cmd.UI, result)
			}
			cmd.UI.Ok()
			cmd.UI.Print(T("The virtual server is migrating."))
			table := cmd.UI.Table([]string{T("id"), T("CreateDate")})
			table.Add(utils.FormatIntPointer(result.Id), utils.FormatSLTimePointer(result.CreateDate))
			table.Print()
		}
	}

	return nil
}

func getMigrationServerList(mask string, filter filter.Filters, cmd *MigrateCommand) []datatypes.Virtual_Guest {
	migrationServerList, err := cmd.VirtualServerManager.GetInstances(mask, filter)
	if err != nil {
		slErrors.NewAPIError(T("Failed to retrieve the virtual server instances.\n"), err.Error(), 2)
	}
	return migrationServerList
}

func showsServerPendingMigration(vsList []datatypes.Virtual_Guest, cmd *MigrateCommand, typeServer string) {
	if typeServer == "vs" {
		table := cmd.UI.Table([]string{T("id"), T("hostname"), T("domain"), T("datacenter"), T("pendingMigrationFlag")})
		cmd.UI.Print("Virtual Server Pending Migration")
		for _, vm := range vsList {
			table.Add(utils.FormatIntPointer(vm.Id), utils.FormatStringPointer(vm.Hostname),
				utils.FormatStringPointer(vm.Domain), utils.FormatStringPointer(vm.Datacenter.Name),
				utils.FormatBoolPointer(vm.PendingMigrationFlag))
		}
		table.Print()
		fmt.Println()
	} else {
		table := cmd.UI.Table([]string{T("id"), T("hostname"), T("domain"), T("datacenter"), T("PendingMigrationFlag"),
			T("HostName"), T("HostId")})
		cmd.UI.Print("Dedicated Hosts")
		for _, vm := range vsList {
			table.Add(utils.FormatIntPointer(vm.Id), utils.FormatStringPointer(vm.Hostname),
				utils.FormatStringPointer(vm.Domain), utils.FormatStringPointer(vm.Datacenter.Name),
				utils.FormatBoolPointer(vm.PendingMigrationFlag), utils.FormatStringPointer(vm.DedicatedHost.Name),
				utils.FormatIntPointer(vm.DedicatedHost.Id))
		}
		table.Print()
	}
}

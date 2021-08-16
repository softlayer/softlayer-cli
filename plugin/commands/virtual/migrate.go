package virtual

import (
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type MigrateCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewMigrageCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *MigrateCommand) {
	return &MigrateCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *MigrateCommand) Run(c *cli.Context) error {
	filters := filter.New()
	vsList := []datatypes.Virtual_Guest{}
	objMask := "mask[id, hostname, domain, datacenter, pendingMigrationFlag, powerState, primaryIpAddress,primaryBackendIpAddress, dedicatedHost]"

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	if !c.IsSet("g") && !c.IsSet("a") && !c.IsSet("host") {
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
		if c.IsSet("all") {
			guestMigration := getMigrationServerList(objMask, nil, cmd)
			if len(guestMigration) == 0 {
				return cli.NewExitError(T("No guests require migration at this time.\n"), 2)
			}
			for _, guest := range guestMigration {
				if *guest.PendingMigrationFlag {
					result, err := cmd.VirtualServerManager.MigrateInstance(*guest.Id)
					if err != nil {
						return cli.NewExitError(T("Failed to migrate the virtual server instance.\n")+err.Error(), 2)
					}

					cmd.UI.Ok()
					cmd.UI.Print(T("The virtual server is migrating: {{.VsId}}.", map[string]interface{}{"VsId": result.Id}))
				}
			}
		}
		if c.IsSet("host") {
			if !c.IsSet("guest") {
				return cli.NewExitError(T("Please add the '--guest' id too.\n"), 2)
			}
			err := cmd.VirtualServerManager.MigrateDedicatedHost(c.Int("guest"), c.Int("host"))
			if err != nil {
				return cli.NewExitError(T("Failed to migrate the dedicated host instance.\n")+err.Error(), 2)
			}

			cmd.UI.Print(T("The dedicated host is migrating: {{.HostId}}.", map[string]interface{}{"HostId": c.Int("host")}))
		}
		if c.IsSet("guest") {
			result, err := cmd.VirtualServerManager.MigrateInstance(c.Int("guest"))
			if err != nil {
				return cli.NewExitError(T("Failed to migrate the virtual server instance.\n")+err.Error(), 2)
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
		cli.NewExitError(T("Failed to retrieve the virtual server instances.\n")+err.Error(), 2)
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

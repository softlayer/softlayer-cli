package user

import (
	"bytes"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DeviceAccessCommand struct {
	UI          terminal.UI
	UserManager managers.UserManager
}

func NewDeviceAccessCommand(ui terminal.UI, userManager managers.UserManager) (cmd *DeviceAccessCommand) {
	return &DeviceAccessCommand{
		UI:          ui,
		UserManager: userManager,
	}
}

func (cmd *DeviceAccessCommand) Run(c *cli.Context) error {

	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	userId := c.Args()[0]
	id, err := strconv.Atoi(userId)
	if err != nil {
		return errors.NewInvalidUsageError(T("User ID should be a number."))
	}

	userpermissions, err := cmd.UserManager.GetUserAllowDevicesPermissions(id)
	if err != nil {
		return cli.NewExitError(T("Failed to get user permissions.\n")+err.Error(), 2)
	}

	dedicatedHosts, err := cmd.UserManager.GetDedicatedHosts(id)
	if err != nil {
		return cli.NewExitError(T("Failed to get dedicated hosts.\n")+err.Error(), 2)
	}

	hardwares, err := cmd.UserManager.GetHardware(id)
	if err != nil {
		return cli.NewExitError(T("Failed to get bare metal servers.\n")+err.Error(), 2)
	}

	virtualGuests, err := cmd.UserManager.GetVirtualGuests(id)
	if err != nil {
		return cli.NewExitError(T("Failed to get virtual servers.\n")+err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("User ID"), utils.FormatIntPointer(&id))
	if len(userpermissions) != 0 {
		buf := new(bytes.Buffer)
		permissionsTable := terminal.NewTable(buf, []string{T("Key Name"), T("Name")})
		for _, permission := range userpermissions {
			permissionsTable.Add(
				utils.FormatStringPointer(permission.KeyName),
				utils.FormatStringPointer(permission.Name),
			)
		}
		permissionsTable.Print()
		table.Add("Permissions", buf.String())
	} else {
		table.Add(T("Permissions"), "User does not have permissions about devices")
	}

	if len(dedicatedHosts) != 0 || len(hardwares) != 0 || len(virtualGuests) != 0 {
		buf := new(bytes.Buffer)
		devicesTable := terminal.NewTable(buf, []string{T("Id"), T("Device Name"), T("Device Type"), T("Public Ip"), T("Private Ip"), T("Notes")})
		if len(dedicatedHosts) != 0 {
			for _, device := range dedicatedHosts {
				notes := ""
				if device.Notes != nil {
					notes = *device.Notes
				}
				devicesTable.Add(
					utils.FormatIntPointer(device.Id),
					utils.FormatStringPointer(device.Name),
					"Dedicated host",
					"",
					"",
					notes,
				)
			}
		}
		if len(hardwares) != 0 {
			for _, device := range hardwares {
				notes := ""
				publicIp := ""
				privateIp := ""
				if device.Notes != nil {
					notes = *device.Notes
				}
				if device.PrimaryIpAddress != nil {
					publicIp = *device.PrimaryIpAddress
				}
				if device.PrimaryBackendIpAddress != nil {
					privateIp = *device.PrimaryBackendIpAddress
				}
				devicesTable.Add(
					utils.FormatIntPointer(device.Id),
					utils.FormatStringPointer(device.FullyQualifiedDomainName),
					"Bare metal server",
					publicIp,
					privateIp,
					notes,
				)
			}
		}
		if len(virtualGuests) != 0 {
			for _, device := range virtualGuests {
				notes := ""
				publicIp := ""
				privateIp := ""
				if device.Notes != nil {
					notes = *device.Notes
				}
				if device.PrimaryIpAddress != nil {
					publicIp = *device.PrimaryIpAddress
				}
				if device.PrimaryBackendIpAddress != nil {
					privateIp = *device.PrimaryBackendIpAddress
				}
				devicesTable.Add(
					utils.FormatIntPointer(device.Id),
					utils.FormatStringPointer(device.FullyQualifiedDomainName),
					"Virtual server",
					publicIp,
					privateIp,
					notes,
				)
			}
		}
		devicesTable.Print()
		table.Add("Devices", buf.String())
	} else {
		table.Add(T("Devices"), "User does not have devices")
	}

	table.Print()
	return nil
}

func UserDeviceAccessMetaData() cli.Command {
	return cli.Command{
		Category:    "user",
		Name:        "device-access",
		Description: T("List all devices the user has access and device access permissions."),
		Usage:       "${COMMAND_NAME} sl user device-access IDENTIFIER",
		Flags:       []cli.Flag{},
	}
}

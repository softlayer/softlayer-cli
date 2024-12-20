package user

import (
	"bytes"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DeviceAccessCommand struct {
	*metadata.SoftlayerCommand
	UserManager managers.UserManager
	Command     *cobra.Command
}

func NewDeviceAccessCommand(sl *metadata.SoftlayerCommand) (cmd *DeviceAccessCommand) {
	thisCmd := &DeviceAccessCommand{
		SoftlayerCommand: sl,
		UserManager:      managers.NewUserManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "device-access " + T("IDENTIFIER"),
		Short: T("List all devices the user has access and device access permissions."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DeviceAccessCommand) Run(args []string) error {
	userId := args[0]
	id, err := strconv.Atoi(userId)
	if err != nil {
		return errors.NewInvalidUsageError(T("User ID should be a number."))
	}

	outputFormat := cmd.GetOutputFlag()

	userpermissions, err := cmd.UserManager.GetUserAllowDevicesPermissions(id)
	if err != nil {
		return errors.NewAPIError(T("Failed to get user permissions.\n"), err.Error(), 2)
	}

	dedicatedHosts, err := cmd.UserManager.GetDedicatedHosts(id, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get dedicated hosts.\n"), err.Error(), 2)
	}

	hardwares, err := cmd.UserManager.GetHardware(id, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get bare metal servers.\n"), err.Error(), 2)
	}

	virtualGuests, err := cmd.UserManager.GetVirtualGuests(id, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get virtual servers.\n"), err.Error(), 2)
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
		table.Add(T("Permissions"), buf.String())
	} else {
		table.Add(T("Permissions"), "-")
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
		table.Add(T("Devices"), buf.String())
	} else {
		table.Add(T("Devices"), "-")
	}

	if outputFormat == "JSON" {
		table.PrintJson()
	} else {
		table.Print()
	}
	return nil
}

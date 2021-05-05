package block

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type AccessRevokeCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
	NetworkManager managers.NetworkManager
}

func NewAccessRevokeCommand(ui terminal.UI, storageManager managers.StorageManager, networkManager managers.NetworkManager) (cmd *AccessRevokeCommand) {
	return &AccessRevokeCommand{
		UI:             ui,
		StorageManager: storageManager,
		NetworkManager: networkManager,
	}
}

func (cmd *AccessRevokeCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	if !c.IsSet("d") && !c.IsSet("v") && !c.IsSet("i") && !c.IsSet("p") && !c.IsSet("hardware-id") && !c.IsSet("virtual-id") && !c.IsSet("ip-address-id") && !c.IsSet("ip-address") {
		return errors.NewInvalidUsageError(T("One of -d | --hardware-id, -v | --virtual-id, -i | --ip-address-id and -p | --ip-address must be specified."))
	}

	IPIds := c.IntSlice("ip-address-id")
	IPs := c.StringSlice("ip-address")
	if len(IPs) > 0 {
		for _, ip := range IPs {
			ipRecord, err := cmd.NetworkManager.IPLookup(ip)
			if err != nil {
				return cli.NewExitError(T("IP address {{.IP}} is not found on your account.Please confirm IP and try again.\n",
					map[string]interface{}{"IP": ip})+err.Error(), 2)
			}
			if ipRecord.Id != nil {
				IPIds = append(IPIds, *ipRecord.Id)
			}

		}
	}
	_, err = cmd.StorageManager.DeauthorizeHostToVolume(volumeID, c.IntSlice("hardware-id"), c.IntSlice("virtual-id"), IPIds, nil)
	if err != nil {
		return cli.NewExitError(T("Failed to revoke access to volume {{.VolumeID}}.\n", map[string]interface{}{"VolumeID": volumeID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	for _, vsID := range c.IntSlice("virtual-id") {
		cmd.UI.Print(T("Access to {{.VolumeId}} was revoked for virtual server {{.VsID}}.", map[string]interface{}{"VolumeId": volumeID, "VsID": vsID}))
	}
	for _, hwID := range c.IntSlice("hardware-id") {
		cmd.UI.Print(T("Access to {{.VolumeId}} was revoked for hardware server {{.HwID}}.", map[string]interface{}{"VolumeId": volumeID, "HwID": hwID}))
	}
	for _, ip := range IPIds {
		cmd.UI.Print(T("Access to {{.VolumeId}} was revoked for IP address {{.IP}}.", map[string]interface{}{"VolumeId": volumeID, "IP": ip}))
	}
	return nil
}

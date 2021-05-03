package file

import (
	"strconv"

	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
)

type AccessAuthorizeCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
	NetworkManager managers.NetworkManager
}

func NewAccessAuthorizeCommand(ui terminal.UI, storageManager managers.StorageManager, networkManager managers.NetworkManager) (cmd *AccessAuthorizeCommand) {
	return &AccessAuthorizeCommand{
		UI:             ui,
		StorageManager: storageManager,
		NetworkManager: networkManager,
	}
}

func (cmd *AccessAuthorizeCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
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

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	resp, err := cmd.StorageManager.AuthorizeHostToVolume(volumeID, c.IntSlice("hardware-id"), c.IntSlice("virtual-id"), IPIds, c.IntSlice("s"))
	if err != nil {
		return cli.NewExitError(T("Failed to authorize host to volume {{.VolumeID}}.\n", map[string]interface{}{"VolumeID": volumeID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	for _, vsID := range c.IntSlice("virtual-id") {
		cmd.UI.Print(T("The virtual server {{.VsID}} was authorized to access {{.VolumeId}}.", map[string]interface{}{"VolumeId": volumeID, "VsID": vsID}))
	}
	for _, hwID := range c.IntSlice("hardware-id") {
		cmd.UI.Print(T("The hardware server {{.HwID}} was authorized to access {{.VolumeId}}.", map[string]interface{}{"VolumeId": volumeID, "HwID": hwID}))
	}
	for _, ip := range IPIds {
		cmd.UI.Print(T("The IP address {{.IP}} was authorized to access {{.VolumeId}}.", map[string]interface{}{"VolumeId": volumeID, "IP": ip}))
	}
	for _, sn := range c.IntSlice("s") {
		cmd.UI.Print(T("The subnet {{.Subnet}} was authorized to access {{.VolumeId}}.", map[string]interface{}{"VolumeId": volumeID, "Subnet": sn}))
	}
	return nil
}

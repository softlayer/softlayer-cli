package block

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
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

func BlockAccessAuthorizeMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "access-authorize",
		Description: T("Authorize hosts to access a given volume"),
		Usage: T(`${COMMAND_NAME} sl block access-authorize VOLUME_ID [OPTIONS]
		
EXAMPLE:
   ${COMMAND_NAME} sl block access-authorize 12345678 --virtual-id 87654321
   This command authorizes virtual server with ID 87654321 to access volume with ID 12345678.`),
		Flags: []cli.Flag{
			cli.IntSliceFlag{
				Name:  "d,hardware-id",
				Usage: T("The ID of one hardware server to authorize"),
			},
			cli.IntSliceFlag{
				Name:  "v,virtual-id",
				Usage: T("The ID of one virtual server to authorize"),
			},
			cli.IntSliceFlag{
				Name:  "i,ip-address-id",
				Usage: T("The ID of one IP address to authorize"),
			},
			cli.StringSliceFlag{
				Name:  "p,ip-address",
				Usage: T("An IP address to authorize"),
			},
			metadata.OutputFlag(),
		},
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

	resp, err := cmd.StorageManager.AuthorizeHostToVolume(volumeID, c.IntSlice("hardware-id"), c.IntSlice("virtual-id"), IPIds, nil)
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
	return nil
}

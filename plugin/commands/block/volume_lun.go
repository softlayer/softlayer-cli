package block

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type VolumeLunCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewVolumeLunCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *VolumeLunCommand) {
	return &VolumeLunCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func BlockVolumeLunMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "volume-set-lun-id",
		Description: T("Set the LUN ID on an existing block storage volume"),
		Usage: T(`${COMMAND_NAME} sl block volume-set-lun-id VOLUME_ID LUN_ID

	The LUN ID only takes effect during the Host Authorization process. It is
	recommended (but not necessary) to de-authorize all hosts before using
	this method.
	VOLUME_ID - the volume ID on which to set the LUN ID
	LUN_ID - recommended range is an integer between 0 and 255. Advanced users
	can use an integer between 0 and 4095`),
	}
}

func (cmd *VolumeLunCommand) Run(c *cli.Context) error {
	if c.NArg() != 2 {
		return errors.NewInvalidUsageError(T("This command requires two arguments."))
	}
	volumeId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	lunId, err := strconv.Atoi(c.Args()[1])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("LUN ID")
	}
	prop, err := cmd.StorageManager.SetLunId(volumeId, lunId)
	if err != nil {
		return cli.NewExitError(T("Failed to set LUN ID for volume {{.VolumeID}}.\n", map[string]interface{}{"VolumeID": volumeId})+err.Error(), 2)
	}
	if prop.Value != nil {
		newLunId, err := strconv.Atoi(*prop.Value)
		if err == nil && newLunId == lunId {
			cmd.UI.Ok()
			cmd.UI.Print(T("Block volume {{.VolumeId}} is reporting LUN ID {{.LunID}}.",
				map[string]interface{}{"VolumeId": volumeId, "LunID": lunId}))
			return nil
		}
	}
	cmd.UI.Failed(T("Failed to confirm the new LUN ID on volume {{.VolumeId}}.", map[string]interface{}{"VolumeId": volumeId}))
	return nil
}

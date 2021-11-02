package file

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type VolumeRefreshCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewVolumeRefreshCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *VolumeRefreshCommand) {
	return &VolumeRefreshCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func (cmd *VolumeRefreshCommand) Run(c *cli.Context) error {

	if c.NArg() != 2 {
		return errors.NewInvalidUsageError(T("This command requires two arguments."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	snapshotId, err := strconv.Atoi(c.Args()[1])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Snapshot ID")
	}

	err = cmd.StorageManager.VolumeRefresh(volumeID, snapshotId)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	cmd.UI.Ok()
	return nil
}

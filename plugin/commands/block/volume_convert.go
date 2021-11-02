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

type VolumeConvertCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewVolumeConvertCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *VolumeConvertCommand) {
	return &VolumeConvertCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func (cmd *VolumeConvertCommand) Run(c *cli.Context) error {

	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	err = cmd.StorageManager.VolumeConvert(volumeID)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	cmd.UI.Ok()
	return nil
}

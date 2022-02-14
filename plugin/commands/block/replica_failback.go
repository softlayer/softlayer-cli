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

type ReplicaFailbackCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewReplicaFailbackCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *ReplicaFailbackCommand) {
	return &ReplicaFailbackCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func BlockReplicaFailbackMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "replica-failback",
		Description: T("Failback a block volume from replica"),
		Usage: T(`${COMMAND_NAME} sl block replica-failback VOLUME_ID
		
EXAMPLE:
   ${COMMAND_NAME} sl block replica-failback 12345678
   This command performs failback operation for volume with ID 12345678.`),
	}
}

func (cmd *ReplicaFailbackCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	err = cmd.StorageManager.FailBackFromReplicant(volumeID)
	if err != nil {
		return cli.NewExitError(T("Failback operation could not be initiated for volume {{.VolumeID}}.\n", map[string]interface{}{"VolumeID": volumeID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Failback of volume {{.VolumeID}} is now in progress.", map[string]interface{}{"VolumeID": volumeID}))
	return nil
}

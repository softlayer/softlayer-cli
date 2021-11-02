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

type DisasterRecoveryFailoverCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewDisasterRecoveryFailoverCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *DisasterRecoveryFailoverCommand) {
	return &DisasterRecoveryFailoverCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func (cmd *DisasterRecoveryFailoverCommand) Run(c *cli.Context) error {
	if c.NArg() != 2 {
		return errors.NewInvalidUsageError(T("This command requires two arguments."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	replicaID, err := strconv.Atoi(c.Args()[1])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Replica ID")
	}
	err = cmd.StorageManager.DisasterRecoveryFailover(volumeID, replicaID)
	if err != nil {
		return cli.NewExitError(T("Failover operation could not be initiated for volume {{.VolumeID}}.\n", map[string]interface{}{"VolumeID": volumeID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Failover of volume {{.VolumeID}} to replica {{.ReplicaID}} is now in progress.",
		map[string]interface{}{"VolumeID": volumeID, "ReplicaID": replicaID}))
	return nil
}

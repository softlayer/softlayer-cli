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

func BlockDisasterRecoveryFailoverMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "disaster-recovery-failover",
		Description: T("Failover an inaccessible volume to its available replicant volume."),
		Usage: T(`${COMMAND_NAME} sl block disaster-recovery-failover VOLUME_ID REPLICA_ID

If a volume (with replication) becomes inaccessible due to a disaster event, this method can be used to immediately
failover to an available replica in another location. This method does not allow for fail back via the API.
To fail back to the original volume after using this method, open a support ticket.
To test failover, use '${COMMAND_NAME} sl block replica-failover' instead.

EXAMPLE:
	${COMMAND_NAME} sl block disaster-recovery-failover 12345678 87654321
	This command performs failover operation for volume with ID 12345678 to replica volume with ID 87654321.`),
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

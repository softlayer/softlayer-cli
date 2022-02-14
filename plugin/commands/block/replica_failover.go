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

type ReplicaFailoverCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewReplicaFailoverCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *ReplicaFailoverCommand) {
	return &ReplicaFailoverCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func BlockReplicaFailOverMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "replica-failover",
		Description: T("Failover a block volume to the given replica volume"),
		Usage: T(`${COMMAND_NAME} sl block replica-failover VOLUME_ID REPLICA_ID
		
EXAMPLE:
   ${COMMAND_NAME} sl block replica-failover 12345678 87654321
   This command performs failover operation for volume with ID 12345678 to replica volume with ID 87654321.`),
	}
}

func (cmd *ReplicaFailoverCommand) Run(c *cli.Context) error {
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
	err = cmd.StorageManager.FailOverToReplicant(volumeID, replicaID)
	if err != nil {
		return cli.NewExitError(T("Failover operation could not be initiated for volume {{.VolumeID}}.\n", map[string]interface{}{"VolumeID": volumeID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Failover of volume {{.VolumeID}} to replica {{.ReplicaID}} is now in progress.",
		map[string]interface{}{"VolumeID": volumeID, "ReplicaID": replicaID}))
	return nil
}

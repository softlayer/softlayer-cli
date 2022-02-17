package file

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SnapshotCreateCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewSnapshotCreateCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *SnapshotCreateCommand) {
	return &SnapshotCreateCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func FileSnapshotCreateMetaData() cli.Command {
	return cli.Command{
		Category:    "file",
		Name:        "snapshot-create",
		Description: T("Create a snapshot on a given volume"),
		Usage: T(`${COMMAND_NAME} sl file snapshot-create VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file snapshot-create 12345678 --note snapshotforibmcloud
   This command creates a snapshot for volume with ID 12345678 and with addition note as snapshotforibmcloud.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,note",
				Usage: T("Notes to set on the new snapshot"),
			},
			metadata.OutputFlag(),
		},
	}
}

func (cmd *SnapshotCreateCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	snapshot, err := cmd.StorageManager.CreateSnapshot(volumeID, c.String("note"))
	if err != nil {
		return cli.NewExitError(T("Error occurred while creating snapshot for volume {{.VolumeID}}.Ensure volume is not failed over or in another state which prevents taking snapshots.\n",
			map[string]interface{}{"VolumeID": volumeID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, snapshot)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("New snapshot {{.SnapshotId}} was created.", map[string]interface{}{"SnapshotId": *snapshot.Id}))
	return nil
}

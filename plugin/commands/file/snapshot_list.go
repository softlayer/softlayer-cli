package file

import (
	"sort"
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

type SnapshotListCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewSnapshotListCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *SnapshotListCommand) {
	return &SnapshotListCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func FileSnapshotListMetaData() cli.Command {
	return cli.Command{
		Category:    "file",
		Name:        "snapshot-list",
		Description: T("List file storage snapshots"),
		Usage: T(`${COMMAND_NAME} sl file snapshot-list VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file snapshot-list 12345678 --sortby id 
   This command lists all snapshots of volume with ID 12345678 and sorts them by ID.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by. Options are: id,name,created,size_bytes"),
			},
			metadata.OutputFlag(),
		},
	}
}

func (cmd *SnapshotListCommand) Run(c *cli.Context) error {
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

	snapshots, err := cmd.StorageManager.GetVolumeSnapshotList(volumeID)
	if err != nil {
		return cli.NewExitError(T("Failed to get snapshot list on your account.\n")+err.Error(), 2)
	}
	sortby := c.String("sortby")
	if sortby == "id" || sortby == "ID" {
		sort.Sort(utils.SnapshotsById(snapshots))
	} else if sortby == "name" {
		sort.Sort(utils.SnapshotsByName(snapshots))
	} else if sortby == "created" {
		sort.Sort(utils.SnapshotsByCreated(snapshots))
	} else if sortby == "size_bytes" {
		sort.Sort(utils.SnapshotsBySize(snapshots))
	} else if sortby == "" {
		//do nothing
	} else {
		return errors.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.", map[string]interface{}{"Column": sortby}))
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, snapshots)
	}

	table := cmd.UI.Table([]string{T("id"), T("user_name"), T("created"), T("size_bytes"), T("notes")})
	for _, sp := range snapshots {
		table.Add(utils.FormatIntPointer(sp.Id),
			utils.FormatStringPointer(sp.Username),
			utils.FormatStringPointer(sp.SnapshotCreationTimestamp),
			utils.FormatStringPointer(sp.SnapshotSizeBytes),
			utils.FormatStringPointer(sp.Notes))
	}
	table.Print()
	return nil
}

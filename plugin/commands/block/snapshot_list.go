package block

import (
	"github.com/spf13/cobra"
	"sort"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SnapshotListCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	SortBy         string
}

func NewSnapshotListCommand(sl *metadata.SoftlayerStorageCommand) *SnapshotListCommand {
	thisCmd := &SnapshotListCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "snapshot-list " + T("IDENTIFIER"),
		Short: T("List {{.storageType}} storage snapshots", sl.StorageI18n),
		Long: T(`
EXAMPLE:
   ${COMMAND_NAME} sl {{.storageType}} snapshot-list 12345678 --sortby id 
   This command lists all snapshots of volume with ID 12345678 and sorts them by ID.`, sl.StorageI18n),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.SortBy, "sortby", "", T("Column to sort by. Options are: id,name,created,size_bytes"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SnapshotListCommand) Run(args []string) error {
	volumeID, err := cmd.StorageManager.GetVolumeId(args[0], cmd.StorageType)
	if err != nil {
		return err
	}
	outputFormat := cmd.GetOutputFlag()

	snapshots, err := cmd.StorageManager.GetVolumeSnapshotList(volumeID)
	if err != nil {
		return slErr.NewAPIError(T("Failed to get snapshot list on your account.\n"), err.Error(), 2)
	}
	sortby := cmd.SortBy
	if sortby == "" {
		// do nothing
	} else if sortby == "id" || sortby == "ID" {
		sort.Sort(utils.SnapshotsById(snapshots))
	} else if sortby == "name" {
		sort.Sort(utils.SnapshotsByName(snapshots))
	} else if sortby == "created" {
		sort.Sort(utils.SnapshotsByCreated(snapshots))
	} else if sortby == "size_bytes" {
		sort.Sort(utils.SnapshotsBySize(snapshots))
	} else {
		return slErr.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.", map[string]interface{}{"Column": sortby}))
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

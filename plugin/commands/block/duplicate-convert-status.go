package block

import (
	"strconv"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DuplicateConvertStatusCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
}

func NewDuplicateConvertStatusCommand(sl *metadata.SoftlayerStorageCommand) *DuplicateConvertStatusCommand {
	thisCmd := &DuplicateConvertStatusCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "duplicate-convert-status " + T("IDENTIFIER"),
		Short: T("Get status for split or move completed percentage of a given block storage duplicate volume."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DuplicateConvertStatusCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	outputFormat := cmd.GetOutputFlag()

	duplicateConversionStatus, err := cmd.StorageManager.GetDuplicateConversionStatus(volumeID, "")
	if err != nil {
		return slErr.NewAPIError(T("Failed to get duplicate conversion status of volume {{.VolumeID}}.\n",
			map[string]interface{}{"VolumeID": volumeID}), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Username"), T("Active Conversion Start Timestamp"), T("Completed Percentage")})
	table.Add(
		utils.FormatStringPointer(duplicateConversionStatus.VolumeUsername),
		utils.FormatStringPointer(duplicateConversionStatus.ActiveConversionStartTime),
		utils.FormatIntPointer(duplicateConversionStatus.DeDuplicateConversionPercentage),
	)

	utils.PrintTable(cmd.UI, table, outputFormat)

	return nil
}

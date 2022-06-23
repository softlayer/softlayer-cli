package block

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

type DuplicateConvertStatusCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewDuplicateConvertStatusCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *DuplicateConvertStatusCommand) {
	return &DuplicateConvertStatusCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func (cmd *DuplicateConvertStatusCommand) Run(c *cli.Context) error {
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

	duplicateConversionStatus, err := cmd.StorageManager.GetDuplicateConversionStatus(volumeID, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get duplicate conversion status of volume {{.VolumeID}}.\n",
			map[string]interface{}{"VolumeID": volumeID})+err.Error(), 2)
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

func BlockDuplicateConvertStatusMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "duplicate-convert-status",
		Description: T("Get status for split or move completed percentage of a given block storage duplicate volume."),
		Usage: T(`${COMMAND_NAME} sl block duplicate-convert-status [OPTIONS] VOLUME_ID

EXAMPLE:
   ${COMMAND_NAME} sl block duplicate-convert-status 12345678`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

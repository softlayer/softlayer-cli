package file

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type VolumeLimitCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewVolumeLimitCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *VolumeLimitCommand) {
	return &VolumeLimitCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func FileVolumeLimitsMetaData() cli.Command {
	return cli.Command{
		Category:    "file",
		Name:        "volume-limits",
		Description: T("Lists the storage limits per datacenter for this account."),
		Usage: T(`${COMMAND_NAME} sl file volume-limits [OPTIONS]

EXAMPLE:
	${COMMAND_NAME} sl file volume-limits
	This command lists the storage limits per datacenter for this account.`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

func (cmd *VolumeLimitCommand) Run(c *cli.Context) error {

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	volumeLimits, err := cmd.StorageManager.GetVolumeCountLimits()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, volumeLimits)
	}

	table := cmd.UI.Table([]string{"Datacenter", "MaximumAvailableCount", "ProvisionedCount"})
	for _, row := range volumeLimits {
		table.Add(
			utils.FormatStringPointer(row.DatacenterName),
			utils.FormatIntPointer(row.MaximumAvailableCount),
			utils.FormatIntPointer(row.ProvisionedCount))
	}
	table.Print()
	return nil
}

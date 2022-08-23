package block

import (
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type VolumeLimitCommand struct {
	*metadata.SoftlayerCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
}

func NewVolumeLimitCommand(sl *metadata.SoftlayerCommand) *VolumeLimitCommand {
	thisCmd := &VolumeLimitCommand{
		SoftlayerCommand: sl,
		StorageManager:   managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "volume-limits",
		Short: T("Lists the storage limits per datacenter for this account."),
		Long: T(`${COMMAND_NAME} sl block volume-limits [OPTIONS]

EXAMPLE:
	${COMMAND_NAME} sl block volume-limits
	This command lists the storage limits per datacenter for this account.`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VolumeLimitCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	volumeLimits, err := cmd.StorageManager.GetVolumeCountLimits()
	if err != nil {
		return err
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

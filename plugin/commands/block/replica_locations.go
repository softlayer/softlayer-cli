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

type ReplicaLocationsCommand struct {
	*metadata.SoftlayerCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
}

func NewReplicaLocationsCommand(sl *metadata.SoftlayerCommand) *ReplicaLocationsCommand {
	thisCmd := &ReplicaLocationsCommand{
		SoftlayerCommand: sl,
		StorageManager:   managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "replica-locations " + T("IDENTIFIER"),
		Short: T("List suitable replication datacenters for the given volume"),
		Long: T(`${COMMAND_NAME} sl block replica-locations VOLUME_ID [OPTIONS]
		
EXAMPLE:
   ${COMMAND_NAME} sl block replica-locations 12345678
   This command lists suitable replication data centers for block volume with ID 12345678.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ReplicaLocationsCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	outputFormat := cmd.GetOutputFlag()

	datacenters, err := cmd.StorageManager.GetReplicationLocations(volumeID)
	if err != nil {
		return slErr.NewAPIError(T("Failed to get datacenters for volume {{.VolumeID}}.\n", map[string]interface{}{"VolumeID": volumeID}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, datacenters)
	}

	if len(datacenters) == 0 {
		cmd.UI.Print(T("No data centers compatible for replication."))
	} else {
		table := cmd.UI.Table([]string{T("ID"), T("Short Name"), T("Long Name")})
		for _, d := range datacenters {
			table.Add(utils.FormatIntPointer(d.Id), utils.FormatStringPointer(d.Name), utils.FormatStringPointer(d.LongName))
		}
		table.Print()
	}
	return nil
}

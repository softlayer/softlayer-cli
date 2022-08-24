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

type ReplicaPartnersCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
}

func NewReplicaPartnersCommand(sl *metadata.SoftlayerStorageCommand) *ReplicaPartnersCommand {
	thisCmd := &ReplicaPartnersCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "replica-partners " + T("IDENTIFIER"),
		Short: T("List existing replicant volumes for a block volume"),
		Long: T(`${COMMAND_NAME} sl block replica-partners VOLUME_ID [OPTIONS]
		
EXAMPLE:
   ${COMMAND_NAME} sl block replica-partners 12345678
   This command lists existing replicant volumes for block volume with ID 12345678.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ReplicaPartnersCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	outputFormat := cmd.GetOutputFlag()

	partners, err := cmd.StorageManager.GetReplicationPartners(volumeID)
	subs := map[string]interface{}{"VolumeID": volumeID}
	if err != nil {
		return slErr.NewAPIError(T("Failed to get replication partners for volume {{.VolumeID}}.\n", subs), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, partners)
	}

	if len(partners) == 0 {
		cmd.UI.Print(T("There are no replication partners for volume {{.VolumeID}}.\n", subs))
	} else {
		table := cmd.UI.Table([]string{T("ID"), T("User name"), T("Account ID"), T("Capacity (GB)"), T("Hardware ID"), T("Guest ID"), T("Host ID")})
		for _, p := range partners {
			table.Add(
				utils.FormatIntPointer(p.Id),
				utils.FormatStringPointer(p.Username),
				utils.FormatIntPointer(p.AccountId),
				utils.FormatIntPointer(p.CapacityGb),
				utils.FormatIntPointer(p.HardwareId),
				utils.FormatIntPointer(p.GuestId),
				utils.FormatIntPointer(p.HostId),
			)
		}
		table.Print()
	}
	return nil
}

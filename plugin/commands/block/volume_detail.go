package block

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type VolumeDetailCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
}

func NewVolumeDetailCommand(sl *metadata.SoftlayerStorageCommand) *VolumeDetailCommand {
	thisCmd := &VolumeDetailCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "volume-detail " + T("IDENTIFIER"),
		Short: T("Display details for a specified volume"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VolumeDetailCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	outputFormat := cmd.GetOutputFlag()

	blockVolume, err := cmd.StorageManager.GetVolumeDetails("block", volumeID, "")
	if err != nil {
		return slErr.NewAPIError(T("Failed to get details of volume {{.VolumeID}}.\n",
			map[string]interface{}{"VolumeID": volumeID}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, blockVolume)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(blockVolume.Id))
	table.Add(T("User name"), utils.FormatStringPointer(blockVolume.Username))
	table.Add(T("Type"), strings.ToLower(utils.FormatStringPointer(blockVolume.StorageType.KeyName)))
	table.Add(T("Capacity (GB)"), utils.FormatIntPointer(blockVolume.CapacityGb))
	table.Add(T("LUN Id"), utils.FormatStringPointer(blockVolume.LunId))
	if blockVolume.Iops != nil {
		table.Add(T("IOPs"), utils.FormatStringPointer(blockVolume.Iops))
	}
	if blockVolume.StorageTierLevel != nil {
		tierLevel := utils.FormatStringPointer(blockVolume.StorageTierLevel)
		tierPerIops := managers.TIER_PER_IOPS[tierLevel]
		table.Add(T("Endurance Tier"), tierLevel)
		if tierPerIops == 0.25 {
			table.Add(T("Endurance Tier Per IOPS"), fmt.Sprintf("%.2f", tierPerIops))
		} else {
			table.Add(T("Endurance Tier Per IOPS"), fmt.Sprintf("%d", int(tierPerIops)))
		}
	}

	if blockVolume.ServiceResource != nil && blockVolume.ServiceResource.Datacenter != nil {
		table.Add(T("Datacenter"), utils.FormatStringPointer(blockVolume.ServiceResource.Datacenter.Name))
	}
	table.Add(T("Target IP"), utils.FormatStringPointer(blockVolume.ServiceResourceBackendIpAddress))

	if blockVolume.SnapshotCapacityGb != nil {
		table.Add(T("Snapshot Size (GB)"), utils.FormatStringPointer(blockVolume.SnapshotCapacityGb))
		if blockVolume.ParentVolume != nil {
			table.Add(T("Snapshot Used (Bytes)"), utils.FormatStringPointer(blockVolume.ParentVolume.SnapshotSizeBytes))
		}
	}

	table.Add(T("# of Active Transactions"), utils.FormatUIntPointer(blockVolume.ActiveTransactionCount))
	if blockVolume.ActiveTransactions != nil && len(blockVolume.ActiveTransactions) > 0 {
		for _, txn := range blockVolume.ActiveTransactions {
			if txn.TransactionStatus != nil {
				table.Add(T("Ongoing Transactions"), utils.FormatStringPointer(txn.TransactionStatus.FriendlyName))
			}
		}
	}

	table.Add(T("Replicant Count"), utils.FormatUIntPointer(blockVolume.ReplicationPartnerCount))
	// #nosec G115 -- Should never be > 2^32
	if blockVolume.ReplicationPartnerCount != nil && int(*blockVolume.ReplicationPartnerCount) > 0 {
		table.Add(T("Replication Status"), utils.FormatStringPointer(blockVolume.ReplicationStatus))
		buf := new(bytes.Buffer)
		repTable := terminal.NewTable(buf, []string{"", ""})
		for _, replicant := range blockVolume.ReplicationPartners {
			repTable.Add(T("Replicant ID"), utils.FormatIntPointer(replicant.Id))
			repTable.Add(T("Volume Name"), utils.FormatStringPointer(replicant.Username))
			repTable.Add(T("Target IP"), utils.FormatStringPointer(replicant.ServiceResourceBackendIpAddress))
			if replicant.ServiceResource != nil && replicant.ServiceResource.Datacenter != nil {
				repTable.Add(T("Datacenter"), utils.FormatStringPointer(replicant.ServiceResource.Datacenter.Name))
			}
			if replicant.ReplicationSchedule != nil && replicant.ReplicationSchedule.Type != nil {
				repTable.Add(T("Schedule"), utils.FormatStringPointer(replicant.ReplicationSchedule.Type.Keyname))
			}
		}
		repTable.Print()
		table.Add(T("Replicant Volumes"), buf.String())
	}

	if blockVolume.OriginalVolumeSize != nil {
		buf := new(bytes.Buffer)
		dupTable := terminal.NewTable(buf, []string{"", ""})
		dupTable.Add(T("Original Volume Name"), utils.FormatStringPointer(blockVolume.OriginalVolumeName))
		dupTable.Add(T("Original Volume Size"), utils.FormatStringPointer(blockVolume.OriginalVolumeSize))
		dupTable.Add(T("Original Snapshot Name"), utils.FormatStringPointer(blockVolume.OriginalSnapshotName))
		dupTable.Print()
		table.Add(T("Duplicate Volume Properties"), buf.String())
	}
	decodedValue, err := url.QueryUnescape(utils.FormatStringPointer(blockVolume.Notes))
	if err != nil {
		return slErr.NewAPIError(T("Failed to decoded the note.\n"), err.Error(), 2)
	}
	table.Add(T("Notes"), decodedValue)
	hasEncryption := T("False")
	if blockVolume.HasEncryptionAtRest != nil && *blockVolume.HasEncryptionAtRest == true {
		hasEncryption = T("True")
	}
	table.Add(T("Encrypted"), hasEncryption)
	table.Print()
	return nil
}

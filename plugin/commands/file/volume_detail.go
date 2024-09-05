package file

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

	fileVolume, err := cmd.StorageManager.GetVolumeDetails("file", volumeID, "")
	if err != nil {
		return slErr.NewAPIError(T("Failed to get details of volume {{.VolumeID}}.\n",
			map[string]interface{}{"VolumeID": volumeID}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, fileVolume)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(fileVolume.Id))
	table.Add(T("User name"), utils.FormatStringPointer(fileVolume.Username))
	table.Add(T("Type"), strings.ToLower(utils.FormatStringPointer(fileVolume.StorageType.KeyName)))
	table.Add(T("Capacity (GB)"), utils.FormatIntPointer(fileVolume.CapacityGb))
	table.Add(T("LUN Id"), utils.FormatStringPointer(fileVolume.LunId))
	if fileVolume.Iops != nil {
		table.Add(T("IOPs"), utils.FormatStringPointer(fileVolume.Iops))
	}
	if fileVolume.StorageTierLevel != nil {
		tierLevel := utils.FormatStringPointer(fileVolume.StorageTierLevel)
		tierPerIops := managers.TIER_PER_IOPS[tierLevel]
		table.Add(T("Endurance Tier"), tierLevel)
		if tierPerIops == 0.25 {
			table.Add(T("Endurance Tier Per IOPS"), fmt.Sprintf("%.2f", tierPerIops))
		} else {
			table.Add(T("Endurance Tier Per IOPS"), fmt.Sprintf("%d", int(tierPerIops)))
		}
	}

	if fileVolume.ServiceResource != nil && fileVolume.ServiceResource.Datacenter != nil {
		table.Add(T("Datacenter"), utils.FormatStringPointer(fileVolume.ServiceResource.Datacenter.Name))
	}
	table.Add(T("Target IP"), utils.FormatStringPointer(fileVolume.ServiceResourceBackendIpAddress))
	table.Add(T("Mount Address"), utils.FormatStringPointer(fileVolume.FileNetworkMountAddress))

	if fileVolume.SnapshotCapacityGb != nil {
		table.Add(T("Snapshot Size (GB)"), utils.FormatStringPointer(fileVolume.SnapshotCapacityGb))
		if fileVolume.ParentVolume != nil {
			table.Add(T("Snapshot Used (Bytes)"), utils.FormatStringPointer(fileVolume.ParentVolume.SnapshotSizeBytes))
		}
	}

	table.Add(T("# of Active Transactions"), utils.FormatUIntPointer(fileVolume.ActiveTransactionCount))
	if fileVolume.ActiveTransactions != nil && len(fileVolume.ActiveTransactions) > 0 {
		for _, txn := range fileVolume.ActiveTransactions {
			if txn.TransactionStatus != nil {
				table.Add(T("Ongoing Transactions"), utils.FormatStringPointer(txn.TransactionStatus.FriendlyName))
			}
		}
	}

	table.Add(T("Replicant Count"), utils.FormatUIntPointer(fileVolume.ReplicationPartnerCount))
	// #nosec G115 -- Should never be > 2^32
	if fileVolume.ReplicationPartnerCount != nil && int(*fileVolume.ReplicationPartnerCount) > 0 {
		table.Add(T("Replication Status"), utils.FormatStringPointer(fileVolume.ReplicationStatus))
		buf := new(bytes.Buffer)
		repTable := terminal.NewTable(buf, []string{"", ""})
		for _, replicant := range fileVolume.ReplicationPartners {
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

	if fileVolume.OriginalVolumeSize != nil {
		buf := new(bytes.Buffer)
		dupTable := terminal.NewTable(buf, []string{"", ""})
		dupTable.Add(T("Original Volume Name"), utils.FormatStringPointer(fileVolume.OriginalVolumeName))
		dupTable.Add(T("Original Volume Size"), utils.FormatStringPointer(fileVolume.OriginalVolumeSize))
		dupTable.Add(T("Original Snapshot Name"), utils.FormatStringPointer(fileVolume.OriginalSnapshotName))
		dupTable.Print()
		table.Add(T("Duplicate Volume Properties"), buf.String())
	}
	decodedValue, err := url.QueryUnescape(utils.FormatStringPointer(fileVolume.Notes))
	if err != nil {
		return slErr.NewAPIError(T("Failed to decoded the note.\n"), err.Error(), 2)
	}
	table.Add(T("Notes"), decodedValue)
	hasEncryption := T("False")
	if fileVolume.HasEncryptionAtRest != nil && *fileVolume.HasEncryptionAtRest == true {
		hasEncryption = T("True")
	}
	table.Add(T("Encrypted"), hasEncryption)
	table.Print()
	return nil
}

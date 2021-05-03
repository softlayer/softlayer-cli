package file

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type VolumeDetailCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewVolumeDetailCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *VolumeDetailCommand) {
	return &VolumeDetailCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func (cmd *VolumeDetailCommand) Run(c *cli.Context) error {
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

	fileVolume, err := cmd.StorageManager.GetVolumeDetails("file", volumeID, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get details of volume {{.VolumeID}}.\n",
			map[string]interface{}{"VolumeID": volumeID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, fileVolume)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(fileVolume.Id))
	table.Add(T("User name"), utils.FormatStringPointer(fileVolume.Username))
	table.Add(T("Type"), strings.ToLower(utils.FormatStringPointer(fileVolume.StorageType.KeyName)))
	table.Add(T("Capacity (GB)"), utils.FormatIntPointer(fileVolume.CapacityGb))
	//table.Add(T("LUN Id"), utils.FormatStringPointer(blockVolume.LunId))
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
	table.Print()
	return nil
}

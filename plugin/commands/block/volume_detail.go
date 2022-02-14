package block

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
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

func BlockVolumeDetailMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "volume-detail",
		Description: T("Display details for a specified volume"),
		Usage: T(`${COMMAND_NAME} sl block volume-detail VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl block volume-detail 12345678 
   This command shows details of volume with ID 12345678.`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
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

	blockVolume, err := cmd.StorageManager.GetVolumeDetails("block", volumeID, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get details of volume {{.VolumeID}}.\n",
			map[string]interface{}{"VolumeID": volumeID})+err.Error(), 2)
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
	table.Print()
	return nil
}

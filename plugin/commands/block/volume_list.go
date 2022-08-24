package block

import (
	"net/url"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type VolumeListCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	Username       string
	Datacenter     string
	StorageType    string
	Notes          string
	Order          int
	SortBy         string
	UserColumns    []string
}

func NewVolumeListCommand(sl *metadata.SoftlayerStorageCommand) *VolumeListCommand {
	thisCmd := &VolumeListCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "volume-list",
		Short: T("List block storage"),
		Long: T(`${COMMAND_NAME} sl {{.storageType}} volume-list [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl {{.storageType}} volume-list -d dal09 -t endurance --sortby capacity_gb
   This command lists all endurance volumes on current account that are located at dal09, and sorts them by capacity.`, sl.StorageI18n),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.Username, "username", "u", "", T("Filter by volume username"))
	cobraCmd.Flags().StringVarP(&thisCmd.Datacenter, "datacenter", "d", "", T("Filter by datacenter shortname"))
	cobraCmd.Flags().StringVarP(&thisCmd.StorageType, "storage-type", "t", "", T("Filter by type of storage volume, options are: performance,endurance"))
	cobraCmd.Flags().StringVarP(&thisCmd.Notes, "notes", "n", "", T("Filter by notes"))
	cobraCmd.Flags().IntVarP(&thisCmd.Order, "order", "o", 0, T("Filter by ID of the order that purchased the block storage"))
	cobraCmd.Flags().StringVar(&thisCmd.SortBy, "sortby", "id", T("Column to sort by, default:id, options are: id,username,datacenter,storage_type,capacity_gb,bytes_used,ip_addr,lunId,active_transactions,created_by"))
	cobraCmd.Flags().StringSliceVar(&thisCmd.UserColumns, "column", []string{}, T("Column to display. Options are: id,username,datacenter,storage_type,capacity_gb,bytes_used,IOPs,ip_addr,lunId,created_by,active_transactions,rep_partner_count,notes. This option can be specified multiple times"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

var maskMap = map[string]string{
	"id":                  "id",
	"username":            "username",
	"datacenter":          "serviceResource.datacenter.name",
	"storage_type":        "storageType.keyName",
	"capacity_gb":         "capacityGb",
	"bytes_used":          "bytesUsed",
	"IOPs":                "iops",
	"ip_addr":             "serviceResourceBackendIpAddress",
	"lunId":               "lunId",
	"active_transactions": "activeTransactionCount",
	"created_by":          "billingItem.orderItem.order.userRecord.username",
	"rep_partner_count":   "replicationPartnerCount",
	"notes":               "notes",
}

func (cmd *VolumeListCommand) Run(args []string) error {
	sortby := cmd.SortBy

	defaultColumns := []string{"id", "username", "datacenter", "storage_type", "capacity_gb", "bytes_used", "IOPs", "ip_addr", "lunId", "active_transactions", "rep_partner_count", "notes"}
	optionalColumns := []string{"notes", "active_transactions", "created_by", "ip_addr"}
	sortColumns := []string{"id", "username", "datacenter", "storage_type", "capacity_gb", "bytes_used", "ip_addr", "lunId", "active_transactions", "created_by"}

	showColumns, err := utils.ValidateColumns2(sortby, cmd.UserColumns, defaultColumns, optionalColumns, sortColumns)
	if err != nil {
		return err
	}

	mask := utils.GetMask(maskMap, showColumns, sortby)

	outputFormat := cmd.GetOutputFlag()

	blockVolumes, err := cmd.StorageManager.ListVolumes("block", cmd.Datacenter, cmd.Username, cmd.StorageType, cmd.Notes, cmd.Order, mask)
	if err != nil {
		return slErr.NewAPIError(T("Failed to list volumes on your account.\n"), err.Error(), 2)
	}

	if sortby == "id" || sortby == "ID" {
		sort.Sort(utils.VolumeById(blockVolumes))
	} else if sortby == "username" {
		sort.Sort(utils.VolumeByUsername(blockVolumes))
	} else if sortby == "datacenter" {
		sort.Sort(utils.VolumeByDatacenter(blockVolumes))
	} else if sortby == "storage_type" {
		sort.Sort(utils.VolumeByStorageType(blockVolumes))
	} else if sortby == "capacity_gb" {
		sort.Sort(utils.VolumeByCapacity(blockVolumes))
	} else if sortby == "bytes_used" {
		sort.Sort(utils.VolumeByBytesUsed(blockVolumes))
	} else if sortby == "ip_addr" {
		sort.Sort(utils.VolumeByIPAddress(blockVolumes))
	} else if sortby == "lunId" {
		sort.Sort(utils.VolumeByLunId(blockVolumes))
	} else if sortby == "active_transactions" {
		sort.Sort(utils.VolumeByTxnCount(blockVolumes))
	} else if sortby == "created_by" {
		sort.Sort(utils.VolumeByCreatedBy(blockVolumes))
	} else {
		return slErr.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.", map[string]interface{}{"Column": sortby}))
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, blockVolumes)
	}

	table := cmd.UI.Table(utils.GetColumnHeader(showColumns))
	for _, blockVolume := range blockVolumes {
		values := make(map[string]string)
		values["id"] = utils.FormatIntPointer(blockVolume.Id)
		values["username"] = utils.FormatStringPointer(blockVolume.Username)
		if blockVolume.ServiceResource != nil && blockVolume.ServiceResource.Datacenter != nil {
			values["datacenter"] = utils.FormatStringPointer(blockVolume.ServiceResource.Datacenter.Name)
		} else {
			values["datacenter"] = "-"
		}

		if blockVolume.StorageType != nil {
			values["storage_type"] = strings.ToLower(utils.FormatStringPointer(blockVolume.StorageType.KeyName))
		} else {
			values["storage_type"] = "-"
		}

		values["capacity_gb"] = utils.FormatIntPointer(blockVolume.CapacityGb)
		values["bytes_used"] = utils.FormatStringPointer(blockVolume.BytesUsed)
		values["IOPs"] = utils.FormatStringPointer(blockVolume.Iops)
		values["ip_addr"] = utils.FormatStringPointer(blockVolume.ServiceResourceBackendIpAddress)
		values["lunId"] = utils.FormatStringPointer(blockVolume.LunId)
		values["active_transactions"] = utils.FormatUIntPointer(blockVolume.ActiveTransactionCount)
		if blockVolume.BillingItem != nil && blockVolume.BillingItem.OrderItem != nil && blockVolume.BillingItem.OrderItem.Order != nil && blockVolume.BillingItem.OrderItem.Order.UserRecord != nil {
			values["created_by"] = utils.FormatStringPointer(blockVolume.BillingItem.OrderItem.Order.UserRecord.Username)
		} else {
			values["created_by"] = "-"
		}
		values["rep_partner_count"] = utils.FormatUIntPointer(blockVolume.ReplicationPartnerCount)
		if blockVolume.Notes != nil {
			decodedValue, err := url.QueryUnescape(utils.FormatStringPointer(blockVolume.Notes))
			if err != nil {
				return slErr.NewAPIError(T("Failed to decoded the note.\n"), err.Error(), 2)
			}
			values["notes"] = decodedValue
		} else {
			values["notes"] = "-"
		}
		row := make([]string, len(showColumns))
		for i, col := range showColumns {
			row[i] = values[col]
		}
		table.Add(row...)
	}
	table.Print()

	return nil
}

package block

import (
	"sort"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type VolumeListCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewVolumeListCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *VolumeListCommand) {
	return &VolumeListCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

var maskMap = map[string]string{
	"id":                  "id",
	"username":            "username",
	"datacenter":          "serviceResource.datacenter.name",
	"storage_type":        "storageType.keyName",
	"capacity_gb":         "capacityGb",
	"bytes_used":          "bytesUsed",
	"ip_addr":             "serviceResourceBackendIpAddress",
	"lunId":               "lunId",
	"active_transactions": "activeTransactionCount",
	"created_by":          "billingItem.orderItem.order.userRecord.username",
	"notes":               "notes",
}

func (cmd *VolumeListCommand) Run(c *cli.Context) error {
	sortby := c.String("sortby")
	if sortby == "" {
		sortby = "id"
	}

	var columns []string
	if c.IsSet("column") {
		columns = c.StringSlice("column")
	} else if c.IsSet("columns") {
		columns = c.StringSlice("columns")
	}

	defaultColumns := []string{"id", "username", "datacenter", "storage_type", "capacity_gb", "bytes_used", "lunId"}
	optionalColumns := []string{"notes", "active_transactions", "created_by", "ip_addr"}
	sortColumns := []string{"id", "username", "datacenter", "storage_type", "capacity_gb", "bytes_used", "ip_addr", "lunId", "active_transactions", "created_by"}

	showColumns, err := utils.ValidateColumns(sortby, columns, defaultColumns, optionalColumns, sortColumns, c)
	if err != nil {
		return err
	}

	mask := utils.GetMask(maskMap, showColumns, sortby)

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	blockVolumes, err := cmd.StorageManager.ListVolumes("block", c.String("datacenter"), c.String("username"), c.String("storage-type"), c.Int("order"), mask)
	if err != nil {
		return cli.NewExitError(T("Failed to list volumes on your account.\n")+err.Error(), 2)
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
		return errors.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.", map[string]interface{}{"Column": sortby}))
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
		values["ip_addr"] = utils.FormatStringPointer(blockVolume.ServiceResourceBackendIpAddress)
		values["lunId"] = utils.FormatStringPointer(blockVolume.LunId)
		values["active_transactions"] = utils.FormatUIntPointer(blockVolume.ActiveTransactionCount)
		if blockVolume.BillingItem != nil && blockVolume.BillingItem.OrderItem != nil && blockVolume.BillingItem.OrderItem.Order != nil && blockVolume.BillingItem.OrderItem.Order.UserRecord != nil {
			values["created_by"] = utils.FormatStringPointer(blockVolume.BillingItem.OrderItem.Order.UserRecord.Username)
		} else {
			values["created_by"] = "-"
		}

		if blockVolume.Notes != nil {
			values["notes"] = utils.FormatStringPointer(blockVolume.Notes)
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

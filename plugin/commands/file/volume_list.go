package file

import (
	"sort"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
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

func FileVolumeListMetaData() cli.Command {
	return cli.Command{
		Category:    "file",
		Name:        "volume-list",
		Description: T("List file storage"),
		Usage: T(`${COMMAND_NAME} sl file volume-list [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file volume-list -d dal09 -t endurance --sortby capacity_gb
   This command lists all endurance volumes on current account that are located at dal09, and sorts them by capacity.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "u,username",
				Usage: T("Filter by volume username"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Filter by datacenter shortname"),
			},
			cli.StringFlag{
				Name:  "t,storage-type",
				Usage: T("Filter by type of storage volume, options are: performance,endurance"),
			},
			cli.IntFlag{
				Name:  "o,order",
				Usage: T("Filter by ID of the order that purchased the file storage"),
			},
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by, default:id, options are: id,username,datacenter,storage_type,capacity_gb,bytes_used,ip_addr,active_transactions,created_by,mount_addr"),
			},
			cli.StringSliceFlag{
				Name:  "column",
				Usage: T("Column to display. Options are: id,username,datacenter,storage_type,capacity_gb,bytes_used,IOPs,ip_addr,lunId,created_by,active_transactions,rep_partner_count,notes. This option can be specified multiple times"),
			},
			cli.StringSliceFlag{
				Name:   "columns",
				Hidden: true,
			},
			metadata.OutputFlag(),
		},
	}
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

	defaultColumns := []string{"id", "username", "datacenter", "storage_type", "capacity_gb", "bytes_used", "IOPs", "ip_addr", "lunId", "active_transactions", "rep_partner_count", "notes"}
	optionalColumns := []string{"created_by", "bytes_used", "active_transactions", "notes"}
	sortColumns := []string{"id", "username", "datacenter", "storage_type", "capacity_gb", "bytes_used", "ip_addr", "active_transactions", "created_by", "mount_addr"}

	showColumns, err := utils.ValidateColumns(sortby, columns, defaultColumns, optionalColumns, sortColumns, c)
	if err != nil {
		return err
	}

	mask := utils.GetMask(maskMap, showColumns, sortby)

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	fileVolumes, err := cmd.StorageManager.ListVolumes("file", c.String("datacenter"), c.String("username"), c.String("storage-type"), c.Int("order"), mask)
	if err != nil {
		return cli.NewExitError(T("Failed to list volumes on your account.\n")+err.Error(), 2)
	}

	if sortby == "id" || sortby == "ID" {
		sort.Sort(utils.VolumeById(fileVolumes))
	} else if sortby == "username" {
		sort.Sort(utils.VolumeByUsername(fileVolumes))
	} else if sortby == "datacenter" {
		sort.Sort(utils.VolumeByDatacenter(fileVolumes))
	} else if sortby == "storage_type" {
		sort.Sort(utils.VolumeByStorageType(fileVolumes))
	} else if sortby == "capacity_gb" {
		sort.Sort(utils.VolumeByCapacity(fileVolumes))
	} else if sortby == "bytes_used" {
		sort.Sort(utils.VolumeByBytesUsed(fileVolumes))
	} else if sortby == "ip_addr" {
		sort.Sort(utils.VolumeByIPAddress(fileVolumes))
	} else if sortby == "active_transactions" {
		sort.Sort(utils.VolumeByTxnCount(fileVolumes))
	} else if sortby == "created_by" {
		sort.Sort(utils.VolumeByCreatedBy(fileVolumes))
	} else if sortby == "mount_addr" {
		sort.Sort(utils.VolumeByMountAddr(fileVolumes))
	} else {
		return errors.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.", map[string]interface{}{"Column": sortby}))
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, fileVolumes)
	}

	table := cmd.UI.Table(utils.GetColumnHeader(showColumns))
	for _, fileVolume := range fileVolumes {
		values := make(map[string]string)
		values["id"] = utils.FormatIntPointer(fileVolume.Id)
		values["username"] = utils.FormatStringPointer(fileVolume.Username)
		if fileVolume.ServiceResource != nil && fileVolume.ServiceResource.Datacenter != nil {
			values["datacenter"] = utils.FormatStringPointer(fileVolume.ServiceResource.Datacenter.Name)
		} else {
			values["datacenter"] = "-"
		}

		if fileVolume.StorageType != nil {
			values["storage_type"] = strings.ToLower(utils.FormatStringPointer(fileVolume.StorageType.KeyName))
		} else {
			values["storage_type"] = "-"
		}

		values["capacity_gb"] = utils.FormatIntPointer(fileVolume.CapacityGb)
		values["bytes_used"] = utils.FormatStringPointer(fileVolume.BytesUsed)
		values["IOPs"] = utils.FormatStringPointer(fileVolume.Iops)
		values["ip_addr"] = utils.FormatStringPointer(fileVolume.ServiceResourceBackendIpAddress)
		values["lunId"] = utils.FormatStringPointer(fileVolume.LunId)
		values["active_transactions"] = utils.FormatUIntPointer(fileVolume.ActiveTransactionCount)
		if fileVolume.BillingItem != nil && fileVolume.BillingItem.OrderItem != nil && fileVolume.BillingItem.OrderItem.Order != nil && fileVolume.BillingItem.OrderItem.Order.UserRecord != nil {
			values["created_by"] = utils.FormatStringPointer(fileVolume.BillingItem.OrderItem.Order.UserRecord.Username)
		} else {
			values["created_by"] = "-"
		}
		values["rep_partner_count"] = utils.FormatUIntPointer(fileVolume.ReplicationPartnerCount)
		values["mount_addr"] = utils.FormatStringPointer(fileVolume.FileNetworkMountAddress)

		if fileVolume.Notes != nil {
			values["notes"] = utils.FormatStringPointer(fileVolume.Notes)
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

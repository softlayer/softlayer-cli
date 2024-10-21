package block

import (
	"github.com/spf13/cobra"
	"sort"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type AccessListCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	UserColumn     []string
	Sortby         string
}

func NewAccessListCommand(sl *metadata.SoftlayerStorageCommand) *AccessListCommand {
	thisCmd := &AccessListCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "access-list " + T("IDENTIFIER"),
		Short: T("List hosts that are authorized to access the volume."),
		Long: T(`EXAMPLE:
   ${COMMAND_NAME} sl {{.storageType}} access-list 12345678 --sortby id 
   This command lists all hosts that are authorized to access volume with ID 12345678 and sorts them by ID.`, sl.StorageI18n),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Sortby, "sortby", "id", T("Column to sort by. Options are: id, name, type, private_ip_address, source_subnet, host_iqn, username, password, allowed_host_id."))
	cobraCmd.Flags().StringSliceVar(&thisCmd.UserColumn, "column", []string{}, T("Column to display. Options are: id, name, type, private_ip_address, source_subnet, host_iqn, username, password, allowed_host_id. This option can be specified multiple times."))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *AccessListCommand) Run(args []string) error {

	volumeID, err := cmd.StorageManager.GetVolumeId(args[0], cmd.StorageType)
	if err != nil {
		return err
	}

	sortby := cmd.Sortby

	outputFormat := cmd.GetOutputFlag()

	defaultColumns := []string{"id", "name", "type", "private_ip_address", "source_subnet", "host_iqn", "username", "password", "allowed_host_id"}
	optionalColumns := []string{}
	sortColumns := []string{"id", "name", "type", "private_ip_address", "source_subnet", "host_iqn", "username", "password", "allowed_host_id"}
	showColumns, err := utils.ValidateColumns2(sortby, cmd.UserColumn, defaultColumns, optionalColumns, sortColumns)
	if err != nil {
		return err
	}

	volume, err := cmd.StorageManager.GetVolumeAccessList(volumeID)
	if err != nil {
		return slErr.NewAPIError(T("Failed to get access list for volume {{.VolumeID}}.\n", map[string]interface{}{"VolumeID": volumeID}), err.Error(), 2)
	}

	accessList := []utils.Access{}
	for _, vs := range volume.AllowedVirtualGuests {
		access := utils.Access{}

		access.ID = utils.FormatIntPointer(vs.Id)
		access.Name = utils.FormatStringPointerName(vs.Hostname) + "." + utils.FormatStringPointerName(vs.Domain)
		access.Type = T("VIRTUAL")
		access.PrivateIPAddress = utils.FormatStringPointer(vs.PrimaryBackendIpAddress)

		if vs.AllowedHost != nil {
			access.SourceSubnet = utils.FormatStringPointer(vs.AllowedHost.SourceSubnet)
			access.HostIQN = utils.FormatStringPointer(vs.AllowedHost.Name)
			access.AllowedHostID = utils.FormatIntPointer(vs.AllowedHost.Id)
			if vs.AllowedHost.Credential != nil {
				credentials := *vs.AllowedHost.Credential
				access.UserName = utils.FormatStringPointer(credentials.Username)
				access.Password = utils.FormatStringPointer(credentials.Password)
			}
		}
		accessList = append(accessList, access)
	}

	for _, hw := range volume.AllowedHardware {
		access := utils.Access{}
		access.ID = utils.FormatIntPointer(hw.Id)
		access.Name = utils.FormatStringPointerName(hw.Hostname) + "." + utils.FormatStringPointerName(hw.Domain)
		access.Type = T("HARDWARE")
		access.PrivateIPAddress = utils.FormatStringPointer(hw.PrimaryBackendIpAddress)

		if hw.AllowedHost != nil {
			access.SourceSubnet = utils.FormatStringPointer(hw.AllowedHost.SourceSubnet)
			access.HostIQN = utils.FormatStringPointer(hw.AllowedHost.Name)
			access.AllowedHostID = utils.FormatIntPointer(hw.AllowedHost.Id)
			if hw.AllowedHost.Credential != nil {
				credentials := *hw.AllowedHost.Credential
				access.UserName = utils.FormatStringPointer(credentials.Username)
				access.Password = utils.FormatStringPointer(credentials.Password)
			}
		}
		accessList = append(accessList, access)
	}

	for _, sn := range volume.AllowedSubnets {

		access := utils.Access{}

		access.ID = utils.FormatIntPointer(sn.Id)

		if utils.FormatStringPointerName(sn.Note) != "" {
			access.Name = utils.FormatStringPointerName(sn.NetworkIdentifier) + "/" + utils.FormatIntPointerName(sn.Cidr) + "(" + utils.FormatStringPointerName(sn.Note) + ")"
		} else {
			access.Name = utils.FormatStringPointerName(sn.NetworkIdentifier) + "/" + utils.FormatIntPointerName(sn.Cidr)
		}

		access.Type = T("SUBNET")

		if sn.EndPointIpAddress != nil {
			access.PrivateIPAddress = utils.FormatStringPointer(sn.EndPointIpAddress.IpAddress)
		}

		if sn.AllowedHost != nil {
			access.SourceSubnet = utils.FormatStringPointer(sn.AllowedHost.SourceSubnet)
			access.HostIQN = utils.FormatStringPointer(sn.AllowedHost.Name)
			access.AllowedHostID = utils.FormatIntPointer(sn.AllowedHost.Id)
			if sn.AllowedHost.Credential != nil {
				credentials := *sn.AllowedHost.Credential
				access.UserName = utils.FormatStringPointer(credentials.Username)
				access.Password = utils.FormatStringPointer(credentials.Password)
			}
		}
		accessList = append(accessList, access)
	}

	for _, ip := range volume.AllowedIpAddresses {
		access := utils.Access{}

		access.ID = utils.FormatIntPointer(ip.Id)

		if utils.FormatStringPointerName(ip.Note) != "" {
			access.Name = utils.FormatStringPointerName(ip.IpAddress) + "(" + utils.FormatStringPointerName(ip.Note) + ")"
		} else {
			access.Name = utils.FormatStringPointerName(ip.IpAddress)
		}

		access.Type = T("IP")
		access.PrivateIPAddress = utils.FormatStringPointer(ip.IpAddress)

		if ip.AllowedHost != nil {
			access.SourceSubnet = utils.FormatStringPointer(ip.AllowedHost.SourceSubnet)
			access.HostIQN = utils.FormatStringPointer(ip.AllowedHost.Name)
			access.AllowedHostID = utils.FormatIntPointer(ip.AllowedHost.Id)
			if ip.AllowedHost.Credential != nil {
				credentials := *ip.AllowedHost.Credential
				access.UserName = utils.FormatStringPointer(credentials.Username)
				access.Password = utils.FormatStringPointer(credentials.Password)
			}
		}
		accessList = append(accessList, access)
	}

	if sortby == "id" || sortby == "ID" {
		sort.Sort(utils.AccessByID(accessList))
	} else if sortby == "name" {
		sort.Sort(utils.AccessByName(accessList))
	} else if sortby == "type" {
		sort.Sort(utils.AccessByType(accessList))
	} else if sortby == "private_ip_address" {
		sort.Sort(utils.AccessByPrivateIPAddress(accessList))
	} else if sortby == "source_subnet" {
		sort.Sort(utils.AccessBySourceSubnet(accessList))
	} else if sortby == "host_iqn" {
		sort.Sort(utils.AccessByHostIQN(accessList))
	} else if sortby == "username" {
		sort.Sort(utils.AccessByUserName(accessList))
	} else if sortby == "password" {
		sort.Sort(utils.AccessByPassword(accessList))
	} else if sortby == "allowed_host_id" {
		sort.Sort(utils.AccessByAllowedHostID(accessList))
	} else if sortby == "" {
		//do nothing
	} else {
		return slErr.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.", map[string]interface{}{"Column": sortby}))
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, accessList)
	}

	table := cmd.UI.Table(utils.GetColumnHeader(showColumns))
	for _, access := range accessList {
		mapValue, err := utils.StructToMap(access)
		if err != nil {
			return err
		}
		row := make([]string, len(showColumns))
		for i, col := range showColumns {
			row[i] = mapValue[col]
		}
		table.Add(row...)
	}

	table.Print()
	return nil
}

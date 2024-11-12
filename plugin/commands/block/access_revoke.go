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

type AccessRevokeCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	NetworkManager managers.NetworkManager
	Hardware_id    []int
	Virtual_id     []int
	Ip_address_id  []int
	Ip_address     []string
	Subnet_id	   int
}

func NewAccessRevokeCommand(sl *metadata.SoftlayerStorageCommand) *AccessRevokeCommand {
	thisCmd := &AccessRevokeCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
		NetworkManager:          managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "access-revoke " + T("IDENTIFIER"),
		Short: T("Revoke authorization for hosts that are accessing a specific volume."),
		Long: T(`EXAMPLE:
   ${COMMAND_NAME} sl {{.storageType}} access-revoke 12345678 --virtual-id 87654321
   This command revokes access of virtual server with ID 87654321 to volume with ID 12345678.`, sl.StorageI18n) + "\n" +
			T(`
   ${COMMAND_NAME} sl {{.storageType}} access-revoke 5555 --subnet-id 1111
   This command removes subnet with id 1111 from the Allowed Host with id 5555. Use 'access-list' to find this id.`, sl.StorageI18n),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().IntSliceVarP(&thisCmd.Hardware_id, "hardware-id", "d", []int{},
		T("The ID of one hardware server to revoke"))
	cobraCmd.Flags().IntSliceVarP(&thisCmd.Virtual_id, "virtual-id", "v", []int{},
		T("The ID of one virtual server to revoke"))
	cobraCmd.Flags().IntSliceVarP(&thisCmd.Ip_address_id, "ip-address-id", "i", []int{},
		T("The ID of one IP address to revoke"))
	cobraCmd.Flags().StringSliceVarP(&thisCmd.Ip_address, "ip-address", "p", []string{},
		T("An IP address to revoke"))
	cobraCmd.Flags().IntVarP(&thisCmd.Subnet_id, "subnet-id", "s", 0,
		T("A Subnet Id. With this option IDENTIFIER should be an 'allowed_host_id' from the access-list command."))
	thisCmd.Command = cobraCmd

	return thisCmd
}

func (cmd *AccessRevokeCommand) Run(args []string) error {

	// Subnets have to get added to an existing authorized host.
	if cmd.Subnet_id > 0 {
		hostId, err := strconv.Atoi(args[0])
		if err != nil {
			return slErr.NewInvalidSoftlayerIdInputError("Allowed Host IDENTIFIER")
		}
		return cmd.RemoveSubnetFromHost(hostId)
	}

	volumeID, err := cmd.StorageManager.GetVolumeId(args[0], cmd.StorageType)
	if err != nil {
		return err
	}

	if len(cmd.Hardware_id) == 0 && len(cmd.Virtual_id) == 0 && len(cmd.Ip_address_id) == 0 && len(cmd.Ip_address) == 0 {
		return slErr.NewInvalidUsageError(T("One of -d | --hardware-id, -v | --virtual-id, -i | --ip-address-id and -p | --ip-address must be specified."))
	}

	IPIds := cmd.Ip_address_id
	IPs := cmd.Ip_address
	if len(IPs) > 0 {
		for _, ip := range IPs {
			ipRecord, err := cmd.NetworkManager.IPLookup(ip)
			if err != nil {
				return slErr.NewAPIError(T("IP address {{.IP}} is not found on your account.Please confirm IP and try again.\n",
					map[string]interface{}{"IP": ip}), err.Error(), 2)
			}
			if ipRecord.Id != nil {
				IPIds = append(IPIds, *ipRecord.Id)
			}

		}
	}
	_, err = cmd.StorageManager.DeauthorizeHostToVolume(volumeID, cmd.Hardware_id, cmd.Virtual_id, IPIds, nil)
	if err != nil {
		subs := map[string]interface{}{"VolumeID": volumeID}
		return slErr.NewAPIError(T("Failed to revoke access to volume {{.VolumeID}}.\n", subs), err.Error(), 2)
	}
	cmd.UI.Ok()
	subs := map[string]interface{}{
		"VolumeId": volumeID,
		"SL_ID": 0,
		"SL_Object": "",
	}
	for _, sl_id := range cmd.Virtual_id {
		subs["SL_Object"] = T("Virtual Server")
		subs["SL_ID"] = sl_id
		cmd.UI.Print(T("Access to {{.VolumeId}} was revoked for {{.SL_Object}} {{.SL_ID}}.", subs))
	}
	for _, sl_id := range cmd.Hardware_id {
		subs["SL_Object"] = T("Hardware Server")
		subs["SL_ID"] = sl_id
		cmd.UI.Print(T("Access to {{.VolumeId}} was revoked for {{.SL_Object}} {{.SL_ID}}.", subs))
	}
	for _, sl_id := range IPIds {
		subs["SL_Object"] = T("IP Address")
		subs["SL_ID"] = sl_id
		cmd.UI.Print(T("Access to {{.VolumeId}} was revoked for {{.SL_Object}} {{.SL_ID}}.", subs))
	}
	return nil
}

func (cmd *AccessRevokeCommand) RemoveSubnetFromHost(host_id int) error {

	outputFormat := cmd.GetOutputFlag()
	subnet_ids := []int{cmd.Subnet_id}
	resp, err := cmd.StorageManager.RemoveSubnetsFromAcl(host_id, subnet_ids)
	if err != nil {
		return err
	}
	// If the API returns an empty array, that means it didn't add the subnet we asked for.
	// Likely because ISCSI Isolation is disabled on the account.
	// ibmcloud sl call-api SoftLayer_Account getObject --mask="mask[id,iscsiIsolationDisabled]"
	if len(resp) == 0 || utils.IntInSlice(cmd.Subnet_id, resp) == -1 {
		subs := map[string]interface{}{"subnetID": cmd.Subnet_id, "accessID": host_id}
		return slErr.NewAPIError(
			T("Failed to remove subnet id: {{.subnetID}} from allowed host id: {{.accessID}}", subs), "", 2,
		)
	}
	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	subs := map[string]interface{}{
		"VolumeId": host_id,
		"SL_ID": cmd.Subnet_id,
		"SL_Object": T("Subnet"),
	}
	cmd.UI.Print(T("Access to {{.VolumeId}} was revoked for {{.SL_Object}} {{.SL_ID}}.", subs))
	return nil
}
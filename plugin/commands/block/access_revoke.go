package block

import (
	"strconv"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type AccessRevokeCommand struct {
	*metadata.SoftlayerCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	NetworkManager managers.NetworkManager
	Hardware_id    []int
	Virtual_id     []int
	Ip_address_id  []int
	Ip_address     []string
}

func NewAccessRevokeCommand(sl *metadata.SoftlayerCommand) *AccessRevokeCommand {
	thisCmd := &AccessRevokeCommand{
		SoftlayerCommand: sl,
		StorageManager:   managers.NewStorageManager(sl.Session),
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "access-revoke " + T("IDENTIFIER"),
		Short: T("Revoke authorization for hosts that are accessing a specific volume"),
		Long: T(`${COMMAND_NAME} sl block access-revoke VOLUME_ID [OPTIONS]
		
EXAMPLE:
   ${COMMAND_NAME} sl block access-revoke 12345678 --virtual-id 87654321
   This command revokes access of virtual server with ID 87654321 to volume with ID 12345678.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().IntSliceVarP(&thisCmd.Hardware_id, "hardware-id", "d", []int{}, T("The ID of one hardware server to revoke"))
	cobraCmd.Flags().IntSliceVarP(&thisCmd.Virtual_id, "virtual-id", "v", []int{}, T("The ID of one virtual server to revoke"))
	cobraCmd.Flags().IntSliceVarP(&thisCmd.Ip_address_id, "ip-address-id", "i", []int{}, T("The ID of one IP address to revoke"))
	cobraCmd.Flags().StringSliceVarP(&thisCmd.Ip_address, "ip-address", "p", []string{}, T("An IP address to revoke"))
	thisCmd.Command = cobraCmd

	return thisCmd
}

func (cmd *AccessRevokeCommand) Run(args []string) error {
	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
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
		return slErr.NewAPIError(T("Failed to revoke access to volume {{.VolumeID}}.\n", map[string]interface{}{"VolumeID": volumeID}), err.Error(), 2)
	}
	cmd.UI.Ok()
	for _, vsID := range cmd.Virtual_id {
		cmd.UI.Print(T("Access to {{.VolumeId}} was revoked for virtual server {{.VsID}}.", map[string]interface{}{"VolumeId": volumeID, "VsID": vsID}))
	}
	for _, hwID := range cmd.Hardware_id {
		cmd.UI.Print(T("Access to {{.VolumeId}} was revoked for hardware server {{.HwID}}.", map[string]interface{}{"VolumeId": volumeID, "HwID": hwID}))
	}
	for _, ip := range IPIds {
		cmd.UI.Print(T("Access to {{.VolumeId}} was revoked for IP address {{.IP}}.", map[string]interface{}{"VolumeId": volumeID, "IP": ip}))
	}
	return nil
}

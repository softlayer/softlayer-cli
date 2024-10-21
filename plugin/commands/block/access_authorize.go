package block

import (
	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type AccessAuthorizeCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	NetworkManager managers.NetworkManager
	Hardware_id    []int
	Virtual_id     []int
	Ip_address_id  []int
	Ip_address     []string
}

func NewAccessAuthorizeCommand(sl *metadata.SoftlayerStorageCommand) *AccessAuthorizeCommand {
	thisCmd := &AccessAuthorizeCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
		NetworkManager:          managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "access-authorize " + T("IDENTIFIER"),
		Short: T("Authorize hosts to access a given volume."),
		Long: T(`EXAMPLE:
   ${COMMAND_NAME} sl {{.storageType}} access-authorize 12345678 --virtual-id 87654321
   This command authorizes virtual server with ID 87654321 to access volume with ID 12345678.`, sl.StorageI18n),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().IntSliceVarP(&thisCmd.Hardware_id, "hardware-id", "d", []int{}, T("The ID of one hardware server to authorize."))
	cobraCmd.Flags().IntSliceVarP(&thisCmd.Virtual_id, "virtual-id", "v", []int{}, T("The ID of one virtual server to authorize."))
	cobraCmd.Flags().IntSliceVarP(&thisCmd.Ip_address_id, "ip-address-id", "i", []int{}, T("The ID of one IP address to authorize."))
	cobraCmd.Flags().StringSliceVarP(&thisCmd.Ip_address, "ip-address", "p", []string{}, T("An IP address to authorize."))
	thisCmd.Command = cobraCmd

	return thisCmd
}

func (cmd *AccessAuthorizeCommand) Run(args []string) error {

	volumeID, err := cmd.StorageManager.GetVolumeId(args[0], cmd.StorageType)
	if err != nil {
		return err
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

	outputFormat := cmd.GetOutputFlag()

	resp, err := cmd.StorageManager.AuthorizeHostToVolume(volumeID, cmd.Hardware_id, cmd.Virtual_id, IPIds, nil)
	if err != nil {
		subs := map[string]interface{}{"VolumeID": volumeID}
		return slErr.NewAPIError(T("Failed to authorize host to volume {{.VolumeID}}.\n", subs), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	for _, vsID := range cmd.Virtual_id {
		cmd.UI.Print(T("The virtual server {{.VsID}} was authorized to access {{.VolumeId}}.", map[string]interface{}{"VolumeId": volumeID, "VsID": vsID}))
	}
	for _, hwID := range cmd.Hardware_id {
		cmd.UI.Print(T("The hardware server {{.HwID}} was authorized to access {{.VolumeId}}.", map[string]interface{}{"VolumeId": volumeID, "HwID": hwID}))
	}
	for _, ip := range IPIds {
		cmd.UI.Print(T("The IP address {{.IP}} was authorized to access {{.VolumeId}}.", map[string]interface{}{"VolumeId": volumeID, "IP": ip}))
	}
	return nil
}

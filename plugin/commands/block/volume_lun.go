package block

import (
	"strconv"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type VolumeLunCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
}

func NewVolumeLunCommand(sl *metadata.SoftlayerStorageCommand) *VolumeLunCommand {
	thisCmd := &VolumeLunCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "volume-set-lun-id " + T("IDENTIFIER") + " " + T("LUN_ID"),
		Short: T("Set the LUN ID on an existing block storage volume"),
		Long: T("Set the LUN ID on an existing block storage volume") + " " + T(`${COMMAND_NAME} sl {{.storageType}} volume-set-lun-id VOLUME_ID LUN_ID

	The LUN ID only takes effect during the Host Authorization process. It is
	recommended (but not necessary) to de-authorize all hosts before using
	this method. See "block access-revoke".
	VOLUME_ID - the volume ID on which to set the LUN ID
	LUN_ID - recommended range is an integer between 0 and 255. Advanced users
	can use an integer between 0 and 4095`, sl.StorageI18n),
		Args: metadata.TwoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VolumeLunCommand) Run(args []string) error {

	volumeID, err := cmd.StorageManager.GetVolumeId(args[0], cmd.StorageType)
	if err != nil {
		return err
	}

	lunId, err := strconv.Atoi(args[1])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("LUN ID")
	}
	prop, err := cmd.StorageManager.SetLunId(volumeID, lunId)

	subs := map[string]interface{}{"VolumeID": volumeID, "VolumeId": volumeID}
	if err != nil {
		return slErr.NewAPIError(T("Failed to set LUN ID for volume {{.VolumeID}}.\n", subs), err.Error(), 2)
	}
	if prop.Value != nil {
		newLunId, err := strconv.Atoi(*prop.Value)
		if err == nil && newLunId == lunId {
			cmd.UI.Ok()
			cmd.UI.Print(T("Block volume {{.VolumeId}} is reporting LUN ID {{.LunID}}.",
				map[string]interface{}{"VolumeId": volumeID, "LunID": lunId}))
			return nil
		}
	}
	cmd.UI.Failed(T("Failed to confirm the new LUN ID on volume {{.VolumeId}}.", subs))
	return nil
}

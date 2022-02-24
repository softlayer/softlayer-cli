package block

import (
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type VolumeCancelCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewVolumeCancelCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *VolumeCancelCommand) {
	return &VolumeCancelCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func BlockVolumeCancelMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "volume-cancel",
		Description: T("Cancel an existing block storage volume"),
		Usage: T(`${COMMAND_NAME} sl block volume-cancel VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl block volume-cancel 12345678 --immediate -f 
   This command cancels volume with ID 12345678 immediately and without asking for confirmation.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "reason",
				Usage: T("An optional reason for cancellation"),
			},
			cli.BoolFlag{
				Name:  "immediate",
				Usage: T("Cancel the block storage volume immediately instead of on the billing anniversary"),
			},
			metadata.ForceFlag(),
		},
	}
}

func (cmd *VolumeCancelCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	if !c.IsSet("f") {
		confirm, err := cmd.UI.Confirm(T("This will cancel the block volume: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": volumeID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	immediate := c.IsSet("immediate")
	err = cmd.StorageManager.CancelVolume("block", volumeID, c.String("reason"), immediate)
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return cli.NewExitError(T("Unable to find volume with ID {{.ID}}.\n", map[string]interface{}{"ID": volumeID})+err.Error(), 0)
		}
		return cli.NewExitError(T("Failed to cancel block volume: {{.ID}}.\n", map[string]interface{}{"ID": volumeID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	if immediate {
		cmd.UI.Print(T("Block volume {{.VolumeId}} has been marked for immediate cancellation.", map[string]interface{}{"VolumeId": volumeID}))
	} else {
		cmd.UI.Print(T("Block volume {{.VolumeId}} has been marked for cancellation.", map[string]interface{}{"VolumeId": volumeID}))
	}
	return nil
}

package block

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type VolumeSetNoteCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewVolumeSetNoteCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *VolumeSetNoteCommand) {
	return &VolumeSetNoteCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func (cmd *VolumeSetNoteCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	if !c.IsSet("n") {
		return errors.NewInvalidUsageError(T("This command requires note flag."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	successful, err := cmd.StorageManager.VolumeSetNote(volumeID, c.String("note"))
	if err != nil {
		return cli.NewExitError(T("Error occurred while note was adding in block volume {{.VolumeID}}.\n",
			map[string]interface{}{"VolumeID": volumeID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, successful)
	}

	response := ""
	if successful {
		cmd.UI.Ok()
		response = "The note was set successfully"
	} else {
		response = "Note could not be set! Please verify your options and try again."
	}

	cmd.UI.Print(T(response))
	return nil
}

func BlockVolumeSetNoteMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "volume-set-note",
		Description: T("Set note for an existing block storage volume."),
		Usage: T(`${COMMAND_NAME} sl block volume-set-note [OPTIONS] VOLUME_ID

EXAMPLE:
   ${COMMAND_NAME} sl block volume-set-note 12345678 --note "this is my note"`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,note",
				Usage: T("Public notes related to a Storage volume  [required]"),
			},
			metadata.OutputFlag(),
		},
	}
}

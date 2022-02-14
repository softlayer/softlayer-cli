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
)

type SnapshotCancelCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewSnapshotCancelCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *SnapshotCancelCommand) {
	return &SnapshotCancelCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func BlockSnapshotCancelMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "snapshot-cancel",
		Description: T("Cancel existing snapshot space for a given volume"),
		Usage: T(`${COMMAND_NAME} sl block snapshot-cancel SNAPSHOT_ID [OPTIONS]
		
EXAMPLE:
   ${COMMAND_NAME} sl block snapshot-cancel 12345678 --immediate -f 
   This command cancels snapshot with ID 12345678 immediately without asking for confirmation.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "reason",
				Usage: T("An optional reason for cancellation"),
			},
			cli.BoolFlag{
				Name:  "immediate",
				Usage: T("Cancel the snapshot space immediately instead of on the billing anniversary"),
			},
			metadata.ForceFlag(),
		},
	}
}

func (cmd *SnapshotCancelCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This will cancel the block volume snapshot space: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": volumeID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	immediate := c.IsSet("immediate")
	err = cmd.StorageManager.CancelSnapshotSpace("block", volumeID, c.String("reason"), immediate)
	if err != nil {
		return cli.NewExitError(T("Failed to cancel snapshot space for volume {{.ID}}.\n", map[string]interface{}{"ID": volumeID})+err.Error(), 2)
	}

	cmd.UI.Ok()
	if immediate {
		cmd.UI.Print(T("Block volume {{.ID}} has been marked for immediate snapshot cancellation.", map[string]interface{}{"ID": volumeID}))
	} else {
		cmd.UI.Print(T("Block volume {{.ID}} has been marked for snapshot cancellation.", map[string]interface{}{"ID": volumeID}))
	}
	return nil
}

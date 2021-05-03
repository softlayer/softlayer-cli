package file

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
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

func (cmd *SnapshotCancelCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This will cancel the file volume snapshot space: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": volumeID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	immediate := c.IsSet("immediate")
	err = cmd.StorageManager.CancelSnapshotSpace("file", volumeID, c.String("reason"), immediate)
	if err != nil {
		return cli.NewExitError(T("Failed to cancel snapshot space for volume {{.ID}}.\n", map[string]interface{}{"ID": volumeID})+err.Error(), 2)
	}

	cmd.UI.Ok()
	if immediate {
		cmd.UI.Print(T("File volume {{.ID}} has been marked for immediate snapshot cancellation.", map[string]interface{}{"ID": volumeID}))
	} else {
		cmd.UI.Print(T("File volume {{.ID}} has been marked for snapshot cancellation.", map[string]interface{}{"ID": volumeID}))
	}
	return nil
}

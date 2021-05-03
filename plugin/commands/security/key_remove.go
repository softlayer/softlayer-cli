package security

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
)

type KeyRemoveCommand struct {
	UI              terminal.UI
	SecurityManager managers.SecurityManager
}

func NewKeyRemoveCommand(ui terminal.UI, securityManager managers.SecurityManager) (cmd *KeyRemoveCommand) {
	return &KeyRemoveCommand{
		UI:              ui,
		SecurityManager: securityManager,
	}
}

func (cmd *KeyRemoveCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	keyID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("SSH Key ID")
	}
	if !c.IsSet("f") {
		confirm, err := cmd.UI.Confirm(T("This will remove SSH key: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": keyID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.SecurityManager.DeleteSSHKey(keyID)
	if err != nil {
		return cli.NewExitError(T("Failed to remove SSH key: {{.ID}}.\n", map[string]interface{}{"ID": keyID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("SSH key {{.ID}} was removed.", map[string]interface{}{"ID": keyID}))
	return nil
}

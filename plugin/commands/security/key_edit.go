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

type KeyEditCommand struct {
	UI              terminal.UI
	SecurityManager managers.SecurityManager
}

func NewKeyEditCommand(ui terminal.UI, securityManager managers.SecurityManager) (cmd *KeyEditCommand) {
	return &KeyEditCommand{
		UI:              ui,
		SecurityManager: securityManager,
	}
}

func (cmd *KeyEditCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	keyID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("SSH Key ID")
	}
	if !c.IsSet("label") && !c.IsSet("note") {
		return errors.NewInvalidUsageError(T("either [--label] or [--note] must be specified to edit SSH key."))
	}
	err = cmd.SecurityManager.EditSSHKey(keyID, c.String("label"), c.String("note"))
	if err != nil {
		return cli.NewExitError(T("Failed to edit SSH key: {{.ID}}.\n", map[string]interface{}{"ID": keyID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("SSH key {{.ID}} was updated.", map[string]interface{}{"ID": keyID}))
	return nil
}

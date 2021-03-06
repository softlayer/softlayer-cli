package security

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

func SecuritySSHKeyRemoveMetaData() cli.Command {
	return cli.Command{
		Category:    "security",
		Name:        "sshkey-remove",
		Description: T("Permanently removes an SSH key"),
		Usage: T(`${COMMAND_NAME} sl security sshkey-remove IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl security sshkey-remove 12345678 -f 
   This command removes the SSH key with ID 12345678 without asking for confirmation.`),
		Flags: []cli.Flag{
			metadata.ForceFlag(),
		},
	}
}

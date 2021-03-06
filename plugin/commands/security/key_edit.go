package security

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
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

func SecuritySSHKeyEditMetaData() cli.Command {
	return cli.Command{
		Category:    "security",
		Name:        "sshkey-edit",
		Description: T("Edit an SSH key"),
		Usage: T(`${COMMAND_NAME} sl security sshkey-edit IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl security sshkey-edit 12345678 --label IBMCloud --note testing
   This command updates the SSH key with ID 12345678 and sets label to "IBMCloud" and note to "testing".`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "label",
				Usage: T("The new label for the key"),
			},
			cli.StringFlag{
				Name:  "note",
				Usage: T("New notes for the key"),
			},
		},
	}
}

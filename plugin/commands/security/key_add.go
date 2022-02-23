package security

import (
	"io/ioutil"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type KeyAddCommand struct {
	UI              terminal.UI
	SecurityManager managers.SecurityManager
}

func NewKeyAddCommand(ui terminal.UI, securityManager managers.SecurityManager) (cmd *KeyAddCommand) {
	return &KeyAddCommand{
		UI:              ui,
		SecurityManager: securityManager,
	}
}

func (cmd *KeyAddCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	if !c.IsSet("k") && !c.IsSet("f") {
		return errors.NewInvalidUsageError(T("either [-k] or [-f|--in-file] is required."))
	}
	if c.IsSet("k") && c.IsSet("f") {
		return errors.NewInvalidUsageError(T("[-k] is not allowed with [-f|--in-file]."))
	}
	var keyText string
	if c.IsSet("k") {
		keyText = c.String("k")
	} else {
		keyFile := c.String("f")
		contents, err := ioutil.ReadFile(keyFile) // #nosec
		if err != nil {
			return cli.NewExitError(T("Failed to read SSH key from file: {{.File}}.\n", map[string]interface{}{"File": keyFile})+err.Error(), 2)
		}
		keyText = string(contents)
	}
	key, err := cmd.SecurityManager.AddSSHKey(keyText, c.Args()[0], c.String("note"))
	if err != nil {
		return cli.NewExitError(T("Failed to add SSH key.\n")+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, key)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("SSH key was added: ") + utils.StringPointertoString(key.Fingerprint))
	return nil
}

func SecuritySSHKeyAddMetaData() cli.Command {
	return cli.Command{
		Category:    "security",
		Name:        "sshkey-add",
		Description: T("Add a new SSH key"),
		Usage: T(`${COMMAND_NAME} sl security sshkey-add LABEL [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl security sshkey-add my_sshkey -f ~/.ssh/id_rsa.pub --note mykey
   This command adds an SSH key from file ~/.ssh/id_rsa.pub with a note "mykey".`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "f,in-file",
				Usage: T("The id_rsa.pub file to import for this key"),
			},
			cli.StringFlag{
				Name:  "k,key",
				Usage: T("The actual SSH key"),
			},
			cli.StringFlag{
				Name:  "note",
				Usage: T("Extra note to be associated with the key"),
			},
			metadata.OutputFlag(),
		},
	}
}

package security

import (
	"io/ioutil"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
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

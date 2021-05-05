package security

import (
	"io/ioutil"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type KeyPrintCommand struct {
	UI              terminal.UI
	SecurityManager managers.SecurityManager
}

func NewKeyPrintCommand(ui terminal.UI, securityManager managers.SecurityManager) (cmd *KeyPrintCommand) {
	return &KeyPrintCommand{
		UI:              ui,
		SecurityManager: securityManager,
	}
}

func (cmd *KeyPrintCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	keyID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("SSH Key ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	key, err := cmd.SecurityManager.GetSSHKey(keyID)
	if err != nil {
		return cli.NewExitError(T("Failed to get SSH Key {{.KeyID}}.\n", map[string]interface{}{"KeyID": keyID})+err.Error(), 1)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, key)
	}

	if file := c.String("f"); file != "" {
		// #nosec G306: write on customer machine
		err = ioutil.WriteFile(file, []byte(utils.FormatStringPointer(key.Key)), 0644)
		if err != nil {
			return cli.NewExitError(T("Failed to write SSH key to file: {{.File}}.\n",
				map[string]interface{}{"File": file})+err.Error(), 1)
		}
	}
	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(key.Id))
	table.Add(T("Label"), utils.FormatStringPointer(key.Label))
	table.Add(T("Notes"), utils.FormatStringPointer(key.Notes))
	table.Print()
	return nil
}

package security

import (
	"sort"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type KeyListCommand struct {
	UI              terminal.UI
	SecurityManager managers.SecurityManager
}

func NewKeyListCommand(ui terminal.UI, securityManager managers.SecurityManager) (cmd *KeyListCommand) {
	return &KeyListCommand{
		UI:              ui,
		SecurityManager: securityManager,
	}
}

func (cmd *KeyListCommand) Run(c *cli.Context) error {

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	keys, err := cmd.SecurityManager.ListSSHKeys("")
	if err != nil {
		return cli.NewExitError(T("Failed to list SSH keys on your account.\n")+err.Error(), 2)
	}

	sortby := c.String("sortby")
	if sortby == "id" || sortby == "ID" {
		sort.Sort(utils.KeyById(keys))
	} else if sortby == "label" {
		sort.Sort(utils.KeyByLabel(keys))
	} else if sortby == "fingerprint" {
		sort.Sort(utils.KeyByFingerprint(keys))
	} else if sortby == "note" {
		sort.Sort(utils.KeyByNotes(keys))
	} else if sortby == "" {
		//do nothing
	} else {
		return errors.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.",
			map[string]interface{}{"Column": sortby}))
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, keys)
	}

	table := cmd.UI.Table([]string{T("ID"), T("label"), T("fingerprint"), T("note")})
	for _, k := range keys {
		table.Add(utils.FormatIntPointer(k.Id),
			utils.FormatStringPointer(k.Label),
			utils.FormatStringPointer(k.Fingerprint),
			utils.FormatStringPointer(k.Notes))
	}
	table.Print()
	return nil
}

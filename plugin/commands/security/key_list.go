package security

import (
	"sort"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
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

func SecuritySSHKeyListMetaData() cli.Command {
	return cli.Command{
		Category:    "security",
		Name:        "sshkey-list",
		Description: T("List SSH keys on your account"),
		Usage: T(`${COMMAND_NAME} sl security sshkey-list [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl security sshkey-list --sortby label
   This command lists all SSH keys on current account and sorts them by label.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by. Options are: id,label,fingerprint,note"),
			},
			metadata.OutputFlag(),
		},
	}
}

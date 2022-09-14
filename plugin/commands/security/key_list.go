package security

import (
	"sort"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type KeyListCommand struct {
	*metadata.SoftlayerCommand
	SecurityManager managers.SecurityManager
	Command         *cobra.Command
	SortBy          string
}

func NewKeyListCommand(sl *metadata.SoftlayerCommand) *KeyListCommand {
	thisCmd := &KeyListCommand{
		SoftlayerCommand: sl,
		SecurityManager:  managers.NewSecurityManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "sshkey-list",
		Short: T("List SSH keys on your account"),
		Long: T(`${COMMAND_NAME} sl security sshkey-list [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl security sshkey-list --sortby label
   This command lists all SSH keys on current account and sorts them by label.`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.SortBy, "sortby", "", T("Column to sort by. Options are: id,label,fingerprint,note"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *KeyListCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	keys, err := cmd.SecurityManager.ListSSHKeys("")
	if err != nil {
		return errors.NewAPIError(T("Failed to list SSH keys on your account.\n"), err.Error(), 2)
	}

	sortby := cmd.SortBy
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

	table := cmd.UI.Table([]string{T("ID"), T("label"), T("fingerprint"), T("note")})
	for _, k := range keys {
		table.Add(utils.FormatIntPointer(k.Id),
			utils.FormatStringPointer(k.Label),
			utils.FormatStringPointer(k.Fingerprint),
			utils.FormatStringPointer(k.Notes))
	}
	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

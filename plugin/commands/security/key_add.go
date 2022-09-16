package security

import (
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type KeyAddCommand struct {
	*metadata.SoftlayerCommand
	SecurityManager managers.SecurityManager
	Command         *cobra.Command
	File            string
	Key             string
	Note            string
}

func NewKeyAddCommand(sl *metadata.SoftlayerCommand) *KeyAddCommand {
	thisCmd := &KeyAddCommand{
		SoftlayerCommand: sl,
		SecurityManager:  managers.NewSecurityManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "sshkey-add " + T("LABEL"),
		Short: T("Add a new SSH key"),
		Long: T(`${COMMAND_NAME} sl security sshkey-add LABEL [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl security sshkey-add my_sshkey -f ~/.ssh/id_rsa.pub --note mykey
   This command adds an SSH key from file ~/.ssh/id_rsa.pub with a note "mykey".`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.File, "in-file", "f", "", T("The id_rsa.pub file to import for this key"))
	cobraCmd.Flags().StringVarP(&thisCmd.Key, "key", "k", "", T("The actual SSH key"))
	cobraCmd.Flags().StringVar(&thisCmd.Note, "note", "", T("Extra note to be associated with the key"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *KeyAddCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	if cmd.Key == "" && cmd.File == "" {
		return errors.NewInvalidUsageError(T("either [-k] or [-f|--in-file] is required."))
	}
	if cmd.Key != "" && cmd.File != "" {
		return errors.NewInvalidUsageError(T("[-k] is not allowed with [-f|--in-file]."))
	}
	var keyText string
	if cmd.Key != "" {
		keyText = cmd.Key
	} else {
		keyFile := cmd.File
		contents, err := ioutil.ReadFile(keyFile) // #nosec
		if err != nil {
			return errors.NewAPIError(T("Failed to read SSH key from file: {{.File}}.\n", map[string]interface{}{"File": keyFile}), err.Error(), 2)
		}
		keyText = string(contents)
	}
	key, err := cmd.SecurityManager.AddSSHKey(keyText, args[0], cmd.Note)
	if err != nil {
		return errors.NewAPIError(T("Failed to add SSH key.\n"), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, key)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("SSH key was added: ") + utils.StringPointertoString(key.Fingerprint))
	return nil
}

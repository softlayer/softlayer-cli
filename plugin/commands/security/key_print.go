package security

import (
	"io/ioutil"
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type KeyPrintCommand struct {
	*metadata.SoftlayerCommand
	SecurityManager managers.SecurityManager
	Command         *cobra.Command
	OutFile         string
}

func NewKeyPrintCommand(sl *metadata.SoftlayerCommand) *KeyPrintCommand {
	thisCmd := &KeyPrintCommand{
		SoftlayerCommand: sl,
		SecurityManager:  managers.NewSecurityManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "sshkey-print",
		Short: T("Prints out an SSH key to the screen"),
		Long: T(`${COMMAND_NAME} sl security sshkey-print IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl security sshkey-print 12345678 -f ~/mykey.pub
   This command shows the ID, label and notes of SSH key with ID 12345678 and write the public key to file: ~/mykey.pub.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.OutFile, "out-file", "f", "", T("The public SSH key will be written to this file"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *KeyPrintCommand) Run(args []string) error {
	keyID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("SSH Key ID")
	}

	outputFormat := cmd.GetOutputFlag()

	key, err := cmd.SecurityManager.GetSSHKey(keyID)
	if err != nil {
		return errors.NewAPIError(T("Failed to get SSH Key {{.KeyID}}.\n", map[string]interface{}{"KeyID": keyID}), err.Error(), 1)
	}

	if file := cmd.OutFile; file != "" {
		// #nosec G306: write on customer machine
		err = ioutil.WriteFile(file, []byte(utils.FormatStringPointer(key.Key)), 0644)
		if err != nil {
			return errors.NewAPIError(T("Failed to write SSH key to file: {{.File}}.\n",
				map[string]interface{}{"File": file}), err.Error(), 1)
		}
	}
	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(key.Id))
	table.Add(T("Label"), utils.FormatStringPointer(key.Label))
	table.Add(T("Notes"), utils.FormatStringPointer(key.Notes))
	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

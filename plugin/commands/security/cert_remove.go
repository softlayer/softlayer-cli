package security

import (
	"strconv"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CertRemoveCommand struct {
	*metadata.SoftlayerCommand
	SecurityManager managers.SecurityManager
	Command         *cobra.Command
	Force           bool
}

func NewCertRemoveCommand(sl *metadata.SoftlayerCommand) *CertRemoveCommand {
	thisCmd := &CertRemoveCommand{
		SoftlayerCommand: sl,
		SecurityManager:  managers.NewSecurityManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "cert-remove " + T("IDENTIFIER"),
		Short: T("Remove SSL certificate"),
		Long: T(`${COMMAND_NAME} sl security cert-remove IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl security cert-remove 12345678 
   This command removes certificate with ID 12345678.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CertRemoveCommand) Run(args []string) error {
	certID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("SSL certificate ID")
	}
	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will remove SSL certificate: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": certID}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.SecurityManager.RemoveCertificate(certID)
	if err != nil {
		return slErr.NewAPIError(T("Failed to remove SSL certificate: {{.ID}}.\n", map[string]interface{}{"ID": certID}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("SSL certificate {{.ID}} was removed.", map[string]interface{}{"ID": certID}))
	return nil
}

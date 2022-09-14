package security

import (
	"io/ioutil"
	"strconv"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type CertEditCommand struct {
	*metadata.SoftlayerCommand
	SecurityManager managers.SecurityManager
	Command         *cobra.Command
	Crt             string
	Csr             string
	Icc             string
	Key             string
	Notes           string
}

func NewCertEditCommand(sl *metadata.SoftlayerCommand) *CertEditCommand {
	thisCmd := &CertEditCommand{
		SoftlayerCommand: sl,
		SecurityManager:  managers.NewSecurityManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "cert-edit " + T("IDENTIFIER"),
		Short: T("Edit SSL certificate"),
		Long: T(`${COMMAND_NAME} sl security cert-edit IDENTIFIER [OPTIONS]

EXAMPLE:
	${COMMAND_NAME} sl security cert-edit 12345678 --key ~/ibm.com.key 
	This command edits certificate with ID 12345678 and updates its private key with file: ~/ibm.com.key.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Crt, "crt", "", T("Certificate file"))
	cobraCmd.Flags().StringVar(&thisCmd.Csr, "csr", "", T("Certificate Signing Request file"))
	cobraCmd.Flags().StringVar(&thisCmd.Icc, "icc", "", T("Intermediate Certificate file"))
	cobraCmd.Flags().StringVar(&thisCmd.Key, "key", "", T("Private Key file"))
	cobraCmd.Flags().StringVar(&thisCmd.Notes, "notes", "", T("Additional notes"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CertEditCommand) Run(args []string) error {
	certID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("SSL certificate ID")
	}
	template := datatypes.Security_Certificate{Id: sl.Int(certID)}
	if cmd.Crt != "" {
		crtFile := cmd.Crt
		certificate, err := ioutil.ReadFile(crtFile) // #nosec
		if err != nil {
			return errors.NewAPIError(T("Failed to read certificate file: {{.File}}.\n",
				map[string]interface{}{"File": crtFile}), err.Error(), 2)
		}
		template.Certificate = sl.String(string(certificate))
	}
	if cmd.Key != "" {
		keyFile := cmd.Key
		key, err := ioutil.ReadFile(keyFile) // #nosec
		if err != nil {
			return errors.NewAPIError(T("Failed to read private key file: {{.File}}.\n",
				map[string]interface{}{"File": keyFile}), err.Error(), 2)
		}
		template.PrivateKey = sl.String(string(key))
	}
	if cmd.Icc != "" {
		iccFile := cmd.Icc
		icc, err := ioutil.ReadFile(iccFile) // #nosec
		if err != nil {
			return errors.NewAPIError(T("Failed to read intermediate certificate file: {{.File}}.\n",
				map[string]interface{}{"File": iccFile}), err.Error(), 2)
		}
		template.IntermediateCertificate = sl.String(string(icc))
	}
	if cmd.Csr != "" {
		csrFile := cmd.Csr
		csr, err := ioutil.ReadFile(csrFile) // #nosec
		if err != nil {
			return errors.NewAPIError(T("Failed to read certificate signing request file: {{.File}}.\n",
				map[string]interface{}{"File": csrFile}), err.Error(), 2)
		}
		template.CertificateSigningRequest = sl.String(string(csr))
	}
	if cmd.Notes != "" {
		template.Notes = sl.String(cmd.Notes)
	}
	err = cmd.SecurityManager.EditCertificate(template)
	if err != nil {
		return errors.NewAPIError(T("Failed to edit SSL certificate: {{.ID}}.\n",
			map[string]interface{}{"ID": certID}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("SSL certificate {{.ID}} was updated.", map[string]interface{}{"ID": certID}))
	return nil
}

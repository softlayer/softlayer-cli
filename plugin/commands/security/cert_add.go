package security

import (
	"io/ioutil"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CertAddCommand struct {
	*metadata.SoftlayerCommand
	SecurityManager managers.SecurityManager
	Command         *cobra.Command
	Crt             string
	Csr             string
	Icc             string
	Key             string
	Notes           string
}

func NewCertAddCommand(sl *metadata.SoftlayerCommand) *CertAddCommand {

	thisCmd := &CertAddCommand{
		SoftlayerCommand: sl,
		SecurityManager:  managers.NewSecurityManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "cert-add",
		Short: T("Add and upload SSL certificate details"),
		Long: T(`${COMMAND_NAME} sl security cert-add [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl security cert-add --crt ~/ibm.com.cert --key ~/ibm.com.key 
   This command adds certificate file: ~/ibm.com.cert and private key file ~/ibm.com.key for domain ibm.com.`),
		Args: metadata.NoArgs,
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

func (cmd *CertAddCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	crtFile := cmd.Crt
	if crtFile == "" {
		return errors.NewMissingInputError("--crt")
	}
	keyFile := cmd.Key
	if keyFile == "" {
		return errors.NewMissingInputError("--key")
	}
	template := datatypes.Security_Certificate{}
	certificate, err := ioutil.ReadFile(crtFile) // #nosec
	if err != nil {
		return errors.NewAPIError(T("Failed to read certificate file: {{.File}}.\n",
			map[string]interface{}{"File": crtFile}), err.Error(), 2)
	}
	template.Certificate = sl.String(string(certificate))
	key, err := ioutil.ReadFile(keyFile) // #nosec
	if err != nil {
		return errors.NewAPIError(T("Failed to read private key file: {{.File}}.\n",
			map[string]interface{}{"File": keyFile}), err.Error(), 2)
	}
	template.PrivateKey = sl.String(string(key))
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

	cert, err := cmd.SecurityManager.AddCertificate(template)
	if err != nil {
		return errors.NewAPIError(T("Failed to add certificate.\n"), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, cert)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("SSL certificate for {{.CommonName}} was added.", map[string]interface{}{"CommonName": utils.StringPointertoString(cert.CommonName)}))
	return nil
}

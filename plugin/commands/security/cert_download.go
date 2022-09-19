package security

import (
	"io/ioutil"
	"strconv"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"

	"github.com/spf13/cobra"
	
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type CertDownloadCommand struct {
	*metadata.SoftlayerCommand
	SecurityManager managers.SecurityManager
	Command         *cobra.Command
}

func NewCertDownloadCommand(sl *metadata.SoftlayerCommand) *CertDownloadCommand {
	thisCmd := &CertDownloadCommand{
		SoftlayerCommand: sl,
		SecurityManager:  managers.NewSecurityManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "cert-download " + T("IDENTIFIER"),
		Short: T("Download SSL certificate and key files"),
		Long: T(`${COMMAND_NAME} sl security cert-download IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl security cert-download 12345678
   This command downloads four files to current directory for certificate with ID 12345678. The four files are: certificate file, certificate signing request file, intermediate certificate file and private key file.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CertDownloadCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	certID, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("SSL certificate ID")
	}
	cert, err := cmd.SecurityManager.GetCertificate(certID)
	if err != nil {
		return errors.NewAPIError(T("Failed to get SSL certificate: {{.ID}}.\n",
			map[string]interface{}{"ID": certID}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, cert)
	}

	var multiErrors []error
	if cert.Certificate == nil || cert.CommonName == nil {
		newError := errors.New(T("Certificate not found"))
		multiErrors = append(multiErrors, newError)
	} else {
		// #nosec G306: write on customer machine
		err = ioutil.WriteFile(*cert.CommonName+".crt", []byte(*cert.Certificate), 0644)
		if err != nil {
			newError := errors.New(T("Failed to write certificate to file: {{.File}}.\n",
				map[string]interface{}{"File": *cert.CommonName + ".crt"}) + err.Error())
			multiErrors = append(multiErrors, newError)
		}
	}

	if cert.PrivateKey == nil || cert.CommonName == nil {
		newError := errors.New(T("Private key not found"))
		multiErrors = append(multiErrors, newError)
	} else {
		// #nosec G306: write on customer machine
		err = ioutil.WriteFile(*cert.CommonName+".key", []byte(*cert.PrivateKey), 0644)
		if err != nil {
			newError := errors.New(T("Failed to write private key to file: {{.File}}.\n",
				map[string]interface{}{"File": *cert.CommonName + ".key"}) + err.Error())
			multiErrors = append(multiErrors, newError)
		}
	}

	if cert.IntermediateCertificate == nil || cert.CommonName == nil {
		newError := errors.New(T("intermediate certificate not found"))
		multiErrors = append(multiErrors, newError)
	} else {
		// #nosec G306: write on customer machine
		err = ioutil.WriteFile(*cert.CommonName+".icc", []byte(*cert.IntermediateCertificate), 0644)
		if err != nil {
			newError := errors.New(T("Failed to write intermediate certificate to file: {{.File}}.\n",
				map[string]interface{}{"File": *cert.CommonName + ".icc"}) + err.Error())
			multiErrors = append(multiErrors, newError)
		}
	}

	if cert.CertificateSigningRequest == nil || cert.CommonName == nil {
		newError := errors.New(T("Certificate signing request not found"))
		multiErrors = append(multiErrors, newError)
	} else {
		// #nosec G306: write on customer machine
		err = ioutil.WriteFile(*cert.CommonName+".csr", []byte(*cert.CertificateSigningRequest), 0644)
		if err != nil {
			newError := errors.New(T("Failed to write certificate signing request to file: {{.File}}.\n",
				map[string]interface{}{"File": *cert.CommonName + ".csr"}) + err.Error())
			multiErrors = append(multiErrors, newError)
		}
	}
	if len(multiErrors) > 0 {
		return errors.CollapseErrors(multiErrors)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("SSL certificate files are downloaded."))
	return nil
}

package security

import (
	"errors"
	"io/ioutil"
	"strconv"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"

	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type CertDownloadCommand struct {
	UI              terminal.UI
	SecurityManager managers.SecurityManager
}

func NewCertDownloadCommand(ui terminal.UI, securityManager managers.SecurityManager) (cmd *CertDownloadCommand) {
	return &CertDownloadCommand{
		UI:              ui,
		SecurityManager: securityManager,
	}
}

func (cmd *CertDownloadCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	certID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("SSL certificate ID")
	}
	cert, err := cmd.SecurityManager.GetCertificate(certID)
	if err != nil {
		return cli.NewExitError(T("Failed to get SSL certificate: {{.ID}}.\n",
			map[string]interface{}{"ID": certID})+err.Error(), 2)
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
		return cli.NewExitError(cli.NewMultiError(multiErrors...).Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("SSL certificate files are downloaded."))
	return nil
}

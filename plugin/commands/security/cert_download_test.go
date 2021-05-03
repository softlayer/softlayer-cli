package security_test

import (
	"errors"
	"os"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/security"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Certificate download", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeSecurityManager *testhelpers.FakeSecurityManager
		cmd                 *security.CertDownloadCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSecurityManager = new(testhelpers.FakeSecurityManager)
		cmd = security.NewCertDownloadCommand(fakeUI, fakeSecurityManager)
		cliCommand = cli.Command{
			Name:        metadata.SecuritySSLCertDownloadMetaData().Name,
			Description: metadata.SecuritySSLCertDownloadMetaData().Description,
			Usage:       metadata.SecuritySSLCertDownloadMetaData().Usage,
			Flags:       metadata.SecuritySSLCertDownloadMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Certificate download", func() {
		Context("Certificate download without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Certificate download with wrong cert ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'SSL certificate ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("Certificate download but server API call fails", func() {
			BeforeEach(func() {
				fakeSecurityManager.GetCertificateReturns(datatypes.Security_Certificate{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())

				Expect(strings.Contains(err.Error(), "Failed to get SSL certificate: 1234")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
		Context("Certificate download", func() {
			BeforeEach(func() {
				fakeSecurityManager.GetCertificateReturns(datatypes.Security_Certificate{
					Id:                        sl.Int(1234),
					CommonName:                sl.String("wilma.org"),
					Certificate:               sl.String("certificate"),
					IntermediateCertificate:   sl.String("intermediatecertificate"),
					PrivateKey:                sl.String("ssh-rsa djghtbtmfhgentongwfrdnglkhsdye"),
					CertificateSigningRequest: sl.String("CertificateSigningRequest"),
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"SSL certificate files are downloaded."}))
			})
		})
		AfterEach(func() {
			os.Remove("wilma.org.crt")
			os.Remove("wilma.org.csr")
			os.Remove("wilma.org.icc")
			os.Remove("wilma.org.key")
		})
	})
})

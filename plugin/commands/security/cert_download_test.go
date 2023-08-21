package security_test

import (
	"errors"
	"os"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/security"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Certificate download", func() {
	var (
		fakeUI              *terminal.FakeUI
		cliCommand          *security.CertDownloadCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
		fakeSecurityManager *testhelpers.FakeSecurityManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = security.NewCertDownloadCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeSecurityManager = new(testhelpers.FakeSecurityManager)
		cliCommand.SecurityManager = fakeSecurityManager
	})

	Describe("Certificate download", func() {
		Context("Certificate download without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("Certificate download with wrong cert ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'SSL certificate ID'. It must be a positive integer."))
			})
		})
		Context("Certificate download but server API call fails", func() {
			BeforeEach(func() {
				fakeSecurityManager.GetCertificateReturns(datatypes.Security_Certificate{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get SSL certificate: 1234"))
			})
		})
		Context("Certificate download is malformed", func() {
			BeforeEach(func() {
				fakeSecurityManager.GetCertificateReturns(datatypes.Security_Certificate{}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Certificate not found"))
				Expect(err.Error()).To(ContainSubstring("Private key not found"))
				Expect(err.Error()).To(ContainSubstring("intermediate certificate not found"))
				Expect(err.Error()).To(ContainSubstring("Certificate signing request not found"))
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
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SSL certificate files are downloaded."))
			})
			It("Output JSON", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--output", "JSON")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("\"commonName\": \"wilma.org\""))
			})
			It("Handle unable to write file", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--output", "JSON")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("\"commonName\": \"wilma.org\""))
			})
		})
		Context("Certificate download unable to write file", func() {
			BeforeEach(func() {
				fakeSecurityManager.GetCertificateReturns(datatypes.Security_Certificate{
					Id:                        sl.Int(1234),
					CommonName:                sl.String("/path/to/nothing/wilma.org"),
					Certificate:               sl.String("certificate"),
					IntermediateCertificate:   sl.String("intermediatecertificate"),
					PrivateKey:                sl.String("ssh-rsa djghtbtmfhgentongwfrdnglkhsdye"),
					CertificateSigningRequest: sl.String("CertificateSigningRequest"),
				}, nil)
			})
			It("Handle unable to write file", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to write certificate to file"))
				Expect(err.Error()).To(ContainSubstring("Failed to write private key to file"))
				Expect(err.Error()).To(ContainSubstring("Failed to write intermediate certificate to file"))
				Expect(err.Error()).To(ContainSubstring("Failed to write certificate signing request to file"))
			})
		})
		AfterEach(func() {
			os.Remove("wilma.org.crt")
			os.Remove("wilma.org.csr")
			os.Remove("wilma.org.icc")
			os.Remove("wilma.org.key")
		})
		Context("Check bad output format", func() {
			It("Error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--output", "text")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format"))
			})
		})
	})
})

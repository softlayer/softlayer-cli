package security_test

import (
	"errors"
	"os"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/security"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Certificate add", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeSecurityManager *testhelpers.FakeSecurityManager
		cmd                 *security.CertAddCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSecurityManager = new(testhelpers.FakeSecurityManager)
		cmd = security.NewCertAddCommand(fakeUI, fakeSecurityManager)

		cliCommand = cli.Command{
			Name:        security.SecuritySSLCertAddMetaData().Name,
			Description: security.SecuritySSLCertAddMetaData().Description,
			Usage:       security.SecuritySSLCertAddMetaData().Usage,
			Flags:       security.SecuritySSLCertAddMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Certificate add", func() {
		Context("Certificate add without crt", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--crt' is required"))
			})
		})
		Context("Certificate add without key", func() {
			It("return error", func() {
				file, _ := os.Create("wilma.org.crt")
				err := testhelpers.RunCommand(cliCommand, "--crt", file.Name())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--key' is required"))
			})
		})
		Context("Certificate add with server fails", func() {
			BeforeEach(func() {
				fakeSecurityManager.AddCertificateReturns(datatypes.Security_Certificate{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				crtFile, _ := os.Create("wilma.org.crt")
				keyFile, _ := os.Create("wilma.org.key")
				err := testhelpers.RunCommand(cliCommand, "--crt", crtFile.Name(), "--key", keyFile.Name())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to add certificate."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})
		Context("Certificate added ", func() {
			crtFile, _ := os.Create(os.TempDir() + "/wilma1.org.crt")
			keyFile, _ := os.Create(os.TempDir() + "/wilma1.org.key")
			BeforeEach(func() {
				fakeSecurityManager.AddCertificateReturns(datatypes.Security_Certificate{
					Id:         sl.Int(1234),
					CommonName: sl.String("wilma.org"),
				}, nil)

			})
			It("Success minimum options", func() {
				err := testhelpers.RunCommand(cliCommand, "--crt", crtFile.Name(), "--key", keyFile.Name())
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SSL certificate for wilma.org was added."))
			})
			It("Success all options", func() {
				err := testhelpers.RunCommand(cliCommand, "--crt", crtFile.Name(), "--key", keyFile.Name(), "--icc",
					keyFile.Name(), "--csr", keyFile.Name(), "--notes", "testNotes")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SSL certificate for wilma.org was added."))
				argsForCall := fakeSecurityManager.AddCertificateArgsForCall(0)
				Expect(*argsForCall.Notes).To(Equal("testNotes"))
			})
			It("Success JSON output", func() {
				err := testhelpers.RunCommand(cliCommand, "--crt", crtFile.Name(), "--key", keyFile.Name(), "--output", "JSON")
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeUI.Outputs()).To(ContainSubstring("\"commonName\": \"wilma.org\""))
			})
			It("Handle Bad file CRT", func() {
				err := testhelpers.RunCommand(cliCommand, "--crt", "./1fakeFile", "--key", keyFile.Name())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to read certificate file"))
			})
			It("Handle Bad file KEY", func() {
				err := testhelpers.RunCommand(cliCommand, "--crt", crtFile.Name(), "--key", "fakeKeyFile")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to read private key file"))
			})
			It("Handle Bad file ICC", func() {
				err := testhelpers.RunCommand(cliCommand, "--crt", crtFile.Name(), "--key", keyFile.Name(), "--icc", "fakeFile")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to read intermediate certificate file"))
			})
			It("Handle Bad file CSR", func() {
				err := testhelpers.RunCommand(cliCommand, "--crt", crtFile.Name(), "--key", keyFile.Name(), "--csr", "fakeFile")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to read certificate signing request file"))
			})
		})
		AfterEach(func() {
			os.Remove(os.TempDir() + "wilma.org.crt")
			os.Remove(os.TempDir() + "wilma.org.key")
		})
		Context("Check bad output format", func() {
			It("Error", func() {
				err := testhelpers.RunCommand(cliCommand, "--output", "text")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format"))
			})
		})
	})
})

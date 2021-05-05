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
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/security"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
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
			Name:        metadata.SecuritySSLCertAddMetaData().Name,
			Description: metadata.SecuritySSLCertAddMetaData().Description,
			Usage:       metadata.SecuritySSLCertAddMetaData().Usage,
			Flags:       metadata.SecuritySSLCertAddMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Certificate add", func() {
		Context("Certificate add without crt", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: '--crt' is required")).To(BeTrue())
			})
		})
		Context("Certificate add without key", func() {
			It("return error", func() {
				file, _ := os.Create("wilma.org.crt")
				err := testhelpers.RunCommand(cliCommand, "--crt", file.Name())
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: '--key' is required")).To(BeTrue())
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
				Expect(strings.Contains(err.Error(), "Failed to add certificate.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
		Context("Certificate added ", func() {
			BeforeEach(func() {
				fakeSecurityManager.AddCertificateReturns(datatypes.Security_Certificate{
					Id:         sl.Int(1234),
					CommonName: sl.String("wilma.org"),
				}, nil)
			})
			It("return error", func() {
				crtFile, _ := os.Create("wilma.org.crt")
				keyFile, _ := os.Create("wilma.org.key")
				err := testhelpers.RunCommand(cliCommand, "--crt", crtFile.Name(), "--key", keyFile.Name())
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"SSL certificate for wilma.org was added."}))
			})
		})
		AfterEach(func() {
			os.Remove("wilma.org.crt")
			os.Remove("wilma.org.key")
		})
	})
})

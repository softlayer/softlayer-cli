package security_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/security"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Certificate edit", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeSecurityManager *testhelpers.FakeSecurityManager
		cmd                 *security.CertEditCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSecurityManager = new(testhelpers.FakeSecurityManager)
		cmd = security.NewCertEditCommand(fakeUI, fakeSecurityManager)
		cliCommand = cli.Command{
			Name:        metadata.SecuritySSLCertEdit().Name,
			Description: metadata.SecuritySSLCertEdit().Description,
			Usage:       metadata.SecuritySSLCertEdit().Usage,
			Flags:       metadata.SecuritySSLCertEdit().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Certificate edit", func() {
		Context("Certificate edit without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Certificate edit without cert ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Certificate add with wrong cert ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'SSL certificate ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("Certificate add with server fails", func() {
			BeforeEach(func() {
				fakeSecurityManager.EditCertificateReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to edit SSL certificate: 1234")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
		Context("Certificate edited ", func() {
			BeforeEach(func() {
				fakeSecurityManager.EditCertificateReturns(nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"SSL certificate 1234 was updated."}))
			})
		})
	})
})

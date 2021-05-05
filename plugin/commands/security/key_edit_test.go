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

var _ = Describe("Key edit", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeSecurityManager *testhelpers.FakeSecurityManager
		cmd                 *security.KeyEditCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSecurityManager = new(testhelpers.FakeSecurityManager)
		cmd = security.NewKeyEditCommand(fakeUI, fakeSecurityManager)
		cliCommand = cli.Command{
			Name:        metadata.SecuritySSHKeyEditMetaData().Name,
			Description: metadata.SecuritySSHKeyEditMetaData().Description,
			Usage:       metadata.SecuritySSHKeyEditMetaData().Usage,
			Flags:       metadata.SecuritySSHKeyEditMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Key edit", func() {
		Context("Key edit without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Key edit with wrong key ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'SSH Key ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("Key edit with no label and no note", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: either [--label] or [--note] must be specified to edit SSH key.")).To(BeTrue())
			})
		})
		Context("Key edit but server API call fails", func() {
			BeforeEach(func() {
				fakeSecurityManager.EditSSHKeyReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--label", "newlabel")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to edit SSH key: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
		Context("Key edit with different parameters", func() {
			BeforeEach(func() {
				fakeSecurityManager.EditSSHKeyReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--label", "newlabel")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"SSH key 1234 was updated."}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--note", "newnote")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"SSH key 1234 was updated."}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--label", "newlabel", "--note", "newnote")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"SSH key 1234 was updated."}))
			})
		})
	})
})

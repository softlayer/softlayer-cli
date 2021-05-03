package security_test

import (
	"errors"
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

var _ = Describe("Key print", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeSecurityManager *testhelpers.FakeSecurityManager
		cmd                 *security.KeyPrintCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSecurityManager = new(testhelpers.FakeSecurityManager)
		cmd = security.NewKeyPrintCommand(fakeUI, fakeSecurityManager)
		cliCommand = cli.Command{
			Name:        metadata.SecuritySSHKeyPrintMetaData().Name,
			Description: metadata.SecuritySSHKeyPrintMetaData().Description,
			Usage:       metadata.SecuritySSHKeyPrintMetaData().Usage,
			Flags:       metadata.SecuritySSHKeyPrintMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Key print", func() {
		Context("Key print without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Key print with wrong key ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'SSH Key ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("Key print but server API call fails", func() {
			BeforeEach(func() {
				fakeSecurityManager.GetSSHKeyReturns(datatypes.Security_Ssh_Key{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get SSH Key 1234")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
		Context("Key print", func() {
			BeforeEach(func() {
				fakeSecurityManager.GetSSHKeyReturns(datatypes.Security_Ssh_Key{
					Id:    sl.Int(1234),
					Label: sl.String("label"),
					Notes: sl.String("notes"),
					Key:   sl.String("ssh-rsa djghtbtmfhgentongwfrdnglkhsdye"),
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"label"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"notes"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"ssh-rsa djghtbtmfhgentongwfrdnglkhsdye"}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f", "/tmp/key")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"label"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"notes"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"ssh-rsa djghtbtmfhgentongwfrdnglkhsdye"}))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f", "/root/key")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to write SSH key to file: /root/key.")).To(BeTrue())
			})
		})
	})
})

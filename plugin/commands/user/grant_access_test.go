package user_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Grant Access", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeUserManager *testhelpers.FakeUserManager
		cmd             *user.GrantAccessCommand
		cliCommand      cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeUserManager = new(testhelpers.FakeUserManager)
		cmd = user.NewGrantAccessCommand(fakeUI, fakeUserManager)
		cliCommand = cli.Command{
			Name:        user.UserGrantAccessMataData().Name,
			Description: user.UserGrantAccessMataData().Description,
			Usage:       user.UserGrantAccessMataData().Usage,
			Flags:       user.UserGrantAccessMataData().Flags,
			Action:      cmd.Run,
		}
	})
	Describe("user grant-access", func() {
		Context("Return error", func() {
			It("Set command without identifier", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one identifier."))
			})

			It("Set command with an invalid identifier", func() {
				err := testhelpers.RunCommand(cliCommand, "abcd", "--hardware=123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: User ID should be a number."))
			})

			It("Set command without options", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one option."))
			})

			It("Set hardware option with an invalid value", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--hardware=abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Hardware ID should be a number."))
			})

			It("Set virtual option with an invalid value", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--virtual=abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Virtual server ID should be a number."))
			})

			It("Set dedicated option with an invalid value", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--dedicated=abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Dedicated host ID should be a number."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeUserManager.AddHardwareAccessReturns(true, nil)
				fakeUserManager.AddDedicatedHostAccessReturns(true, nil)
				fakeUserManager.AddVirtualGuestAccessReturns(true, nil)
			})

			It("Set command with valid user and hardware", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--hardware=123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Access was granted"))
			})

			It("Set command with valid user and virtual guest", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--virtual=123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Access was granted"))
			})

			It("Set command with valid user and dedicated host", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--dedicated=123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Access was granted"))
			})
		})
	})
})

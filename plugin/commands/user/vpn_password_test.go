package user_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("user vpn-password", func() {
	var (
		fakeUI          *terminal.FakeUI
		cliCommand      *user.VpnPasswordCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
		fakeUserManager *testhelpers.FakeUserManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeUserManager = new(testhelpers.FakeUserManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = user.NewVpnPasswordCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.UserManager = fakeUserManager
	})

	Describe("user vpn-password", func() {

		Context("Return error", func() {
			It("Set command without Argument", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})

			It("Set command with an invalid user Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcde", "--password=Mypassword1.")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'User ID'. It must be a positive integer."))
			})

			It("Set without any option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "111111")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`required flag(s) "password" not set`))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeUserManager.UpdateVpnPasswordReturns(false, errors.New("Failed to update user vpn password."))
			})
			It("Failed edit user", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "111111", "--password=Mypassword1.")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to update user vpn password."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeUserManager.UpdateVpnPasswordReturns(true, nil)
			})
			It("enable", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "111111", "--password=Mypassword1.")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Successfully updated user VPN password"))
			})
		})
	})
})

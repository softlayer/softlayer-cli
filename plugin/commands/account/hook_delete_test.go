package account_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("account hook-delete", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *account.HookDeleteCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeAccountManager *testhelpers.FakeAccountManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeAccountManager = new(testhelpers.FakeAccountManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = account.NewHookDeleteCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.AccountManager = fakeAccountManager
	})

	Describe("account hook-delete", func() {

		Context("Return error", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage : This command requires one argument"))
			})

			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hook ID'. It must be a positive integer."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeAccountManager.DeleteProvisioningScriptReturns(false, errors.New("Failed to delete Provisioning Hook"))
			})
			It("Failed to delete Provisioning Hook", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to delete Provisioning Hook"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeAccountManager.DeleteProvisioningScriptReturns(true, nil)
			})
			It("Return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Successfully removed Provisioning Hook."))
			})
		})
	})
})

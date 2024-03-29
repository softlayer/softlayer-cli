package account_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Account cancel-item", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *account.CancelItemCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeAccountManager *testhelpers.FakeAccountManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeAccountManager = new(testhelpers.FakeAccountManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = account.NewCancelItemCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.AccountManager = fakeAccountManager
	})

	Describe("Account cancel-item", func() {
		Context("Account cancel-item, Invalid Usage", func() {
			It("Set command without id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
			It("Set command with id like letters", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Item ID'. It must be a positive integer."))
			})
		})

		Context("Account cancel-item, softlayer errors", func() {
			It("Set command with unknow item ID", func() {
				fakeAccountManager.CancelItemReturns(errors.New("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123'. (HTTP 404)"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to find item with ID: 123."))
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123'. (HTTP 404)"))
			})
			It("Set command with used ID", func() {
				fakeAccountManager.CancelItemReturns(errors.New("SoftLayer_Exception_Public: This cancellation could not be processed please contact support. This billing item is already canceled. (HTTP 500)"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to cancel item: 123."))
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_Public: This cancellation could not be processed please contact support. This billing item is already canceled. (HTTP 500)"))
			})
		})

		Context("Account cancel-item, correct use", func() {
			It("return account cancel-item", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Item: 123 was cancelled."))
			})
		})
	})
})

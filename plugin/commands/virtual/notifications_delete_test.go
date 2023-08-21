package virtual_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("vs notifications-delete", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cliCommand    *virtual.NotificationsDeleteCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewNotificationsDeleteCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
	})

	Describe("vs notifications-delete", func() {
		Context("vs notifications-delete without arguments", func() {
			It("vs notifications-delete without virtual server ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})

		Context("vs notifications-delete with wrong virtual server ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Virtual server ID'. It must be a positive integer."))
			})
		})
		Context("failed to delete vs notification", func() {
			BeforeEach(func() {
				fakeVSManager.DeleteUserCustomerNotificationReturns(false, errors.New("Failed to delete User Customer Notification"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Errors()).To(ContainSubstring("Failed to delete User Customer Notification."))
			})
		})

		Context("vs notifications-delete with correct vs ID ", func() {
			BeforeEach(func() {
				fakeVSManager.DeleteUserCustomerNotificationReturns(true, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Successfully removed User Customer notification."))
			})
		})
	})
})

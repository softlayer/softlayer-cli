package virtual_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS notifications", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *virtual.NotifiactionsCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeVSManager *testhelpers.FakeVirtualServerManager
	)

	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewNotifiactionsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
	})

	Describe("vs notifications", func() {
		Context("vs notifications without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("vs notifications with wrong virtual server ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Virtual server ID'. It must be a positive integer."))
			})
		})

		Context("failed to get vs notifications", func() {
			BeforeEach(func() {
				fakeVSManager.GetUserCustomerNotificationsByVirtualGuestIdReturns([]datatypes.User_Customer_Notification_Virtual_Guest{}, errors.New("Failed to get User Customer Notifications."))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get User Customer Notifications."))
			})
		})

		Context("vs notifications with correct virtual guest ID ", func() {
			BeforeEach(func() {
				fakerUserCustomerNotifications := []datatypes.User_Customer_Notification_Virtual_Guest{
					datatypes.User_Customer_Notification_Virtual_Guest{
						User: &datatypes.User_Customer{
							LastName:  sl.String("Jhon"),
							FirstName: sl.String("Smith"),
							Email:     sl.String("jhonsmith@email.com"),
							Username:  sl.String("jhonsmith"),
						},
					},
				}
				fakeVSManager.GetUserCustomerNotificationsByVirtualGuestIdReturns(fakerUserCustomerNotifications, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Jhon"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Smith"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("jhonsmith@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("jhonsmith"))
			})
		})
	})
})

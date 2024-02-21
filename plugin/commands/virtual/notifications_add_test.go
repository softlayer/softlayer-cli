package virtual_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("vs notifications-add", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cliCommand    *virtual.NotificationsAddCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewNotificationsAddCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
	})

	Describe("vs notifications-add", func() {
		Context("vs notifications-add without arguments", func() {
			It("vs notifications-add without virtual server ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})

		Context("vs notifications-add with wrong virtual server ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Virtual server ID'. It must be a positive integer."))
			})
		})

		Context("failed to add vs notification", func() {
			BeforeEach(func() {
				fakeVSManager.CreateUserCustomerNotificationReturns(datatypes.User_Customer_Notification_Virtual_Guest{}, errors.New("Failed to create User Customer Notification"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--users=111111")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Errors()).To(ContainSubstring("Failed to create User Customer Notification"))
			})
		})

		Context("hardware notifications-add with correct hardware ID ", func() {
			BeforeEach(func() {
				fakerUserCustomerNotification := datatypes.User_Customer_Notification_Virtual_Guest{

					Id: sl.Int(123456),
					Guest: &datatypes.Virtual_Guest{
						FullyQualifiedDomainName: sl.String("virtualserver@domain.com"),
					},
					User: &datatypes.User_Customer{
						LastName:  sl.String("Jhon"),
						FirstName: sl.String("Smith"),
						Email:     sl.String("jhonsmith@email.com"),
						Username:  sl.String("jhonsmith"),
					},
				}
				fakeVSManager.CreateUserCustomerNotificationReturns(fakerUserCustomerNotification, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--users=111111")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtualserver@domain.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Jhon"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Smith"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("jhonsmith@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("jhonsmith"))
			})
		})
	})
})

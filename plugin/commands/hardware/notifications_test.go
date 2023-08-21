package hardware_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware notifications", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cliCommand          *hardware.NotificationsCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewNotificationsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.HardwareManager = fakeHardwareManager
	})

	Describe("hardware notifications", func() {
		Context("hardware notifications without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("hardware notifications with wrong hardware ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'. It must be a positive integer."))
			})
		})

		Context("failed to get hardware notifications", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetUserCustomerNotificationsByHardwareIdReturns([]datatypes.User_Customer_Notification_Hardware{}, errors.New("Failed to get User Customer Notifications."))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get User Customer Notifications."))
			})
		})

		Context("hardware notifications with correct hardware ID ", func() {
			BeforeEach(func() {
				fakerUserCustomerNotifications := []datatypes.User_Customer_Notification_Hardware{
					datatypes.User_Customer_Notification_Hardware{
						User: &datatypes.User_Customer{
							LastName:  sl.String("Jhon"),
							FirstName: sl.String("Smith"),
							Email:     sl.String("jhonsmith@email.com"),
							Username:  sl.String("jhonsmith"),
						},
						Id: sl.Int(111111),
					},
				}
				fakeHardwareManager.GetUserCustomerNotificationsByHardwareIdReturns(fakerUserCustomerNotifications, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Jhon"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Smith"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("jhonsmith@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("jhonsmith"))
			})
		})
	})
})

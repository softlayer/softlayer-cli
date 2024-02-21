package hardware_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware notifications-add", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cliCommand          *hardware.NotificationsAddCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewNotificationsAddCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.HardwareManager = fakeHardwareManager
	})

	Describe("hardware notifications-add", func() {
		Context("hardware notifications-add without arguments", func() {
			It("hardware notifications-add without hardware ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})

		Context("hardware notifications-add with wrong hardware ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'. It must be a positive integer."))
			})
		})

		Context("failed to add hardware notification", func() {
			BeforeEach(func() {
				fakeHardwareManager.CreateUserCustomerNotificationReturns(datatypes.User_Customer_Notification_Hardware{}, errors.New("Failed to create User Customer Notification"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--users=111111")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Errors()).To(ContainSubstring("Failed to create User Customer Notification"))
			})
		})

		Context("hardware notifications-add with correct hardware ID ", func() {
			BeforeEach(func() {
				fakerUserCustomerNotification := datatypes.User_Customer_Notification_Hardware{

					Id: sl.Int(123456),
					Hardware: &datatypes.Hardware{
						FullyQualifiedDomainName: sl.String("hardware@domain.com"),
					},
					User: &datatypes.User_Customer{
						LastName:  sl.String("Jhon"),
						FirstName: sl.String("Smith"),
						Email:     sl.String("jhonsmith@email.com"),
						Username:  sl.String("jhonsmith"),
					},
				}
				fakeHardwareManager.CreateUserCustomerNotificationReturns(fakerUserCustomerNotification, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--users=111111")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("jhonsmith@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Jhon"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Smith"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("jhonsmith@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("jhonsmith"))
			})
		})
	})
})

package user_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("User Notifications", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeUserManager *testhelpers.FakeUserManager
		cliCommand      *user.NotificationsCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeUserManager = new(testhelpers.FakeUserManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = user.NewNotificationsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.UserManager = fakeUserManager
	})

	Describe("user list ", func() {

		Context("Return error", func() {
			BeforeEach(func() {
				fakeUserManager.GetAllNotificationsReturns([]datatypes.Email_Subscription{}, errors.New("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
			It("An invalid output id is set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return no error", func() {
			fakeNotifications := []datatypes.Email_Subscription{}
			BeforeEach(func() {
				fakeNotifications = []datatypes.Email_Subscription{
					datatypes.Email_Subscription{
						Id:          sl.Int(1),
						Name:        sl.String("Order Being Reviewed"),
						Description: sl.String("Email about your order."),
						Enabled:     sl.Bool(true),
					},
					datatypes.Email_Subscription{
						Id:          sl.Int(12),
						Name:        sl.String("Severity 2"),
						Description: sl.String("Incidents that cause measurable service degradation yet not an actual outage."),
						Enabled:     sl.Bool(false),
					},
				}
				fakeUserManager.GetAllNotificationsReturns(fakeNotifications, nil)
			})

			It("list notifications", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order Being Reviewed"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email about your order."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("12"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Severity 2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Incidents that cause measurable service degradation yet not an actual outage."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("false"))
			})

			It("list notifications in json format", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"id": 1`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"name": "Order Being Reviewed"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"description": "Email about your order."`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"enabled": true`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"id": 12`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"name": "Severity 2"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"description": "Incidents that cause measurable service degradation yet not an actual outage."`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"enabled": false`))
			})
		})
	})
})

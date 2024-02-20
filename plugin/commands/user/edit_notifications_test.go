package user_test

import (
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

var _ = Describe("User Edit Notifications", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeUserManager *testhelpers.FakeUserManager
		cliCommand      *user.EditNotificationsCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeUserManager = new(testhelpers.FakeUserManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = user.NewEditNotificationsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.UserManager = fakeUserManager
	})

	Describe("user Edit Notifications", func() {

		Context("Return error", func() {
			It("An invalid output id is set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})

			It("Set --enable and --disable options", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--disable='Order Being Reviewed'", "--enable='High Impact'")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Only set --enable or --disable options."))
			})

			It("Set --enable and --disable options", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--disable='Order Being Reviewed'", "--enable='High Impact'")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Only set --enable or --disable options."))
			})

			It("Set --enable and --disable options", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "'Order Being Reviewed'", "--disable='High Impact'")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Only set --enable or --disable options."))
			})

			It("No set options and arguments", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires notification names as arguments and options flags."))
			})
		})

		Context("Return error", func() {
			fakeNotifications := []datatypes.Email_Subscription{}
			BeforeEach(func() {
				fakeNotifications = []datatypes.Email_Subscription{
					datatypes.Email_Subscription{
						Id:          sl.Int(1),
						Name:        sl.String("Order Being Reviewed"),
						Description: sl.String("Email about your order."),
						Enabled:     sl.Bool(false),
					},
					datatypes.Email_Subscription{
						Id:          sl.Int(12),
						Name:        sl.String("Severity 2"),
						Description: sl.String("Incidents that cause measurable service degradation yet not an actual outage."),
						Enabled:     sl.Bool(false),
					},
				}
				fakeUserManager.GetAllNotificationsReturns(fakeNotifications, nil)
				fakeUserManager.EnableEmailSubscriptionNotificationReturns(false, nil)
			})

			It("Enable notification that does not exist", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--enable='Order Email'")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Notifications updated unsuccessfully: 'Order Email'. Review if already set or if the name is correct."))
			})

			It("Enable notification that does not exist in json output format", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--enable='Order Email'", "--output=json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("false"))
			})

			It("Disable notification that does not exist in json output format", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--disable='OrderEmail'", "--output=json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("false"))
			})

			It("Disable notification that does not exist", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--disable='Order Email'")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Notifications updated unsuccessfully: 'Order Email'. Review if already set or if the name is correct."))
			})
		})
	})
})

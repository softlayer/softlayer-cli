package account_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Account list Events", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *account.EventsCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = account.NewEventsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("Account events", func() {
		Context("Account events, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
			It("Set command with an invalid date option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--date-min", "abcd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid format date."))
			})
		})

		Context("Account events, correct use", func() {
			It("return account events", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--date-min", "2022-03-12")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Planned"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Event Data             Id       Event ID    Subject                Status   Items   Start Date             End Date               Acknowledged   Updates"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2022-04-08T00:30:00Z   341058   144369902   Maintenance - Zone 2   Active   2       2022-04-08T00:30:00Z   2022-04-08T06:00:00Z   false          1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Unplanned"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Id       Event ID    Subject                Status   Items   Start Date             Last Updated           Acknowledged   Updates"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("341058   144369902   Maintenance - Zone 2   Active   2       2022-04-08T00:30:00Z   2022-03-24T17:34:32Z   false          1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Announcement"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Id       Event ID    Subject                Status   Items   Acknowledged   Updates"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("341058   144369902   Maintenance - Zone 2   Active   2       false          1"))
			})
			It("return account events in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`Planned`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Event Data": "2022-04-08T00:30:00Z",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`Unplanned`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Id": "341058",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`Announcement`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Id": "341058",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})

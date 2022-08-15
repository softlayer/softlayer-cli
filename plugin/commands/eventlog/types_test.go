package eventlog_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/eventlog"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("event-log types", func() {
	var (
		fakeUI              *terminal.FakeUI
		cliCommand          *eventlog.TypesCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
		fakeEventLogManager *testhelpers.FakeEventLogManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeEventLogManager = new(testhelpers.FakeEventLogManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = eventlog.NewTypesCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.EventLogManager = fakeEventLogManager
	})

	Describe("event-log types", func() {

		Context("Return error", func() {
			BeforeEach(func() {
				fakeEventLogManager.GetEventLogTypesReturns([]string{}, errors.New("Failed to get Event Logs types"))
			})
			It("Failed get Event Logs types", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get Event Logs types"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerTypes := []string{
					"API Authentication",
					"User",
					"Bluemix LB",
				}
				fakeEventLogManager.GetEventLogTypesReturns(fakerTypes, nil)
			})

			It("Set command with only --date-min", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("API Authentication"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("User"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Bluemix LB"))
			})
		})
	})
})

package eventlog_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/eventlog"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("event-log get", func() {
	var (
		fakeUI              *terminal.FakeUI
		cliCommand          *eventlog.GetCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
		fakeEventLogManager *testhelpers.FakeEventLogManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeEventLogManager = new(testhelpers.FakeEventLogManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = eventlog.NewGetCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.EventLogManager = fakeEventLogManager
	})

	Describe("event-log get", func() {

		Context("Return error", func() {

			It("Set invalid --date-min value", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--date-min=05/10/2022")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid format date to --date-min."))
			})

			It("Set invalid --date-max value", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--date-max=05/10/2022")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid format date to --date-max."))
			})
		})
		Context("Return error", func() {
			BeforeEach(func() {
				fakeEventLogManager.GetEventLogsReturns([]datatypes.Event_Log{}, errors.New("Failed to get Event Logs."))
			})
			It("Failed get event logs", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--limit=10")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get Event Logs."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2017-01-01T00:00:00Z")
				fakerLogs := []datatypes.Event_Log{
					datatypes.Event_Log{
						EventName:       sl.String("IAM Token validation successful"),
						Label:           sl.String("sl307608-chechu"),
						ObjectName:      sl.String("User"),
						EventCreateDate: sl.Time(created),
						UserId:          sl.Int(123456),
						User: &datatypes.User_Customer{
							Username: sl.String("user1234"),
						},
						MetaData: sl.String("metadata"),
					},
				}
				fakeEventLogManager.GetEventLogsReturns(fakerLogs, nil)
			})

			It("Set command with all options", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--date-min=2016-01-01", "--date-max=2017-02-01", "--obj-id=123456", "--obj-event=Create", "--obj-type=Create", "--metadata", "--utc-offset=-0000")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("IAM Token validation successful"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("sl307608-chechu"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("User"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-01-01T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("user1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("metadata"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2017-01-01T00:00:00Z")
				fakerLogs := []datatypes.Event_Log{
					datatypes.Event_Log{
						EventName:       sl.String("Power On"),
						Label:           sl.String("testvs-ab8s.domain.com"),
						ObjectName:      sl.String("CCI"),
						EventCreateDate: sl.Time(created),
						UserId:          nil,
						UserType:        sl.String("SYSTEM"),
					},
				}
				fakeEventLogManager.GetEventLogsReturns(fakerLogs, nil)
			})

			It("Set command with only --date-min", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--date-min=2016-01-01")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Power On"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("testvs-ab8s.domain.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("CCI"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-01-01T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SYSTEM"))
			})

			It("Set command with only --date-max", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--date-max=2016-01-01")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Power On"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("testvs-ab8s.domain.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("CCI"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-01-01T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SYSTEM"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerLogs := []datatypes.Event_Log{
					datatypes.Event_Log{
						EventName: nil,
					},
				}
				fakeEventLogManager.GetEventLogsReturns(fakerLogs, nil)
			})

			It("Set command with all options", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--date-min=2016-01-01", "--date-max=2017-02-01", "--obj-id=123456", "--obj-event=Create", "--obj-type=Create", "--metadata", "--utc-offset=-0000")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("No logs available for filter"))
			})
		})
	})
})

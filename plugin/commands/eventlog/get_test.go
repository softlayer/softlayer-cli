package eventlog_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/eventlog"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("event-log get", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeEventLogManager *testhelpers.FakeEventLogManager
		cmd                 *eventlog.GetCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeEventLogManager = new(testhelpers.FakeEventLogManager)
		cmd = eventlog.NewGetCommand(fakeUI, fakeEventLogManager)
		cliCommand = cli.Command{
			Name:        eventlog.EventLogGetMetaData().Name,
			Description: eventlog.EventLogGetMetaData().Description,
			Usage:       eventlog.EventLogGetMetaData().Usage,
			Flags:       eventlog.EventLogGetMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("event-log get", func() {

		Context("Return error", func() {
			It("Set command with invalid limit", func() {
				err := testhelpers.RunCommand(cliCommand, "--limit=abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'limit'. It must be a positive integer."))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCommand(cliCommand, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})

			It("Set invalid --date-min value", func() {
				err := testhelpers.RunCommand(cliCommand, "--date-min=05/10/2022")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid format date to --date-min."))
			})

			It("Set invalid --date-max value", func() {
				err := testhelpers.RunCommand(cliCommand, "--date-max=05/10/2022")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid format date to --date-max."))
			})
		})
		Context("Return error", func() {
			BeforeEach(func() {
				fakeEventLogManager.GetEventLogsReturns([]datatypes.Event_Log{}, errors.New("Failed to get Event Logs."))
			})
			It("Failed get event logs", func() {
				err := testhelpers.RunCommand(cliCommand, "--limit=10")
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
				err := testhelpers.RunCommand(cliCommand, "--date-min=2016-01-01", "--date-max=2017-02-01", "--obj-id=123456", "--obj-event=Create", "--obj-type=Create", "--metadata", "--utc-offset=-0000")
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
				err := testhelpers.RunCommand(cliCommand, "--date-min=2016-01-01")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Power On"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("testvs-ab8s.domain.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("CCI"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-01-01T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SYSTEM"))
			})

			It("Set command with only --date-max", func() {
				err := testhelpers.RunCommand(cliCommand, "--date-max=2016-01-01")
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
				err := testhelpers.RunCommand(cliCommand, "--date-min=2016-01-01", "--date-max=2017-02-01", "--obj-id=123456", "--obj-event=Create", "--obj-type=Create", "--metadata", "--utc-offset=-0000")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("No logs available for filter"))
			})
		})
	})
})

package eventlog_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/eventlog"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("event-log types", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeEventLogManager *testhelpers.FakeEventLogManager
		cmd                 *eventlog.TypesCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeEventLogManager = new(testhelpers.FakeEventLogManager)
		cmd = eventlog.NewTypesCommand(fakeUI, fakeEventLogManager)
		cliCommand = cli.Command{
			Name:        eventlog.EventLogTypesMetaData().Name,
			Description: eventlog.EventLogTypesMetaData().Description,
			Usage:       eventlog.EventLogTypesMetaData().Usage,
			Flags:       eventlog.EventLogTypesMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("event-log types", func() {

		Context("Return error", func() {

			It("Set invalid output", func() {
				err := testhelpers.RunCommand(cliCommand, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeEventLogManager.GetEventLogTypesReturns([]string{}, errors.New("Failed to get Event Logs types"))
			})
			It("Failed get Event Logs types", func() {
				err := testhelpers.RunCommand(cliCommand)
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
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("API Authentication"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("User"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Bluemix LB"))
			})
		})
	})
})

package block_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Snapshot enable", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *block.SnapshotEnableCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
		FakeStorageManager *testhelpers.FakeStorageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewSnapshotEnableCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("Snapshot enable", func() {
		Context("Incorrect Usage Tests", func() {
			It("No arguments", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
			It("Bad Volume ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "a1234", "-c=100", "-s=INTERVAL")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
			})
			It("Bad Interval", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--schedule-type=FAKE", "-c=100")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("needs to be one of INTERVAL, HOURLY, DAILY, WEEKLY, not FAKE."))
			})
			It("No Retention Count", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "HOURLY")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`required flag(s) "retention-count" not set`))
			})
			It("Bad Minutes Value", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "HOURLY", "-c", "3", "-m", "100")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-m|--minute] value must be between 0 and 59."))
			})
			It("Wrong Hour", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "HOURLY", "-c", "3", "-m", "20", "--hour", "50")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-r|--hour] value must be between 0 and 23."))
			})
			It("Wrong Days of the Week", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "HOURLY", "-c", "3", "-m", "20", "--hour", "10", "-d", "10")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-d|--day-of-week] value must be between 0 and 6."))
			})
		})

		Context("Snapshot enable with correct parameters", func() {
			BeforeEach(func() {
				FakeStorageManager.EnableSnapshotReturns(nil)
			})
			It("Happy Path All Options", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "HOURLY", "-c", "3", "-m", "20", "--hour", "10", "-d", "0")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("HOURLY snapshots have been enabled for volume 1234."))
			})
			It("Happy Path min Options", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "9999", "-s", "INTERVAL", "-c", "500")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("INTERVAL snapshots have been enabled for volume 9999."))
			})
		})

		Context("Snapshot enable with correct parameters but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.EnableSnapshotReturns(errors.New("Internal Server Error"))
			})
			It("API error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "HOURLY", "-c", "3", "-m", "20", "--hour", "10", "-d", "0")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to enable HOURLY snapshot for volume 1234."))
			})
		})
	})
})

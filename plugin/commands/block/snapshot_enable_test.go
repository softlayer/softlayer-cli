package block_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
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
		Context("Snapshot enable without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage : This command requires one argument"))
			})
		})
		Context("Snapshot enable with wrong volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
			})
		})

		Context("Snapshot enable without -s", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-s|--schedule-type] is required, options are: HOURLY, DAILY, WEEKLY."))
			})
		})

		Context("Snapshot enable with wrong -s", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "MONTHLY")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-s|--schedule-type] must be HOURLY, DAILY, or WEEKLY."))
			})
		})

		Context("Snapshot enable without -c", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "HOURLY")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-c|--retention-count' is required"))
			})
		})

		Context("Snapshot enable with wrong -m", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "HOURLY", "-c", "3", "-m", "100")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-m|--minute] value must be between 0 and 59."))
			})
		})

		Context("Snapshot enable with wrong --hour", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "HOURLY", "-c", "3", "-m", "20", "--hour", "50")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-r|--hour] value must be between 0 and 23."))
			})
		})

		Context("Snapshot enable with wrong -d", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "HOURLY", "-c", "3", "-m", "20", "--hour", "10", "-d", "10")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-d|--day-of-week] value must be between 0 and 6."))
			})
		})

		Context("Snapshot enable with correct parameters", func() {
			BeforeEach(func() {
				FakeStorageManager.EnableSnapshotReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "HOURLY", "-c", "3", "-m", "20", "--hour", "10", "-d", "0")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("HOURLY snapshots have been enabled for volume 1234."))
			})
		})

		Context("Snapshot enable with correct parameters but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.EnableSnapshotReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "HOURLY", "-c", "3", "-m", "20", "--hour", "10", "-d", "0")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to enable HOURLY snapshot for volume 1234."))

			})
		})
	})
})

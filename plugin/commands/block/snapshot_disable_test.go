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

var _ = Describe("Snapshot disable", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *block.SnapshotDisableCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
		FakeStorageManager *testhelpers.FakeStorageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewSnapshotDisableCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})
	Describe("Snapshot disable", func() {
		Context("Snapshot disable without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("Snapshot disable with wrong volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
			})
		})

		Context("Snapshot disable without -s", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [--schedule-type] is required, options are: HOURLY, DAILY, WEEKLY."))
			})
		})

		Context("Snapshot disable with wrong -s", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "MONTHLY")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [--schedule-type] must be HOURLY, DAILY, or WEEKLY."))
			})
		})

		Context("Snapshot disable with correct volume id and -s", func() {
			BeforeEach(func() {
				FakeStorageManager.DisableSnapshotsReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "DAILY")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("DAILY snapshots have been disabled for volume 1234."))
			})
		})

		Context("Snapshot disable with correct volume id and -s  but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.DisableSnapshotsReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "DAILY")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to disable DAILY snapshot for volume 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})
	})
})

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

var _ = Describe("Snapshot Cancel", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *block.SnapshotCancelCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
		FakeStorageManager *testhelpers.FakeStorageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewSnapshotCancelCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("Snapshot cancel", func() {
		Context("Snapshot cancel without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("Snapshot cancel with wrong volume id", func() {
			It("error resolving volume ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
			})
		})

		Context("Snapshot cancel with correct volume id without -f and not continue", func() {
			BeforeEach(func() {
				FakeStorageManager.CancelSnapshotSpaceReturns(nil)
			})
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("This will cancel the block volume snapshot space: 1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
			})
		})

		Context("Snapshot cancel with correct volume id", func() {
			BeforeEach(func() {
				FakeStorageManager.CancelSnapshotSpaceReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Block volume 1234 has been marked for snapshot cancellation."))
			})
		})

		Context("Snapshot cancel with correct volume id and immediate", func() {
			BeforeEach(func() {
				FakeStorageManager.CancelSnapshotSpaceReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--immediate", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Block volume 1234 has been marked for immediate snapshot"))
			})
		})

		Context("Snapshot cancel with correct volume id but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.CancelSnapshotSpaceReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to cancel snapshot space for volume 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})
	})
})

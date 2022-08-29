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

var _ = Describe("Volume set snapshot notification status", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *block.SnapshotSetNotificationCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
		FakeStorageManager *testhelpers.FakeStorageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewSnapshotSetNotificationCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})
	Describe("Volume set snapshot notification status", func() {
		Context("Volume set snapshot notification status without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("Volume set snapshot notification status with wrong volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
			})
		})

		Context("Volume set snapshot notification status without --enable or --disable", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234567")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Either '--enable' or '--disable' is required."))
			})
		})

		Context("Volume set snapshot notification status with --enable and --disable", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--enable", "--disable", "1234567")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--enable', '--disable' are exclusive"))
			})
		})

		Context("Volume set notification status with an error", func() {
			BeforeEach(func() {
				FakeStorageManager.SetSnapshotNotificationReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--enable", "1234567")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to set the snapshort notification  for volume '1234567'.\n"))
			})
		})

		Context("Volume set notification status with correct volume id and status enabled", func() {
			BeforeEach(func() {
				FakeStorageManager.SetSnapshotNotificationReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--enable", "1234567")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Snapshots space usage threshold warning notification has been set to 'true' for volume '1234567'."))
			})
		})

		Context("Volume set notification status with correct volume id and status disable", func() {
			BeforeEach(func() {
				FakeStorageManager.SetSnapshotNotificationReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--disable", "1234567")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Snapshots space usage threshold warning notification has been set to 'false' for volume '1234567'."))
			})
		})
	})
})

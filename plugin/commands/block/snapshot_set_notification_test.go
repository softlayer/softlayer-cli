package block_test

import (
	"errors"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Volume set snapshot notification status", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.SnapshotSetNotificationCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewSnapshotSetNotificationCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        block.BlockVolumeSnapshotSetNotificationMetaData().Name,
			Description: block.BlockVolumeSnapshotSetNotificationMetaData().Description,
			Usage:       block.BlockVolumeSnapshotSetNotificationMetaData().Usage,
			Flags:       block.BlockVolumeSnapshotSetNotificationMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Volume set snapshot notification status", func() {
		Context("Volume set snapshot notification status without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("Volume set snapshot notification status with wrong volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
			})
		})

		Context("Volume set snapshot notification status without --enable or --disable", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234567")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Either '--enable' or '--disable' is required."))
			})
		})

		Context("Volume set snapshot notification status with --enable and --disable", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--enable", "--disable", "1234567")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--enable', '--disable' are exclusive"))
			})
		})

		Context("Volume set notification status with an error", func() {
			BeforeEach(func() {
				FakeStorageManager.SetSnapshotNotificationReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--enable", "1234567")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to set the snapshort notification  for volume '1234567'.\n"))
			})
		})

		Context("Volume set notification status with correct volume id and status enabled", func() {
			BeforeEach(func() {
				FakeStorageManager.SetSnapshotNotificationReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--enable", "1234567")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Snapshots space usage threshold warning notification has been set to 'true' for volume '1234567'."))
			})
		})

		Context("Volume set notification status with correct volume id and status disable", func() {
			BeforeEach(func() {
				FakeStorageManager.SetSnapshotNotificationReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--disable", "1234567")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Snapshots space usage threshold warning notification has been set to 'false' for volume '1234567'."))
			})
		})
	})
})

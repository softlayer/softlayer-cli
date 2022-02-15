package file_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/file"
	
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Volume snapshot notification status", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *file.SnapshotGetNotificationStatusCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = file.NewSnapshotGetNotificationStatusCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        file.FileVolumeSnapshotGetNotificationStatusMetaData().Name,
			Description: file.FileVolumeSnapshotGetNotificationStatusMetaData().Description,
			Usage:       file.FileVolumeSnapshotGetNotificationStatusMetaData().Usage,
			Flags:       file.FileVolumeSnapshotGetNotificationStatusMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Get volume snapshot notification status", func() {
		Context("Volume get notification status without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})

		Context("Volume get notification status with wrong volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
			})
		})

		Context("Volume get notification status with an error", func() {
			BeforeEach(func() {
				FakeStorageManager.GetSnapshotNotificationStatusReturns(-1, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234567")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get the snapshot notification status for volume '1234567'.\n"))
			})
		})

		Context("Volume get notification status with correct volume id and status enabled", func() {
			BeforeEach(func() {
				FakeStorageManager.GetSnapshotNotificationStatusReturns(1, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234567")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Enabled: Snapshots space usage threshold is enabled for volume '1234567'."))
			})
		})

		Context("Volume get notification status with correct volume id and status disabled", func() {
			BeforeEach(func() {
				FakeStorageManager.GetSnapshotNotificationStatusReturns(0, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234567")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Disabled: Snapshots space usage threshold is disabled for volume '1234567'."))
			})
		})
	})
})

package file_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/file"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Snapshot disable", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *file.SnapshotDisableCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = file.NewSnapshotDisableCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        metadata.FileSnapshotDisableMetaData().Name,
			Description: metadata.FileSnapshotDisableMetaData().Description,
			Usage:       metadata.FileSnapshotDisableMetaData().Usage,
			Flags:       metadata.FileSnapshotDisableMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Snapshot disable", func() {
		Context("Snapshot disable without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Snapshot disable with wrong volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Volume ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Snapshot disable without -s", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [--schedule-type] is required, options are: HOURLY, DAILY, WEEKLY.")).To(BeTrue())
			})
		})

		Context("Snapshot disable with wrong -s", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "-s", "MONTHLY")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [--schedule-type] must be HOURLY, DAILY, or WEEKLY.")).To(BeTrue())
			})
		})

		Context("Snapshot disable with correct volume id and -s", func() {
			BeforeEach(func() {
				FakeStorageManager.DisableSnapshotsReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "DAILY")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"DAILY snapshots have been disabled for volume 1234."}))
			})
		})

		Context("Snapshot disable with correct volume id and -s  but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.DisableSnapshotsReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "DAILY")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to disable DAILY snapshot for volume 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
	})
})

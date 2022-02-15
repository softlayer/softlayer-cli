package file_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/file"
	
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Snapshot enable", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *file.SnapshotEnableCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = file.NewSnapshotEnableCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        file.FileSnapshotEnableMetaData().Name,
			Description: file.FileSnapshotEnableMetaData().Description,
			Usage:       file.FileSnapshotEnableMetaData().Usage,
			Flags:       file.FileSnapshotEnableMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Snapshot enable", func() {
		Context("Snapshot enable without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Snapshot enable with wrong volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Volume ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Snapshot enable without -s", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-s|--schedule-type] is required, options are: HOURLY, DAILY, WEEKLY.")).To(BeTrue())
			})
		})

		Context("Snapshot enable with wrong -s", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "MONTHLY")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-s|--schedule-type] must be HOURLY, DAILY, or WEEKLY.")).To(BeTrue())
			})
		})

		Context("Snapshot enable without -c", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "HOURLY")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: '-c|--retention-count' is required")).To(BeTrue())
			})
		})

		Context("Snapshot enable with wrong -m", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "HOURLY", "-c", "3", "-m", "100")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-m|--minute] value must be between 0 and 59.")).To(BeTrue())
			})
		})

		Context("Snapshot enable with wrong --hour", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "HOURLY", "-c", "3", "-m", "20", "--hour", "50")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-r|--hour] value must be between 0 and 23.")).To(BeTrue())
			})
		})

		Context("Snapshot enable with wrong -d", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "HOURLY", "-c", "3", "-m", "20", "--hour", "10", "-d", "10")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-d|--day-of-week] value must be between 0 and 6.")).To(BeTrue())
			})
		})

		Context("Snapshot enable with correct parameters", func() {
			BeforeEach(func() {
				FakeStorageManager.EnableSnapshotReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "HOURLY", "-c", "3", "-m", "20", "--hour", "10", "-d", "0")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"HOURLY snapshots have been enabled for volume 1234."}))
			})
		})

		Context("Snapshot enable with correct parameters but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.EnableSnapshotReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "HOURLY", "-c", "3", "-m", "20", "--hour", "10", "-d", "0")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to enable HOURLY snapshot for volume 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
	})
})

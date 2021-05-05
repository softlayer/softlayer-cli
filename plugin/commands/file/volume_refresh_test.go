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
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("File Volume Refresh", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *file.VolumeRefreshCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = file.NewVolumeRefreshCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        metadata.FileVolumeRefreshMetaData().Name,
			Description: metadata.FileVolumeRefreshMetaData().Description,
			Usage:       metadata.FileVolumeRefreshMetaData().Usage,
			Action:      cmd.Run,
		}
	})

	Describe("File Volume Refresh", func() {
		Context("No Arguments Error", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires two arguments.")).To(BeTrue())
			})
		})
		Context("Bad VolumeId", func() {
			It("error resolving volume ID", func() {
				err := testhelpers.RunCommand(cliCommand, "abc", "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Volume ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("Bad SnapshotId", func() {
			It("error resolving snapshot ID", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Snapshot ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("Proper Usage", func() {
			BeforeEach(func() {
				FakeStorageManager.VolumeRefreshReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
			})
		})

		Context("Proper Usage, but API error", func() {
			BeforeEach(func() {
				FakeStorageManager.VolumeRefreshReturns(errors.New("Fake Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "5678")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"OK"}))
				Expect(strings.Contains(err.Error(), "Fake Internal Server Error")).To(BeTrue())
			})
		})
	})
})

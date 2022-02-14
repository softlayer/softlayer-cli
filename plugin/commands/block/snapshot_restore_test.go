package block_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Snapshot restore", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.SnapshotRestoreCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewSnapshotRestoreCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        block.BlockSnapshotRestoreMetaData().Name,
			Description: block.BlockSnapshotRestoreMetaData().Description,
			Usage:       block.BlockSnapshotRestoreMetaData().Usage,
			Flags:       block.BlockSnapshotRestoreMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Snapshot restore", func() {
		Context("Snapshot restore without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires two arguments.")).To(BeTrue())
			})
		})
		Context("Snapshot order without snapshot id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires two arguments.")).To(BeTrue())
			})
		})
		Context("Snapshot order with wrong volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc", "123")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Volume ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("Snapshot order with wrong snapshot id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Snapshot ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Snapshot order with correct parameters", func() {
			BeforeEach(func() {
				FakeStorageManager.RestoreFromSnapshotReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Block volume 123 is being restored using snapshot 456."}))
			})
		})

		Context("Snapshot order with correct parameters but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.RestoreFromSnapshotReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "456")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{
					"OK",
				}))
				Expect(strings.Contains(err.Error(), "Failed to restore volume 123 from snapshot 456.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
	})
})

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

var _ = Describe("Snapshot Cancel", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.SnapshotCancelCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewSnapshotCancelCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        block.BlockSnapshotCancelMetaData().Name,
			Description: block.BlockSnapshotCancelMetaData().Description,
			Usage:       block.BlockSnapshotCancelMetaData().Usage,
			Flags:       block.BlockSnapshotCancelMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Snapshot cancel", func() {
		Context("Snapshot cancel without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Snapshot cancel with wrong volume id", func() {
			It("error resolving volume ID", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Volume ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Snapshot cancel with correct volume id without -f and not continue", func() {
			BeforeEach(func() {
				FakeStorageManager.CancelSnapshotSpaceReturns(nil)
			})
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This will cancel the block volume snapshot space: 1234 and cannot be undone. Continue?"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted."}))
			})
		})

		Context("Snapshot cancel with correct volume id", func() {
			BeforeEach(func() {
				FakeStorageManager.CancelSnapshotSpaceReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Block volume 1234 has been marked for snapshot cancellation."}))
			})
		})

		Context("Snapshot cancel with correct volume id and immediate", func() {
			BeforeEach(func() {
				FakeStorageManager.CancelSnapshotSpaceReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--immediate", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Block volume 1234 has been marked for immediate snapshot cancellation."}))
			})
		})

		Context("Snapshot cancel with correct volume id but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.CancelSnapshotSpaceReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to cancel snapshot space for volume 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
	})
})

package block_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/block"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Snapshot Delete", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.SnapshotDeleteCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewSnapshotDeleteCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        metadata.BlockSnapshotDeleteMetaData().Name,
			Description: metadata.BlockSnapshotDeleteMetaData().Description,
			Usage:       metadata.BlockSnapshotDeleteMetaData().Usage,
			Flags:       metadata.BlockSnapshotDeleteMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Snapshot delete", func() {
		Context("Snapshot delete without snapshot id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Snapshot delete with wrong snapshot id", func() {
			It("error resolving snapshot ID", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Snapshot ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Snapshot delete with correct snapshot id", func() {
			BeforeEach(func() {
				FakeStorageManager.DeleteSnapshotReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Snapshot 1234 was deleted."}))
			})
		})

		Context("Snapshot delete with correct snapshot id but snapshot not found", func() {
			BeforeEach(func() {
				FakeStorageManager.DeleteSnapshotReturns(errors.New("SoftLayer_Exception_ObjectNotFound"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Unable to find snapshot with ID 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "SoftLayer_Exception_ObjectNotFound")).To(BeTrue())
			})
		})

		Context("Snapshot delete with correct snapshot id but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.DeleteSnapshotReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to delete snapshot 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
	})
})

package block_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Snapshot Create", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.SnapshotCreateCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewSnapshotCreateCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        block.BlockSnapshotCreateMetaData().Name,
			Description: block.BlockSnapshotCreateMetaData().Description,
			Usage:       block.BlockSnapshotCreateMetaData().Usage,
			Flags:       block.BlockSnapshotCreateMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Snapshot create", func() {
		Context("Snapshot create without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Snapshot create with wrong volume id", func() {
			It("error resolving volume ID", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Volume ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Snapshot create with correct volume id", func() {
			BeforeEach(func() {
				FakeStorageManager.CreateSnapshotReturns(datatypes.Network_Storage{Id: sl.Int(5678)}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"New snapshot 5678 was created."}))
			})
		})

		Context("Snapshot create with correct volume id and additional note", func() {
			BeforeEach(func() {
				FakeStorageManager.CreateSnapshotReturns(datatypes.Network_Storage{Id: sl.Int(5678), Notes: sl.String("my note to create snapshot")}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--note", "my note to create snapshot")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"New snapshot 5678 was created."}))
			})
		})

		Context("Snapshot create with correct volume id but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.CreateSnapshotReturns(datatypes.Network_Storage{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"OK"}))
				Expect(strings.Contains(err.Error(), "Error occurred while creating snapshot for volume 1234.Ensure volume is not failed over or in another state which prevents taking snapshots.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
	})
})

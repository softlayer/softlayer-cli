package block_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Snapshot Create", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *block.SnapshotCreateCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
		FakeStorageManager *testhelpers.FakeStorageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewSnapshotCreateCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("Snapshot create", func() {
		Context("Snapshot create without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("Snapshot create with wrong volume id", func() {
			It("error resolving volume ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
			})
		})

		Context("Snapshot create with correct volume id", func() {
			BeforeEach(func() {
				FakeStorageManager.CreateSnapshotReturns(datatypes.Network_Storage{Id: sl.Int(5678)}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("New snapshot 5678 was created."))
			})
		})

		Context("Snapshot create with correct volume id and additional note", func() {
			BeforeEach(func() {
				FakeStorageManager.CreateSnapshotReturns(datatypes.Network_Storage{Id: sl.Int(5678), Notes: sl.String("my note to create snapshot")}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--note", "my note to create snapshot")
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeUI.Outputs()).To(ContainSubstring("New snapshot 5678 was created."))
			})
		})

		Context("Snapshot create with correct volume id but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.CreateSnapshotReturns(datatypes.Network_Storage{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Error occurred while creating snapshot for volume 1234.Ensure volume is not failed over or in another state which prevents taking snapshots."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})
	})
})

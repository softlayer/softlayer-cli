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

var _ = Describe("Replica locations", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *block.ReplicaLocationsCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
		FakeStorageManager *testhelpers.FakeStorageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewReplicaLocationsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("Replicant locations", func() {
		Context("replicant locations without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage : This command requires one argument"))

			})
		})
		Context("Replicant locations with server error", func() {
			BeforeEach(func() {
				FakeStorageManager.GetReplicationLocationsReturns(nil, errors.New("Internal server error"))
			})
			It("error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("OK"))
				Expect(err.Error()).To(ContainSubstring("Failed to get datacenters for volume 1234."))
			})
		})
		Context("Replicant locations", func() {
			BeforeEach(func() {
				FakeStorageManager.GetReplicationLocationsReturns(
					[]datatypes.Location{
						datatypes.Location{
							Id:       sl.Int(352494),
							Name:     sl.String("hkg02"),
							LongName: sl.String("Hong Kong 2")},
					}, nil)
			})
			It("succeed", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("352494"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("hkg02"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hong Kong 2"))

			})
		})
		Context("Replicant locations", func() {
			BeforeEach(func() {
				FakeStorageManager.GetReplicationLocationsReturns(
					[]datatypes.Location{}, nil)
			})
			It("succeed", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("No data centers compatible for replication."))
			})
		})
	})
})

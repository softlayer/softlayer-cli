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

var _ = Describe("Replica locations", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.ReplicaLocationsCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewReplicaLocationsCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        block.BlockReplicaLocationsMetaData().Name,
			Description: block.BlockReplicaLocationsMetaData().Description,
			Usage:       block.BlockReplicaLocationsMetaData().Usage,
			Flags:       block.BlockReplicaLocationsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Replicant locations", func() {
		Context("replicant locations without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Replicant locations with server error", func() {
			BeforeEach(func() {
				FakeStorageManager.GetReplicationLocationsReturns(nil, errors.New("Internal server error"))
			})
			It("error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{
					"OK",
				}))
				Expect(strings.Contains(err.Error(), "Failed to get datacenters for volume 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
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
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"352494"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"hkg02"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Hong Kong 2"}))

			})
		})
		Context("Replicant locations", func() {
			BeforeEach(func() {
				FakeStorageManager.GetReplicationLocationsReturns(
					[]datatypes.Location{}, nil)
			})
			It("succeed", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{
					"No data centers compatible for replication.",
				}))
			})
		})
	})
})

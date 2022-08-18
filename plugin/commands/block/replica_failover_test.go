package block_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Replica failover", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *block.ReplicaFailoverCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		FakeStorageManager *testhelpers.FakeStorageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = block.NewReplicaFailoverCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("Replicant failover", func() {
		Context("replicant failover without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires two arguments."))
			})
		})
		Context("replicant failover without replicant id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires two arguments."))
			})
		})
		Context("Replicant fail over with wrong volume id", func() {
			It("error resolving volume ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
			})
		})
		Context("Replicant fail over with wrong replica id", func() {
			It("error resolving volume ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Replica ID'. It must be a positive integer."))
			})
		})

		Context("Replicant fail over with correct volume id and replica id", func() {
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Failover of volume 1234 to replica 5678 is now in progress."))
			})
		})

		Context("Replicant fail over with correct volume id but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.FailOverToReplicantReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "5678")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("OK"))
				Expect(err.Error()).To(ContainSubstring("Failover operation could not be initiated for volume 1234."))
			})
		})
	})
})

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

var _ = Describe("block subnets-assign", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cliCommand         *block.SubnetsAssignCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewSubnetsAssignCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("block subnets-assign", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--subnet-id", "111111")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcde", "--subnet-id", "111111")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Access ID'. It must be a positive integer."))
			})

			It("Set without option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`"subnet-id" not set`))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				FakeStorageManager.AssignSubnetsToAclReturns([]int{}, errors.New("Failed to get subnets."))
			})
			It("Failed get subnets", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--subnet-id", "111111")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to assign subnets."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerSubnets := []int{
					111111,
				}
				FakeStorageManager.AssignSubnetsToAclReturns(fakerSubnets, nil)
			})
			It("Assign subnets", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--subnet-id", "111111")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Successfully assigned subnet"))
			})
		})

	})
})

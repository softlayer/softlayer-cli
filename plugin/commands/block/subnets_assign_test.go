package block_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("block subnets-assign", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.SubnetsAssignCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewSubnetsAssignCommand(fakeUI, fakeStorageManager)
		cliCommand = cli.Command{
			Name:        block.BlockSubnetsAssignMetaData().Name,
			Description: block.BlockSubnetsAssignMetaData().Description,
			Usage:       block.BlockSubnetsAssignMetaData().Usage,
			Flags:       block.BlockSubnetsAssignMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("block subnets-assign", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCommand(cliCommand, "--subnet-id", "111111")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCommand(cliCommand, "abcde", "--subnet-id", "111111")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Access ID'. It must be a positive integer."))
			})

			It("Set without option", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`Required flag "subnet-id" not set`))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeStorageManager.AssignSubnetsToAclReturns([]int{}, errors.New("Failed to get subnets."))
			})
			It("Failed get subnets", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--subnet-id", "111111")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to assign subnets."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerSubnets := []int{
					111111,
				}
				fakeStorageManager.AssignSubnetsToAclReturns(fakerSubnets, nil)
			})
			It("Assign subnets", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--subnet-id", "111111")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Successfully assigned subnet"))
			})
		})

	})
})

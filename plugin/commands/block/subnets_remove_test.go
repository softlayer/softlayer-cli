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

var _ = Describe("block subnets-remove", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.SubnetsRemoveCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewSubnetsRemoveCommand(fakeUI, fakeStorageManager)
		cliCommand = cli.Command{
			Name:        block.BlockSubnetsRemoveMetaData().Name,
			Description: block.BlockSubnetsRemoveMetaData().Description,
			Usage:       block.BlockSubnetsRemoveMetaData().Usage,
			Flags:       block.BlockSubnetsRemoveMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("block subnets-remove", func() {

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
				fakeStorageManager.RemoveSubnetsFromAclReturns([]int{}, errors.New("Failed to get subnets."))
			})
			It("Failed get subnets", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--subnet-id", "111111")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to remove subnets."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerSubnets := []int{
					111111,
				}
				fakeStorageManager.RemoveSubnetsFromAclReturns(fakerSubnets, nil)
			})
			It("remove subnets", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--subnet-id", "111111")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Successfully removed subnet"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeStorageManager.RemoveSubnetsFromAclReturns([]int{}, nil)
			})
			It("not remove subnets", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--subnet-id", "111111")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Failed to remove subnet"))
			})
		})

	})
})

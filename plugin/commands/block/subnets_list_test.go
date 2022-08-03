package block_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("block subnets-list", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.SubnetsListCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewSubnetsListCommand(fakeUI, fakeStorageManager)
		cliCommand = cli.Command{
			Name:        block.BlockSubnetsListMetaData().Name,
			Description: block.BlockSubnetsListMetaData().Description,
			Usage:       block.BlockSubnetsListMetaData().Usage,
			Flags:       block.BlockSubnetsListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("block subnets-list", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCommand(cliCommand, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Access ID'. It must be a positive integer."))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeStorageManager.GetSubnetsInAclReturns([]datatypes.Network_Subnet{}, errors.New("Failed to get subnets."))
			})
			It("Failed get subnets", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get subnets."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerSubnets := []datatypes.Network_Subnet{
					datatypes.Network_Subnet{
						Id:                sl.Int(123456),
						NetworkIdentifier: sl.String("11.22.33.44"),
						Cidr:              sl.Int(111111),
					},
				}
				fakeStorageManager.GetSubnetsInAclReturns(fakerSubnets, nil)
			})
			It("List subnets", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("11.22.33.44"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
			})
		})

	})
})

package subnet_test

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
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/subnet"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Subnet create", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *subnet.CreateCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = subnet.NewCreateCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        metadata.SubnetCreateMetaData().Name,
			Description: metadata.SubnetCreateMetaData().Description,
			Usage:       metadata.SubnetCreateMetaData().Usage,
			Flags:       metadata.SubnetCreateMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Subnet create", func() {
		Context("Subnet create with not enough parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires three arguments.")).To(BeTrue())
			})
		})

		Context("Subnet create with wrong network", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc", "8", "123")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: NETWORK has to be either public or private.")).To(BeTrue())
			})
		})

		Context("Subnet create with wrong quantity", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "public", "abc", "123")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'QUANTITY'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Subnet create with wrong vlanID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "public", "8", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'VLAN ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Subnet create without -f", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "public", "8", "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This action will incur charges on your account. Continue?"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted."}))
			})
		})

		Context("Subnet create with -test", func() {
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "public", "8", "123", "--test")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The order is correct."}))
			})
		})

		Context("Subnet create with correct parameters but server fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.AddSubnetReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "public", "8", "123", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to add subnet.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Subnet create with correct parameters", func() {
			BeforeEach(func() {
				fakeNetworkManager.AddSubnetReturns(datatypes.Container_Product_Order_Receipt{OrderId: sl.Int(12345678)}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "public", "8", "123", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Order 12345678 was placed."}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "public", "8", "123", "-f", "--v6")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Order 12345678 was placed."}))
			})
		})
	})
})

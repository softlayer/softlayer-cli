package subnet_test

import (
	"errors"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/subnet"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Subnet Route", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *subnet.RouteCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = subnet.NewRouteCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        metadata.SubnetRouteMetaData().Name,
			Description: metadata.SubnetRouteMetaData().Description,
			Usage:       metadata.SubnetRouteMetaData().Usage,
			Flags:       metadata.SubnetRouteMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Subnet Route", func() {
		Context("Subnet route without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("Subnet route with wrong subnet ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Subnet ID'. It must be a positive integer."))
			})
		})

		Context("Subnet route without -t", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-i", "1234", "1234567")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-t/--type] is required."))
			})
		})

		Context("Subnet route without -i", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "SoftLayer_Network_Subnet_IpAddress", "1234567")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-i/--type-id] is required."))
			})
		})

		Context("Subnet route failed", func() {
			BeforeEach(func() {
				fakeNetworkManager.RouteReturns(false, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-i", "1234", "-t", "SoftLayer_Network_Subnet_IpAddress", "1234567")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to route using the type: SoftLayer_Network_Subnet_IpAddress and identifier: 1234.\n"))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Subnet route with correct parameters", func() {
			BeforeEach(func() {
				fakeNetworkManager.RouteReturns(true, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "-i", "1234", "-t", "SoftLayer_Network_Subnet_IpAddress", "1234567")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The transaction to route is created, routes will be updated in one or two minutes."}))
			})
		})
	})
})

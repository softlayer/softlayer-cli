package subnet_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/subnet"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Subnet Clear Route", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *subnet.ClearRouteCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = subnet.NewClearRouteCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        metadata.SubnetClearRouteMetaData().Name,
			Description: metadata.SubnetClearRouteMetaData().Description,
			Usage:       metadata.SubnetClearRouteMetaData().Usage,
			Flags:       metadata.SubnetClearRouteMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Subnet Clear Route", func() {
		Context("Subnet clear route without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Subnet clear route with wrong subnet ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Subnet ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Subnet clear route failed", func() {
			BeforeEach(func() {
				fakeNetworkManager.ClearRouteReturns(false, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234567")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to clear the route for the subnet: 1234567.\n")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Subnet clear route with correct parameters", func() {
			BeforeEach(func() {
				fakeNetworkManager.RouteReturns(true, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234567")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The transaction to clear the route is created, routes will be updated in one or two minutes."}))
			})
		})
	})
})

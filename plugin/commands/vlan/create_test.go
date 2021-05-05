package vlan_test

import (
	"errors"
	"fmt"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/vlan"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VLAN create", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *vlan.CreateCommand
		cliCommand         cli.Command
		context            plugin.PluginContext
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		context = plugin.InitPluginContext("softlayer")
		cmd = vlan.NewCreateCommand(fakeUI, fakeNetworkManager, context)
		cliCommand = cli.Command{
			Name:        metadata.VlanCreateMetaData().Name,
			Description: metadata.VlanCreateMetaData().Description,
			Usage:       metadata.VlanCreateMetaData().Usage,
			Flags:       metadata.VlanCreateMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VLAN create", func() {
		Context("VLAN create with -r and -d", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-r", "router123", "-d", "dal09")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-r|--router] is not allowed with [-d|--datacenter] or [-t|--vlan-type].")).To(BeTrue())
				Expect(err.Error()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl vlan options' to check available options.", cmd.Context.CLIName())}))
			})
		})

		Context("VLAN create with -r and -t", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-r", "router123", "-t", "public")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-r|--router] is not allowed with [-d|--datacenter] or [-t|--vlan-type].")).To(BeTrue())
				Expect(err.Error()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl vlan options' to check available options.", cmd.Context.CLIName())}))
			})
		})

		Context("VLAN create with -t but no -d", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "public")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-d|--datacenter] and [-t|--vlan-type] are required.")).To(BeTrue())
				Expect(err.Error()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl vlan options' to check available options.", cmd.Context.CLIName())}))
			})
		})

		Context("VLAN create with -d but no -t", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-d", "dal10")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-d|--datacenter] and [-t|--vlan-type] are required.")).To(BeTrue())
				Expect(err.Error()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl vlan options' to check available options.", cmd.Context.CLIName())}))
			})
		})

		Context("VLAN create with -d and -t but not continue", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "-t", "public", "-d", "dal10")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This action will incur charges on your account. Continue?"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted."}))
			})
		})

		Context("VLAN create with correct parameters but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.AddVlanReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "public", "-d", "dal10", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to add VLAN.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("VLAN create with correct parameters", func() {
			BeforeEach(func() {
				fakeNetworkManager.AddVlanReturns(datatypes.Container_Product_Order_Receipt{OrderId: sl.Int(12345678)}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "public", "-d", "dal10", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The order 12345678 was placed."}))
			})
		})
	})
})

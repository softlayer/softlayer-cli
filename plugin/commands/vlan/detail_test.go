package vlan_test

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
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/vlan"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VLAN Detail", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *vlan.DetailCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = vlan.NewDetailCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        vlan.VlanDetailMetaData().Name,
			Description: vlan.VlanDetailMetaData().Description,
			Usage:       vlan.VlanDetailMetaData().Usage,
			Flags:       vlan.VlanDetailMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VLAN detail", func() {
		Context("VLAN detail without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("VLAN detail with wrong vlan id", func() {
			It("error resolving vlan ID", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'VLAN ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("VLAN detail with correct vlan id but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetVlanReturns(datatypes.Network_Vlan{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get VLAN: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("VLAN detail with correct vlan id", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetVlanReturns(datatypes.Network_Vlan{
					Id:         sl.Int(1234),
					VlanNumber: sl.Int(100),
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"100"}))
			})
		})
	})
})

package globalip_test

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
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/globalip"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("GlobalIP list", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *globalip.ListCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = globalip.NewListCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        globalip.GlobalIpListMetaData().Name,
			Description: globalip.GlobalIpListMetaData().Description,
			Usage:       globalip.GlobalIpListMetaData().Usage,
			Flags:       globalip.GlobalIpListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("GlobalIP list", func() {
		Context("GlobalIP list with both v4 and v6", func() {
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--v4", "--v6")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [--v4] is not allowed with [--v6].")).To(BeTrue())
			})
		})

		Context("GlobalIP list server fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListGlobalIPsReturns([]datatypes.Network_Subnet_IpAddress_Global{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to list global IPs on your account.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("GlobalIP list ", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListGlobalIPsReturns([]datatypes.Network_Subnet_IpAddress_Global{
					datatypes.Network_Subnet_IpAddress_Global{
						Id: sl.Int(123456),
						IpAddress: &datatypes.Network_Subnet_IpAddress{
							IpAddress: sl.String("5.6.7.8"),
						},
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"123456"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"5.6.7.8"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"No"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"None"}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--v4")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"123456"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"5.6.7.8"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"No"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"None"}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--v6")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"123456"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"5.6.7.8"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"No"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"None"}))
			})
		})
	})
})

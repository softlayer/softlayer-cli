package globalip_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/globalip"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var globalIpReturn = []datatypes.Network_Subnet_IpAddress_Global{
	datatypes.Network_Subnet_IpAddress_Global{
		Id: sl.Int(123456),
		IpAddress: &datatypes.Network_Subnet_IpAddress{
			IpAddress: sl.String("5.6.7.8"),
			SubnetId: sl.Int(998877),
		},
	},
}

var _ = Describe("GlobalIP list", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *globalip.ListCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeNetworkManager *testhelpers.FakeNetworkManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = globalip.NewListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cliCommand.NetworkManager = fakeNetworkManager
	})

	Describe("GlobalIP list", func() {
		Context("GlobalIP list with both v4 and v6", func() {
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--v4", "--v6")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [--v4] is not allowed with [--v6].")).To(BeTrue())
			})
		})

		Context("GlobalIP list server fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListGlobalIPsReturns([]datatypes.Network_Subnet_IpAddress_Global{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to list global IPs on your account.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("GlobalIP list ", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListGlobalIPsReturns(globalIpReturn, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("998877"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("5.6.7.8"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("No"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("None"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--v4")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("998877"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("5.6.7.8"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("No"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("None"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--v6")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("998877"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("5.6.7.8"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("No"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("None"))
			})
		})
	})
})

package subnet_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/subnet"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("subnet edit-ip", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *subnet.EditIpCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeNetworkManager *testhelpers.FakeNetworkManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = subnet.NewEditIpCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cliCommand.NetworkManager = fakeNetworkManager
	})

	Describe("subnet edit-ip", func() {

		Context("Return error", func() {
			It("Set command without Argument", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})

			It("Set command without option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--note' is required"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetIpByAddressReturns(datatypes.Network_Subnet_IpAddress{}, errors.New("Failed to get Subnet IP by address"))
			})
			It("Failed get IP object by address", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "11.22.33.44", "--note=myNote")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get Subnet IP by address"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetIpByAddressReturns(datatypes.Network_Subnet_IpAddress{}, nil)
			})
			It("Set command with and inexistent ip address", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "11.22.33.44", "--note=myNote")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to find object"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakerSubnetIp := datatypes.Network_Subnet_IpAddress{
					Id: sl.Int(123456),
				}
				fakeNetworkManager.GetIpByAddressReturns(fakerSubnetIp, nil)
				fakeNetworkManager.EditSubnetIpAddressReturns(false, errors.New("Failed to set note"))
			})
			It("Failed set note", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "11.22.33.44", "--note=myNote")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to set note"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeNetworkManager.EditSubnetIpAddressReturns(true, nil)
			})

			It("Set note", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--note=myNote")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Set note successfully"))
			})
		})
	})
})

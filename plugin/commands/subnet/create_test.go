package subnet_test

import (
	"errors"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/subnet"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Subnet create", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *subnet.CreateCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeNetworkManager *testhelpers.FakeNetworkManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = subnet.NewCreateCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cliCommand.NetworkManager = fakeNetworkManager
	})

	Describe("Subnet create", func() {
		Context("Subnet create with not enough parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires three arguments."))
			})
		})

		Context("Subnet create with wrong network", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc", "8", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: NETWORK has to be either public or private."))
			})
		})

		Context("Subnet create with wrong quantity", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "public", "abc", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'QUANTITY'. It must be a positive integer."))
			})
		})

		Context("Subnet create with wrong vlanID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "public", "8", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'VLAN ID'. It must be a positive integer."))
			})
		})

		Context("Subnet create without -f", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "public", "8", "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This action will incur charges on your account. Continue?"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted."}))
			})
		})

		Context("Subnet create with -test", func() {
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "public", "8", "123", "--test")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The order is correct."}))
			})
		})

		Context("Subnet create with correct parameters but server fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.AddSubnetReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "public", "8", "123", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to add subnet."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Subnet create with correct parameters", func() {
			BeforeEach(func() {
				fakeNetworkManager.AddSubnetReturns(datatypes.Container_Product_Order_Receipt{OrderId: sl.Int(12345678)}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "public", "8", "123", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Order 12345678 was placed."}))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "public", "8", "123", "-f", "--ipv6")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Order 12345678 was placed."}))
			})
		})
	})
})

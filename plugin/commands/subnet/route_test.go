package subnet_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/subnet"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Subnet Route", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *subnet.RouteCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeNetworkManager *testhelpers.FakeNetworkManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = subnet.NewRouteCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cliCommand.NetworkManager = fakeNetworkManager
	})

	Describe("Subnet Route", func() {
		Context("Subnet route without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage : This command requires one argument"))
			})
		})
		Context("Subnet route with bad id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Subnet ID'"))
			})
		})
		Context("Subnet route without args", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--ip, --server, --vsi or --vlan' is required"))
			})
		})
		Context("IP Route happy path", func() {
			It("IP Success", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345", "--ip", "192.168.1.1")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The transaction to route is created"))
			})
		})
		Context("Server Route happy path", func() {
			It("Server Success", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345", "--server", "test<domain.com>")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The transaction to route is created"))
			})
		})
		Context("VSI Route happy path", func() {
			It("VSI Success", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345", "--vsi", "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The transaction to route is created"))
			})
		})
		Context("VLAN Route happy path", func() {
			It("VLAN Success", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345", "--vlan", "999999")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The transaction to route is created"))
			})
		})
		Context("VLAN Route error returned", func() {
			It("Vlan Route Failure", func() {
				fakeNetworkManager.RouteReturns(false, errors.New("SoftLayer_API_Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345", "--vlan", "999999")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_API_Error"))
			})
		})
		Context("Prints Help Text", func() {
			It("Help Text", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345", "--vroom", "999999")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown flag: --vroom"))
			})
		})
	})
})

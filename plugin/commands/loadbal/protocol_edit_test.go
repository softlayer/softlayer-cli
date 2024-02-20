package loadbal_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
)

type OptionMapping struct {
	SLApiConfig datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration
	CLIValue    string
}

var _ = Describe("LoadBal_protocol-edit_Test", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *loadbal.ProtocolEditCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeLBManager *testhelpers.FakeLoadBalancerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = loadbal.NewProtocolEditCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cliCommand.LoadBalancerManager = fakeLBManager
	})

	Context("No LB ID", func() {
		It("Error No Id", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--id' is required"))
		})
		It("Error unable to find Id", func() {
			fakeLBManager.GetLoadBalancerUUIDReturns("-", errors.New("SoftLayer_Exception_ApiError"))
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--id", "12345")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to get load balancer: SoftLayer_Exception_ApiError"))
		})
	})

	Context("No Listener UUID", func() {
		It("Error no UUID", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--id", "12345")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("'--protocol-uuid' is required"))
		})
	})

	Context("Testing Options", func() {
		It("with all arguments", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--id", "12345", "--protocol-uuid", "abc123", "--front-protocol", "HTTP", "--back-protocol", "HTTP", "--front-port", "80", "--back-port", "80", "--method", "ROUNDROBIN", "--client-timeout", "100", "--server-timeout", "100", "--sticky", "cookie", "--connections", "5")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Protocol edited"))
		})
		It("with sticky as source-ip", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--id", "12345", "--protocol-uuid", "abc123", "--front-protocol", "HTTP", "--back-protocol", "HTTP", "--front-port", "80", "--back-port", "80", "--method", "ROUNDROBIN", "--client-timeout", "100", "--server-timeout", "100", "--sticky", "source-ip", "--connections", "5")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Protocol edited"))
		})
		It("with wrong sticky", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--id", "12345", "--protocol-uuid", "abc123", "--sticky", "abc")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Value of option '--sticky' should be cookie or source-ip"))
		})
	})

	Context("API Error", func() {
		It("Handles API Error", func() {
			fakeLBManager.AddLoadBalancerListenerReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("SL_API_ERROR"))
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--id", "12345", "--protocol-uuid", "aa1122", "--server-timeout", "100")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("SL_API_ERROR"))
		})
	})
})

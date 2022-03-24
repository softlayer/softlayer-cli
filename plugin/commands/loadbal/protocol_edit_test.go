package loadbal_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"github.com/softlayer/softlayer-go/datatypes"
)

type OptionMapping struct {
	SLApiConfig datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration
	CLIValue    string
}

var _ = Describe("LoadBal_protocol-edit_Test", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeLBManager *testhelpers.FakeLoadBalancerManager
		cmd           *loadbal.ProtocolEditCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewProtocolEditCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        loadbal.LoadbalProtocolEditMetadata().Name,
			Description: loadbal.LoadbalProtocolEditMetadata().Description,
			Usage:       loadbal.LoadbalProtocolEditMetadata().Usage,
			Flags:       loadbal.LoadbalProtocolEditMetadata().Flags,
			Action:      cmd.Run,
		}
	})

	Context("No LB ID", func() {
		It("Error No Id", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--id' is required"))
		})
		It("Error unable to find Id", func() {
			fakeLBManager.GetLoadBalancerUUIDReturns("-", errors.New("SoftLayer_Exception_ApiError"))
			err := testhelpers.RunCommand(cliCommand, "--id", "12345")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to get load balancer: SoftLayer_Exception_ApiError"))
		})
	})

	Context("No Listener UUID", func() {
		It("Error no UUID", func() {
			err := testhelpers.RunCommand(cliCommand, "--id", "12345")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("'--protocol-uuid' is required"))
		})
	})

	Context("Testing Options", func() {
		It("with all arguments", func() {
			err := testhelpers.RunCommand(cliCommand, "--id", "12345", "--protocol-uuid", "abc123", "--front-protocol", "HTTP", "--back-protocol", "HTTP", "--front-port", "80", "--back-port", "80", "--method", "ROUNDROBIN", "--client-timeout", "100", "--server-timeout", "100", "--sticky", "cookie", "--connections", "5")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Protocol edited"))
		})
		It("with sticky as source-ip", func() {
			err := testhelpers.RunCommand(cliCommand, "--id", "12345", "--protocol-uuid", "abc123", "--front-protocol", "HTTP", "--back-protocol", "HTTP", "--front-port", "80", "--back-port", "80", "--method", "ROUNDROBIN", "--client-timeout", "100", "--server-timeout", "100", "--sticky", "source-ip", "--connections", "5")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Protocol edited"))
		})
		It("with wrong sticky", func() {
			err := testhelpers.RunCommand(cliCommand, "--id", "12345", "--protocol-uuid", "abc123", "--sticky", "abc")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Value of option '--sticky' should be cookie or source-ip"))
		})
	})
	
	Context("API Error", func() {
		It("Handles API Error", func() {
			fakeLBManager.AddLoadBalancerListenerReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("SL_API_ERROR"))
			err := testhelpers.RunCommand(cliCommand, "--id", "12345", "--protocol-uuid", "aa1122", "--server-timeout", "100")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("SL_API_ERROR"))
		})
	})
})

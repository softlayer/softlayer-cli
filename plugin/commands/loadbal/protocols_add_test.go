package loadbal_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"github.com/softlayer/softlayer-go/datatypes"
)

var _ = Describe("LoadBal_protocol-add_Test", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeLBManager *testhelpers.FakeLoadBalancerManager
		cmd           *loadbal.ProtocolAddCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewProtocolAddCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        metadata.LoadbalProtocolAddMetadata().Name,
			Description: metadata.LoadbalProtocolAddMetadata().Description,
			Usage:       metadata.LoadbalProtocolAddMetadata().Usage,
			Flags:       metadata.LoadbalProtocolAddMetadata().Flags,
			Action:      cmd.Run,
		}
	})

	Context("CLI Usage Errors", func() {
		It("Error No Id", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--id' is required"))
		})
		It("Error unable to find Id", func(){
			fakeLBManager.GetLoadBalancerUUIDReturns("-", errors.New("SoftLayer_Exception_ApiError"))
			err := testhelpers.RunCommand(cliCommand, "--id", "12345")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to get load balancer: SoftLayer_Exception_ApiError"))
		})
		It("Error bad stick option", func(){
			err := testhelpers.RunCommand(cliCommand, "--id", "12345", "--sticky", "bad_option")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Value of option '--sticky' should be cookie or source-ip"))
		})
	})
	Context("Setting Options", func() {
		BeforeEach(func() {
			fakeLBManager.GetLoadBalancerUUIDReturns("aaa-bbb-111", nil)
		})
		It("All Options", func(){
			err := testhelpers.RunCommand(cliCommand, "--id", "12345", "--front-protocol", "HTTP",
				"--back-protocol", "HTTPS", "--front-port", "99", "--back-port", "81", "-m", "TEST", "--sticky", "source-ip",
				"-c", "500", "--client-timeout", "100", "--server-timeout", "200",
			)
			Expect(err).NotTo(HaveOccurred())
			lbUUID, argsForCall := fakeLBManager.AddLoadBalancerListenerArgsForCall(0)
			Expect(*lbUUID).To(Equal("aaa-bbb-111"))
			Expect(len(argsForCall)).To(Equal(1))
			Expect(*argsForCall[0].BackendPort).To(Equal(81))
			Expect(*argsForCall[0].FrontendPort).To(Equal(99))
			Expect(*argsForCall[0].FrontendProtocol).To(Equal("HTTP"))
			Expect(*argsForCall[0].BackendProtocol).To(Equal("HTTPS"))
			Expect(*argsForCall[0].LoadBalancingMethod).To(Equal("TEST"))
			Expect(*argsForCall[0].SessionType).To(Equal("SOURCE_IP"))
			Expect(*argsForCall[0].MaxConn).To(Equal(500))
			Expect(*argsForCall[0].ClientTimeout).To(Equal(100))
			Expect(*argsForCall[0].ServerTimeout).To(Equal(200))
			Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
		})
		It("No Options", func(){
			err := testhelpers.RunCommand(cliCommand, "--id", "12345")
			Expect(err).NotTo(HaveOccurred())
			lbUUID, argsForCall := fakeLBManager.AddLoadBalancerListenerArgsForCall(0)
			Expect(*lbUUID).To(Equal("aaa-bbb-111"))
			Expect(len(argsForCall)).To(Equal(1))
			Expect(*argsForCall[0].BackendPort).To(Equal(80))
			Expect(*argsForCall[0].FrontendPort).To(Equal(80))
			Expect(*argsForCall[0].FrontendProtocol).To(Equal("HTTP"))
			Expect(*argsForCall[0].BackendProtocol).To(Equal("HTTP"))
			Expect(*argsForCall[0].LoadBalancingMethod).To(Equal("ROUNDROBIN"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
		})

	})
	Context("API Error", func() {
		It("Handles API Error", func() {
			fakeLBManager.AddLoadBalancerListenerReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("SL_API_ERROR"))
			err := testhelpers.RunCommand(cliCommand, "--id", "12345", "--server-timeout", "100")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("SL_API_ERROR"))
		})
	})
})

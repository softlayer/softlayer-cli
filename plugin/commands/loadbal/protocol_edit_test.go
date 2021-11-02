package loadbal_test

import (
	"errors"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"github.com/softlayer/softlayer-go/datatypes"
)

type OptionMapping struct {
	SLApiConfig datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration
	CLIValue    string
}

var stringValue = "HTTP"
var intValue = 80

// This lets us not have to spell out every test to test every option.
var optionMap = map[string]OptionMapping{
	"front-protocol": OptionMapping{
		SLApiConfig: datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{FrontendProtocol: &stringValue},
		CLIValue:    stringValue,
	},
	"back-protocol": OptionMapping{
		SLApiConfig: datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{BackendProtocol: &stringValue},
		CLIValue:    stringValue,
	},
	"front-port": OptionMapping{
		SLApiConfig: datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{FrontendPort: &intValue},
		CLIValue:    strconv.Itoa(intValue),
	},
	"back-port": OptionMapping{
		SLApiConfig: datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{BackendPort: &intValue},
		CLIValue:    strconv.Itoa(intValue),
	},
	"m": OptionMapping{
		SLApiConfig: datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{LoadBalancingMethod: &stringValue},
		CLIValue:    stringValue,
	},
	"client-timeout": OptionMapping{
		SLApiConfig: datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{ClientTimeout: &intValue},
		CLIValue:    strconv.Itoa(intValue),
	},
	"server-timeout": OptionMapping{
		SLApiConfig: datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{ServerTimeout: &intValue},
		CLIValue:    strconv.Itoa(intValue),
	},
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
			Name:        metadata.LoadbalProtocolEditMetadata().Name,
			Description: metadata.LoadbalProtocolEditMetadata().Description,
			Usage:       metadata.LoadbalProtocolEditMetadata().Usage,
			Flags:       metadata.LoadbalProtocolEditMetadata().Flags,
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
		var listenerUUID = "aaasssbbb-123"
		BeforeEach(func() {
			fakeLBManager.GetLoadBalancerUUIDReturns("aaa-bbb-111", nil)
		})
		for opt, slOpt := range optionMap {
			Context("Test "+opt, func() {
				It("Matches API call", func() {
					err := testhelpers.RunCommand(cliCommand, "--id", "12345", "--protocol-uuid", listenerUUID, "--"+opt, slOpt.CLIValue)
					slOpt.SLApiConfig.ListenerUuid = &listenerUUID
					Expect(err).NotTo(HaveOccurred())
					lbUUID, argsForCall := fakeLBManager.AddLoadBalancerListenerArgsForCall(0)
					Expect(*lbUUID).To(Equal("aaa-bbb-111"))
					Expect(len(argsForCall)).To(Equal(1))
					Expect(argsForCall[0]).To(Equal(slOpt.SLApiConfig))
					Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				})
			})
		}
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

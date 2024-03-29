package loadbal_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Load balancer add policies", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *loadbal.L7PolicyAddCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeLBManager *testhelpers.FakeLoadBalancerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = loadbal.NewL7PolicyAddCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cliCommand.LoadBalancerManager = fakeLBManager
	})

	Context("CLI Usage Errors", func() {
		It("Error No protocol-uuid", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("--protocol-uuid"))
		})
		It("Error missing name", func() {
			command := "--protocol-uuid uuid-12345 -a test-action"
			command_args := strings.Fields(command)
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("-n, --name"))
		})
		It("Error missing action", func() {
			command := "--protocol-uuid uuid-12345 -n test-name"
			command_args := strings.Fields(command)
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("-a, --action"))
		})

		It("No valid action", func() {
			command := "--protocol-uuid uuid-12345 -n test-name -a unknown-action"
			command_args := strings.Fields(command)
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(
				ContainSubstring("-a, --action should be REJECT | REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"))
		})
		It("Error invalid usage for REJECT", func() {
			command := "--protocol-uuid test-12345 -n test-name -a REJECT -r REDIRECT_URL"
			command_args := strings.Fields(command)
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(
				"-r, --redirect is only available with action REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"))
		})
		It("Error No --redirect", func() {
			command := "--protocol-uuid uuis-12345 -n test-name -a REDIRECT_URL"
			command_args := strings.Fields(command)
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(
				"-r, --redirect is required with action REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"))
		})

	})

	Context("API Error", func() {
		It("Handles API Error", func() {
			command := "--protocol-uuid uuid-12345 -n test-name -a REJECT"
			command_args := strings.Fields(command)
			println(command_args)
			fakeLBManager.AddL7PolicyReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("SL_API_ERROR"))
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("SL_API_ERROR"))

		})
	})

	Context("CLI Usage", func() {
		It("REJECT", func() {
			command := "--protocol-uuid uuid-12345 -n test-reject -a REJECT"
			command_args := strings.Fields(command)
			fakeLBManager.AddL7PolicyReturns(datatypes.Network_LBaaS_LoadBalancer{}, nil)
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).NotTo(HaveOccurred())
		})
		It("REDIRECT_POOL", func() {
			command := "--protocol-uuid uuid-12345 -n test-pool -a REDIRECT_POOL -r uuid-pool"
			command_args := strings.Fields(command)
			fakeLBManager.AddL7PolicyReturns(datatypes.Network_LBaaS_LoadBalancer{}, nil)
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).NotTo(HaveOccurred())
		})
		It("REDIRECT_URL", func() {
			command := "--protocol-uuid uuid-12345 -n test-url -a REDIRECT_URL -r http://example.com"
			command_args := strings.Fields(command)
			fakeLBManager.AddL7PolicyReturns(datatypes.Network_LBaaS_LoadBalancer{}, nil)
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).NotTo(HaveOccurred())
		})

		It("REDIRECT_HTTPS", func() {
			command := "--protocol-uuid uuid-12345 -n test-https -a REDIRECT_HTTPS -r uuid-https-protocol"
			command_args := strings.Fields(command)
			fakeLBManager.AddL7PolicyReturns(datatypes.Network_LBaaS_LoadBalancer{}, nil)
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).NotTo(HaveOccurred())
		})

	})

})

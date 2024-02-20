package loadbal_test

import (
	"errors"
	"strings"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Load balancer edit policies", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *loadbal.L7PolicyEditCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeLBManager *testhelpers.FakeLoadBalancerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = loadbal.NewL7PolicyEditCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cliCommand.LoadBalancerManager = fakeLBManager
	})

	Context("CLI Usage Errors", func() {
		It("Error No policy-id", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("--policy-id"))
		})
		It("No valid action", func() {
			command := "--policy-id 12345 -n test-name -a unknown-action"
			command_args := strings.Fields(command)
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(
				ContainSubstring("-a, --action should be REJECT | REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"))
		})
		It("Error invalid usage for REJECT", func() {
			command := "--policy-id 12345 -n test-name -a REJECT -r REDIRECT_URL"
			command_args := strings.Fields(command)
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(
				"-r, --redirect is only available with action REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"))
		})
		It("Error No --redirect", func() {
			command := "--policy-id 12345 -n test-name -a REDIRECT_URL"
			command_args := strings.Fields(command)
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(
				"-r, --redirect is required with action REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"))
		})
	})

	Context("API Error", func() {
		It("Handles API Error", func() {
			command := "--policy-id 12345 -n test-name -a REJECT"
			command_args := strings.Fields(command)
			fakeLBManager.EditL7PolicyReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("SL_API_ERROR"))
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to edit l7 policy"))
		})
		It("Handles API Error", func() {
			command := "--policy-id 12345 -n test-name -a REJECT"
			command_args := strings.Fields(command)
			fakeLBManager.GetL7PolicyReturns(datatypes.Network_LBaaS_L7Policy{}, errors.New("SL_API_ERROR"))
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			println(err)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to get l7 policy"))
		})
	})

	Context("CLI Usage", func() {
		It("REJECT", func() {
			command := "--policy-id 12345 -n test-reject -a REJECT"
			command_args := strings.Fields(command)
			fakeLBManager.EditL7PolicyReturns(datatypes.Network_LBaaS_LoadBalancer{}, nil)
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).NotTo(HaveOccurred())
		})
		It("REDIRECT_POOL", func() {
			command := "--policy-id 12345 -n test-pool -a REDIRECT_POOL -r uuid-pool"
			command_args := strings.Fields(command)
			fakeLBManager.EditL7PolicyReturns(datatypes.Network_LBaaS_LoadBalancer{}, nil)
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).NotTo(HaveOccurred())
		})
		It("REDIRECT_URL", func() {
			command := "--policy-id 12345 -n test-url -a REDIRECT_URL -r http://example.com"
			command_args := strings.Fields(command)
			fakeLBManager.EditL7PolicyReturns(datatypes.Network_LBaaS_LoadBalancer{}, nil)
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).NotTo(HaveOccurred())
		})
		It("REDIRECT_HTTPS", func() {
			command := "--policy-id 12345 -n test-https -a REDIRECT_HTTPS -r uuid-https-protocol"
			command_args := strings.Fields(command)
			fakeLBManager.EditL7PolicyReturns(datatypes.Network_LBaaS_LoadBalancer{}, nil)
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).NotTo(HaveOccurred())
		})
		It("priority 1 REDIRECT_POOL", func() {
			command := "--policy-id 12345 -n test-https -a REDIRECT_POOL -r uuid-https-protocol -p 1"
			command_args := strings.Fields(command)
			createDateTest, _ := time.Parse(time.RFC3339, "2022-02-01T00:00:00Z")
			fakeLBManager.GetL7PolicyReturns(datatypes.Network_LBaaS_L7Policy{
				Id:          sl.Int(123),
				Uuid:        sl.String("abc123"),
				Name:        sl.String("policy name"),
				Action:      sl.String("REDIRECT_HTTPS"),
				RedirectUrl: sl.String("/"),
				Priority:    sl.Int(1),
				CreateDate:  sl.Time(createDateTest),
			}, nil)
			err := testhelpers.RunCobraCommand(cliCommand.Command, command_args...)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("L7 policy edited"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("ID    UUID     Name          Action           Redirect   Priority   Create Date"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("123   abc123   policy name   REDIRECT_HTTPS   /          1          2022-02-01T00:00:00Z"))
		})
	})
})

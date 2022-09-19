package loadbal_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Load balancer edit policies", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *loadbal.L7RuleAddCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeLBManager *testhelpers.FakeLoadBalancerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = loadbal.NewL7RuleAddCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cliCommand.LoadBalancerManager = fakeLBManager
	})

	Describe("l7 rule add", func() {
		Context("l7 rule add, missing arguments error", func() {
			It("policy-uuid is required", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--policy-uuid' is required"))
			})
			It("-t, --type is required", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--policy-uuid", "abc123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-t, --type' is required"))
			})
			It("-c, --compare-type is required", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--policy-uuid", "abc123", "--type", "HOST_NAME")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-c, --compare-type' is required"))
			})
			It("-v, --value is required", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--policy-uuid", "abc123", "--type", "HOST_NAME", "--compare-type", "EQUAL_TO")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-v, --value' is required"))
			})
		})

		Context("l7 pool edited", func() {
			It("with all attributes", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--policy-uuid", "abc123", "--type", "COOKIE", "--compare-type", "EQUAL_TO", "--value", "value", "--key", "abc123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("L7 rule added"))
			})
		})

		Context("errors", func() {
			It("type wrong", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--policy-uuid", "abc123", "--type", "ABC", "--compare-type", "EQUAL_TO", "--value", "value", "--key", "abc123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("The value of option -t, --type should be HOST_NAME | FILE_TYPE | HEADER | COOKIE | PATH."))
			})
			It("compare-type wrong", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--policy-uuid", "abc123", "--type", "COOKIE", "--compare-type", "ABC", "--value", "value", "--key", "abc123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("The value of option -c, --compare-type should be EQUAL_TO | ENDS_WITH | STARTS_WITH | REGEX | CONTAINS."))
			})
			It("compare-type wrong", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--policy-uuid", "abc123", "--type", "HOST_NAME", "--compare-type", "EQUAL_TO", "--value", "value", "--key", "abc123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("-k, --key is only available in HEADER or COOKIE type."))
			})
			It("Failed to add l7 rule", func() {
				fakeLBManager.AddL7RuleReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--policy-uuid", "abc123", "--type", "COOKIE", "--compare-type", "EQUAL_TO", "--value", "value", "--key", "abc123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to add l7 rule: Internal server error"))
			})
		})
	})
})

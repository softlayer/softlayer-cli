package loadbal_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
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
		cliCommand    *loadbal.ProtocolDeleteCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeLBManager *testhelpers.FakeLoadBalancerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = loadbal.NewProtocolDeleteCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cliCommand.LoadBalancerManager = fakeLBManager
	})

	Describe("protocol del", func() {
		Context("protocol del, missing arguments error", func() {
			It("lb-id is required", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--lb-id' is required"))
			})
			It("protocol-uuid is required", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--lb-id", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--protocol-uuid' is required"))
			})
		})

		Context("protocol del input confirmation error", func() {
			It("Input No, Aborted", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--lb-id", "123", "--protocol-uuid", "abc123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This will delete the load balancer protocol: abc123 and cannot be undone. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted"))
			})
			It("Input wrong, error", func() {
				fakeUI.Inputs("abc")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--lb-id", "123", "--protocol-uuid", "abc123")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This will delete the load balancer protocol: abc123 and cannot be undone. Continue?"))
				Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
			})
		})

		Context("protocol deleted", func() {
			It("with all attributes", func() {
				fakeUI.Inputs("yes")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--lb-id", "123", "--protocol-uuid", "abc123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Protocol abc123 removed"))
			})
		})

		Context("errors", func() {
			It("Failed to get load balancer", func() {
				fakeUI.Inputs("yes")
				fakeLBManager.GetLoadBalancerUUIDReturns("123", errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--lb-id", "123", "--protocol-uuid", "abc123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get load balancer: Internal server error"))
			})
			It("Failed to delete protocol", func() {
				fakeUI.Inputs("yes")
				fakeLBManager.DeleteLoadBalancerListenerReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--lb-id", "123", "--protocol-uuid", "abc123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to delete protocol abc123: Internal server error."))
			})
		})
	})
})

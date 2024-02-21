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

var _ = Describe("Load balancer add policies", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *loadbal.L7PolicyDeleteCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeLBManager *testhelpers.FakeLoadBalancerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = loadbal.NewL7PolicyDeleteCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cliCommand.LoadBalancerManager = fakeLBManager
	})

	Context("policy del without policy-id", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--policy-id' is required"))
		})
	})
	Context("policy del with No as confirmation", func() {
		It("return error", func() {
			fakeUI.Inputs("No")
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--policy-id", "123")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("This will cancel the load balancer policy: 123 and cannot be undone. Continue?"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
		})
	})
	Context("policy del with confirmation error", func() {
		It("return error", func() {
			fakeUI.Inputs("123456")
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--policy-id", "123")
			Expect(err).To(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("This will cancel the load balancer policy: 123 and cannot be undone. Continue?"))
			Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
		})
	})
	Context("policy del with confirmation error", func() {
		It("return error", func() {
			fakeUI.Inputs("Yes")
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--policy-id", "123")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("This will cancel the load balancer policy: 123 and cannot be undone. Continue?"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("L7 policy deleted"))
		})
	})
	Context("policy del with confirmation error", func() {
		BeforeEach(func() {
			fakeLBManager.DeleteL7PolicyReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("Internal server error"))
		})
		It("return error", func() {
			fakeUI.Inputs("Yes")
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--policy-id", "123")
			Expect(err).To(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("This will cancel the load balancer policy: 123 and cannot be undone. Continue?"))
			Expect(err.Error()).To(ContainSubstring("Failed to delete l7 policy: Internal server error."))
		})
	})
})

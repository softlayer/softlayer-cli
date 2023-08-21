package loadbal_test

import (
	"errors"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Load balancer cancel", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *loadbal.CancelCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeLBManager *testhelpers.FakeLoadBalancerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = loadbal.NewCancelCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cliCommand.LoadBalancerManager = fakeLBManager
	})

	Context("cancel without loadbalID", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
		})
	})
	Context("cancel without confirmation", func() {
		It("return aborted", func() {
			fakeUI.Inputs("No")
			err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This will cancel the load balancer: 1234 and cannot be undone. Continue?"}))
		})
	})
	Context("cancel with confirmation error", func() {
		It("return error", func() {
			fakeUI.Inputs("123456")
			err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
			Expect(err).To(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("This will cancel the load balancer: 1234 and cannot be undone. Continue?"))
			Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
		})
	})
	Context("cancel with server fails", func() {
		BeforeEach(func() {
			fakeLBManager.CancelLoadBalancerReturns(true, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to cancel load balancer 1234."))
			Expect(err.Error()).To(ContainSubstring("Internal server error"))
		})
	})
	Context("cancel with correct load balancer ID", func() {
		BeforeEach(func() {
			fakeLBManager.CancelLoadBalancerReturns(true, nil)
		})
		It("return no error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Load balancer 1234 is cancelled."}))
		})
	})
	Context("cancel with server error, load balancer with UUID fail", func() {
		BeforeEach(func() {
			fakeLBManager.GetLoadBalancerUUIDReturns("", errors.New("Internal server error"))
		})
		It("return no error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to get load balancer: Internal server error."))
		})
	})
})

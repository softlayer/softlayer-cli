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

var _ = Describe("Load balancer cancel", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *loadbal.L7MembersAddCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeLBManager *testhelpers.FakeLoadBalancerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = loadbal.NewL7MembersAddCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cliCommand.LoadBalancerManager = fakeLBManager
	})

	Context("member add without pool-uuid", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--pool-uuid' is required"))
		})
	})
	Context("member add without address", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--pool-uuid", "abc123")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--address' is required"))
		})
	})
	Context("member add without port", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--pool-uuid", "abc123", "--address", "address")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--port' is required"))
		})
	})
	Context("member add with server fail", func() {
		BeforeEach(func() {
			fakeLBManager.AddL7MemberReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--pool-uuid", "abc123", "--address", "address", "--port", "123")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to add L7 member: Internal server error."))
		})
	})
	Context("member add Ok", func() {
		It("return no error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--pool-uuid", "abc123", "--address", "address pool", "--port", "123")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("L7 Member address pool added in pool abc123"))
		})
	})
})

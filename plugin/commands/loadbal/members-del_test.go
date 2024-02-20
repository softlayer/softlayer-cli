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
		cliCommand    *loadbal.MembersDelCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeLBManager *testhelpers.FakeLoadBalancerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = loadbal.NewMembersDelCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cliCommand.LoadBalancerManager = fakeLBManager
	})

	Describe("members del", func() {
		Context("members del, missing arguments error", func() {
			It("lb-id is required", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--lb-id' is required"))
			})
			It("member-uuid is required", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--lb-id", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-m, --member-uuid' is required"))
			})
		})

		Context("member del input confirmation error", func() {
			It("Input No, Aborted", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--lb-id", "123", "--member-uuid", "abc123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This will delete the load balancer member: abc123 and cannot be undone. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted"))
			})
			It("Input wrong, error", func() {
				fakeUI.Inputs("abc")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--lb-id", "123", "--member-uuid", "abc123")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This will delete the load balancer member: abc123 and cannot be undone. Continue?"))
				Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
			})
		})

		Context("members deleted", func() {
			It("with all attributes", func() {
				fakeUI.Inputs("yes")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--lb-id", "123", "--member-uuid", "abc123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Member abc123 removed"))
			})
		})

		Context("errors", func() {
			It("Failed to get load balancer", func() {
				fakeUI.Inputs("yes")
				fakeLBManager.GetLoadBalancerUUIDReturns("123", errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--lb-id", "123", "--member-uuid", "abc123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get load balancer: Internal server error"))
			})
			It("Failed to delete load balancer member", func() {
				fakeUI.Inputs("yes")
				fakeLBManager.DeleteLoadBalancerMemberReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--lb-id", "123", "--member-uuid", "abc123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to delete load balancer member abc123: Internal server error."))
			})
		})
	})
})

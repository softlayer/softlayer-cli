package loadbal_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Load balancer edit policies", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeLBManager *testhelpers.FakeLoadBalancerManager
		cmd           *loadbal.L7RuleDelCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewL7RuleDelCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        loadbal.LoadbalL7RuleDelMetadata().Name,
			Description: loadbal.LoadbalL7RuleDelMetadata().Description,
			Usage:       loadbal.LoadbalL7RuleDelMetadata().Usage,
			Flags:       loadbal.LoadbalL7RuleDelMetadata().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("l7 rule del", func() {
		Context("l7 rule del, missing arguments error", func() {
			It("policy-uuid is required", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--policy-uuid' is required"))
			})
			It("rule-uuid is required", func() {
				err := testhelpers.RunCommand(cliCommand, "--policy-uuid", "abc123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--rule-uuid' is required"))
			})
		})

		Context("l7 rule del input confirmation error", func() {
			It("Input No, Aborted", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "--policy-uuid", "abc123", "--rule-uuid", "abcd1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This will delete the load balancer L7 rule: abcd1234 and cannot be undone. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted"))
			})
			It("Input No, Aborted", func() {
				fakeUI.Inputs("abc")
				err := testhelpers.RunCommand(cliCommand, "--policy-uuid", "abc123", "--rule-uuid", "abcd1234")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This will delete the load balancer L7 rule: abcd1234 and cannot be undone. Continue?"))
				Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
			})
		})

		Context("l7 rule deleted", func() {
			It("with all attributes", func() {
				fakeUI.Inputs("yes")
				err := testhelpers.RunCommand(cliCommand, "--policy-uuid", "abc123", "--rule-uuid", "abcd1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("L7Rule abcd1234 removed"))
			})
		})

		Context("errors", func() {
			It("Failed to del l7 rule", func() {
				fakeUI.Inputs("yes")
				fakeLBManager.DeleteL7RuleReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("Internal server error"))
				err := testhelpers.RunCommand(cliCommand, "--policy-uuid", "abc123", "--rule-uuid", "abcd1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to delete L7Rule abcd1234: Internal server error"))
			})
		})
	})
})

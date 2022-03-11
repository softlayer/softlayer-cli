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

var _ = Describe("Load balancer cancel", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeLBManager *testhelpers.FakeLoadBalancerManager
		cmd           *loadbal.L7MembersDelCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewL7MembersDelCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        loadbal.LoadbalL7MemberDeleteMetadata().Name,
			Description: loadbal.LoadbalL7MemberDeleteMetadata().Description,
			Usage:       loadbal.LoadbalL7MemberDeleteMetadata().Usage,
			Flags:       loadbal.LoadbalL7MemberDeleteMetadata().Flags,
			Action:      cmd.Run,
		}
	})

	Context("member del without pool-uuid", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--pool-uuid' is required"))
		})
	})
	Context("member del without member-uuid", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "--pool-uuid", "abc123")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--member-uuid' is required"))
		})
	})
	Context("member del with No as confirmation", func() {
		It("return error", func() {
			fakeUI.Inputs("No")
			err := testhelpers.RunCommand(cliCommand, "--pool-uuid", "abc123", "--member-uuid", "abcde123456")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("This will delete the load balancer L7 member: abcde123456 and cannot be undone. Continue?"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
		})
	})
	Context("member del with confirmation error", func() {
		It("return error", func() {
			fakeUI.Inputs("123456")
			err := testhelpers.RunCommand(cliCommand, "--pool-uuid", "abc123", "--member-uuid", "abcde123456")
			Expect(err).To(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("This will delete the load balancer L7 member: abcde123456 and cannot be undone. Continue?"))
			Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
		})
	})
	Context("member del OK", func() {
		It("return error", func() {
			fakeUI.Inputs("Yes")
			err := testhelpers.RunCommand(cliCommand, "--pool-uuid", "abc123", "--member-uuid", "abcde123456")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("This will delete the load balancer L7 member: abcde123456 and cannot be undone. Continue?"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Member abcde123456 removed from abc123"))
		})
	})
	Context("member del with server fail", func() {
		BeforeEach(func() {
			fakeLBManager.DeleteL7MemberReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("Internal server error"))
		})
		It("return error", func() {
			fakeUI.Inputs("Yes")
			err := testhelpers.RunCommand(cliCommand, "--pool-uuid", "abc123", "--member-uuid", "abcde123456")
			Expect(err).To(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("This will delete the load balancer L7 member: abcde123456 and cannot be undone. Continue?"))
			Expect(err.Error()).To(ContainSubstring("Failed to delete L7member abcde123456: Internal server error."))
		})
	})
})

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

var _ = Describe("Load balancer add policies", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeLBManager *testhelpers.FakeLoadBalancerManager
		cmd           *loadbal.L7PolicyDeleteCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewL7PolicyDeleteCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        loadbal.LoadbalL7PolicyDeleteMetadata().Name,
			Description: loadbal.LoadbalL7PolicyDeleteMetadata().Description,
			Usage:       loadbal.LoadbalL7PolicyDeleteMetadata().Usage,
			Flags:       loadbal.LoadbalL7PolicyDeleteMetadata().Flags,
			Action:      cmd.Run,
		}
	})

	Context("policy del without policy-id", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--policy-id' is required"))
		})
	})
	Context("policy del with No as confirmation", func() {
		It("return error", func() {
			fakeUI.Inputs("No")
			err := testhelpers.RunCommand(cliCommand, "--policy-id", "123")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("This will cancel the load balancer policy: 123 and cannot be undone. Continue?"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
		})
	})
	Context("policy del with confirmation error", func() {
		It("return error", func() {
			fakeUI.Inputs("123456")
			err := testhelpers.RunCommand(cliCommand, "--policy-id", "123")
			Expect(err).To(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("This will cancel the load balancer policy: 123 and cannot be undone. Continue?"))
			Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
		})
	})
	Context("policy del with confirmation error", func() {
		It("return error", func() {
			fakeUI.Inputs("Yes")
			err := testhelpers.RunCommand(cliCommand, "--policy-id", "123")
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
			err := testhelpers.RunCommand(cliCommand, "--policy-id", "123")
			Expect(err).To(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("This will cancel the load balancer policy: 123 and cannot be undone. Continue?"))
			Expect(err.Error()).To(ContainSubstring("Failed to delete l7 policy: Internal server error."))
		})
	})
})

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
		cmd           *loadbal.L7MembersAddCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewL7MembersAddCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        loadbal.LoadbalL7MemberAddMetadata().Name,
			Description: loadbal.LoadbalL7MemberAddMetadata().Description,
			Usage:       loadbal.LoadbalL7MemberAddMetadata().Usage,
			Flags:       loadbal.LoadbalL7MemberAddMetadata().Flags,
			Action:      cmd.Run,
		}
	})

	Context("member add without pool-uuid", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--pool-uuid' is required"))
		})
	})
	Context("member add without address", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "--pool-uuid", "abc123")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--address' is required"))
		})
	})
	Context("member add without port", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "--pool-uuid", "abc123", "--address", "address")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--port' is required"))
		})
	})
	Context("member add with server fail", func() {
		BeforeEach(func() {
			fakeLBManager.AddL7MemberReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "--pool-uuid", "abc123", "--address", "address", "--port", "123")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to add L7 member: Internal server error."))
		})
	})
	Context("member add Ok", func() {
		It("return no error", func() {
			err := testhelpers.RunCommand(cliCommand, "--pool-uuid", "abc123", "--address", "address pool", "--port", "123")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("L7 Member address pool added in pool abc123"))
		})
	})
})

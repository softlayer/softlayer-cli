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
		cmd           *loadbal.MembersAddCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewMembersAddCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        loadbal.LoadbalMemberAddMetadata().Name,
			Description: loadbal.LoadbalMemberAddMetadata().Description,
			Usage:       loadbal.LoadbalMemberAddMetadata().Usage,
			Flags:       loadbal.LoadbalMemberAddMetadata().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("members add", func() {
		Context("members add, missing arguments error", func() {
			It("id is required", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--id' is required"))
			})
			It("ip is required", func() {
				err := testhelpers.RunCommand(cliCommand, "--id", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--ip' is required"))
			})
		})

		Context("members added", func() {
			It("with all attributes", func() {
				err := testhelpers.RunCommand(cliCommand, "--id", "123", "--ip", "10.10.10.10")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Member 10.10.10.10 added"))
			})
		})

		Context("errors", func() {
			It("Failed to get load balancer", func() {
				fakeLBManager.GetLoadBalancerUUIDReturns("123", errors.New("Internal server error"))
				err := testhelpers.RunCommand(cliCommand, "--id", "123", "--ip", "10.10.10.10")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get load balancer: Internal server error"))
			})
			It("Failed to add load balancer member", func() {
				fakeLBManager.AddLoadBalancerMemberReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("Internal server error"))
				err := testhelpers.RunCommand(cliCommand, "--id", "123", "--ip", "10.10.10.10")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to add load balancer member: Internal server error."))
			})
		})
	})
})

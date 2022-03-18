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
		cmd           *loadbal.L7PoolDelCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewL7PoolDelCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        loadbal.LoadbalL7PoolDelMetadata().Name,
			Description: loadbal.LoadbalL7PoolDelMetadata().Description,
			Usage:       loadbal.LoadbalL7PoolDelMetadata().Usage,
			Flags:       loadbal.LoadbalL7PoolDelMetadata().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("l7 pool del", func() {
		Context("l7 pool del, missing arguments error", func() {
			It("--id is required", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--pool-id' is required"))
			})
		})

		Context("l7 pool input confirmation error", func() {
			It("Input No, Aborted", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "--pool-id", "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This will delete the load balancer L7 pool: 123 and cannot be undone. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted"))
			})
			It("Input No, Aborted", func() {
				fakeUI.Inputs("abc")
				err := testhelpers.RunCommand(cliCommand, "--pool-id", "123")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This will delete the load balancer L7 pool: 123 and cannot be undone. Continue?"))
				Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
			})

		})

		Context("l7 pool deleted", func() {
			It("with all attributes", func() {
				fakeUI.Inputs("yes")
				err := testhelpers.RunCommand(cliCommand, "--pool-id", "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("L7Pool 123 removed"))
			})

		})

		Context("errors", func() {
			It("failed to delete l7pool", func() {
				fakeUI.Inputs("yes")
				fakeLBManager.DeleteLoadBalancerL7PoolReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("Internal server error"))
				err := testhelpers.RunCommand(cliCommand, "--pool-id", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to delete L7Pool 123: Internal server error."))
			})
		})
	})
})

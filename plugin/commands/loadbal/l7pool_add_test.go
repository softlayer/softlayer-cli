package loadbal_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
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
		cliCommand    *loadbal.L7PoolAddCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeLBManager *testhelpers.FakeLoadBalancerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = loadbal.NewL7PoolAddCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cliCommand.LoadBalancerManager = fakeLBManager
	})

	Describe("l7 pool add", func() {
		Context("l7 pool add, missing arguments error", func() {
			It("--id is required", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--id' is required"))
			})
			It("--name is required", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--id", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-n, --name' is required"))
			})
		})

		Context("l7 pool added", func() {
			It("with all attributes", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--id", "123", "--name", "poolTest", "--method", "", "--protocol", "", "-s", "10.0.0.1:80", "--health-path", "", "--health-interval", "0", "--health-retry", "0", "--health-timeout", "0", "--sticky", "cookie")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("L7 pool added"))
			})
			It("with sticky as source-ip", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--id", "123", "--name", "poolTest", "--method", "", "--protocol", "", "-s", "10.0.0.1:80", "--health-path", "", "--health-interval", "0", "--health-retry", "0", "--health-timeout", "0", "--sticky", "source-ip")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("L7 pool added"))
			})
		})

		Context("errors", func() {
			It("get load balancer error", func() {
				fakeLBManager.GetLoadBalancerUUIDReturns("", errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--id", "123", "--name", "poolTest", "--method", "", "--protocol", "", "-s", "10.0.0.1:80", "--health-path", "", "--health-interval", "0", "--health-retry", "0", "--health-timeout", "0", "--sticky", "cookie")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get load balancer: Internal server error."))
			})
			It("sticky wrong", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--id", "123", "--name", "poolTest", "--method", "", "--protocol", "", "-s", "10.0.0.1:80", "--health-path", "", "--health-interval", "0", "--health-retry", "0", "--health-timeout", "0", "--sticky", "abcd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Value of option '--sticky' should be cookie or source-ip"))
			})
			It("add load balancer l7 pool error", func() {
				fakeLBManager.AddLoadBalancerL7PoolReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--id", "123", "--name", "poolTest", "--method", "", "--protocol", "", "-s", "10.0.0.1:80", "--health-path", "", "--health-interval", "0", "--health-retry", "0", "--health-timeout", "0", "--sticky", "cookie")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to add load balancer l7 pool: Internal server error."))
			})
			It("options server needs a port", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--id", "123", "--name", "poolTest", "--method", "", "--protocol", "", "-s", "10.0.0.1", "--health-path", "", "--health-interval", "0", "--health-retry", "0", "--health-timeout", "0", "--sticky", "cookie")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("--server needs a port. 10.0.0.1 improperly formatted"))
			})
			It("add load balancer, options server needs a port", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--id", "123", "--name", "poolTest", "--method", "", "--protocol", "", "-s", "10.0.0.1:abc", "--health-path", "", "--health-interval", "0", "--health-retry", "0", "--health-timeout", "0", "--sticky", "cookie")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("The port has to be a positive integer."))
			})
		})
	})
})

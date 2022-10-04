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
		cliCommand    *loadbal.L7PoolEditCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeLBManager *testhelpers.FakeLoadBalancerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = loadbal.NewL7PoolEditCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cliCommand.LoadBalancerManager = fakeLBManager
	})

	Describe("l7 pool edit", func() {
		Context("l7 pool edit, missing arguments error", func() {
			It("--pool-uuid is required", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--pool-uuid' is required"))
			})
			It("pass at least one of the flags.", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--pool-uuid", "abc123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Please pass at least one of the flags."))
			})
		})

		Context("l7 pool edited", func() {
			It("with all attributes", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--pool-uuid", "123", "--name", "NewName", "--method", "ROUNDROBIN", "--protocol", "80", "-s", "10.0.0.1:80", "--health-path", "/", "--health-interval", "100", "--health-retry", "5", "--health-timeout", "200", "--sticky", "cookie")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("L7 pool updated"))
			})
			It("with sticky as source-ip", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--pool-uuid", "123", "--method", "ROUNDROBIN", "--protocol", "80", "-s", "10.0.0.1:80", "--health-path", "/", "--health-interval", "100", "--health-retry", "5", "--health-timeout", "200", "--sticky", "source-ip")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("L7 pool updated"))
			})
		})

		Context("errors", func() {
			It("sticky wrong", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--pool-uuid", "123", "--method", "ROUNDROBIN", "--protocol", "80", "-s", "10.0.0.1:80", "--health-path", "/", "--health-interval", "100", "--health-retry", "5", "--health-timeout", "200", "--sticky", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Value of option '--sticky' should be cookie or source-ip"))
			})
			It("options server needs a port", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--pool-uuid", "123", "--method", "ROUNDROBIN", "--protocol", "80", "-s", "10.0.0.1", "--health-path", "/", "--health-interval", "100", "--health-retry", "5", "--health-timeout", "200", "--sticky", "cookie")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("--server needs a port. 10.0.0.1 improperly formatted"))
			})
			It("add load balancer l7 pool error", func() {
				fakeLBManager.UpdateLoadBalancerL7PoolReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--pool-uuid", "123", "--method", "ROUNDROBIN", "--protocol", "80", "-s", "10.0.0.1:80", "--health-path", "/", "--health-interval", "100", "--health-retry", "5", "--health-timeout", "200", "--sticky", "cookie")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to update l7 pool: Internal server error."))
			})
		})
	})
})

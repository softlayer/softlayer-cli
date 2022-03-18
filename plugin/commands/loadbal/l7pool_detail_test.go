package loadbal_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Load balancer edit policies", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeLBManager *testhelpers.FakeLoadBalancerManager
		cmd           *loadbal.L7PoolDetailCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewL7PoolDetailCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        loadbal.LoadbalL7PoolDetailMetadata().Name,
			Description: loadbal.LoadbalL7PoolDetailMetadata().Description,
			Usage:       loadbal.LoadbalL7PoolDetailMetadata().Usage,
			Flags:       loadbal.LoadbalL7PoolDetailMetadata().Flags,
			Action:      cmd.Run,
		}
		modifyDateTest, _ := time.Parse(time.RFC3339, "2022-02-01T00:00:00Z")
		fakeLBManager.GetLoadBalancerL7PoolReturns(datatypes.Network_LBaaS_L7Pool{
			Name:                   sl.String("Name"),
			Id:                     sl.Int(123),
			Uuid:                   sl.String("abc123"),
			LoadBalancingAlgorithm: sl.String("REDIRECT_URL"),
			Protocol:               sl.String("80"),
		}, nil)
		fakeLBManager.GetL7SessionAffinityReturns(datatypes.Network_LBaaS_L7SessionAffinity{
			Type: sl.String("type"),
		}, nil)
		fakeLBManager.GetL7HealthMonitorReturns(datatypes.Network_LBaaS_L7HealthMonitor{
			Interval:           sl.Int(100),
			MaxRetries:         sl.Int(5),
			MonitorType:        sl.String("Monitor"),
			Timeout:            sl.Int(150),
			UrlPath:            sl.String("/"),
			ModifyDate:         sl.Time(modifyDateTest),
			ProvisioningStatus: sl.String("ACTIVE"),
		}, nil)
		fakeLBManager.ListL7MembersReturns([]datatypes.Network_LBaaS_L7Member{
			datatypes.Network_LBaaS_L7Member{
				Id:                 sl.Int(12345),
				Uuid:               sl.String("abcd1234"),
				Address:            sl.String("Address member"),
				Weight:             sl.Int(100),
				ModifyDate:         sl.Time(modifyDateTest),
				ProvisioningStatus: sl.String("ACTIVE"),
			},
		}, nil)
	})

	Describe("l7 pool detail", func() {
		Context("l7 pool detail, missing arguments error", func() {
			It("--pool-id is required", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--pool-id' is required"))
			})
		})

		Context("l7 pool deleted", func() {
			It("with all attributes", func() {
				err := testhelpers.RunCommand(cliCommand, "--pool-id", "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name                 Value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name                 Name"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID                   123"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("UUID                 abc123"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Method               REDIRECT_URL"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Protocol             80"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Session Stickiness   type"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Health Check:        Interval   Retries   Type      Timeout   URL   Modify                 Active"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("100        5         Monitor   150       /     2022-02-01T00:00:00Z   ACTIVE"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Members:             ID      UUID       Address          Weight   Modify                 Active"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("12345   abcd1234   Address member   100      2022-02-01T00:00:00Z   ACTIVE"))
			})
		})

		Context("errors", func() {
			It("failed to get the detail to l7pool", func() {
				fakeUI.Inputs("yes")
				fakeLBManager.GetLoadBalancerL7PoolReturns(datatypes.Network_LBaaS_L7Pool{}, errors.New("Internal server error"))
				err := testhelpers.RunCommand(cliCommand, "--pool-id", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get L7 Pool 123: Internal server error."))
			})
			It("Failed to get L7 Pool Session Affinity", func() {
				fakeUI.Inputs("yes")
				fakeLBManager.GetL7SessionAffinityReturns(datatypes.Network_LBaaS_L7SessionAffinity{}, errors.New("Internal server error"))
				err := testhelpers.RunCommand(cliCommand, "--pool-id", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get L7 Pool Session Affinity: Internal server error."))
			})
			It("Failed to get L7 Health Monitor", func() {
				fakeUI.Inputs("yes")
				fakeLBManager.GetL7HealthMonitorReturns(datatypes.Network_LBaaS_L7HealthMonitor{}, errors.New("Internal server error"))
				err := testhelpers.RunCommand(cliCommand, "--pool-id", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get L7 Health Monitor: Internal server error."))
			})
			It("Failed to get L7 Members", func() {
				fakeUI.Inputs("yes")
				fakeLBManager.ListL7MembersReturns([]datatypes.Network_LBaaS_L7Member{}, errors.New("Internal server error"))
				err := testhelpers.RunCommand(cliCommand, "--pool-id", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get L7 Members: Internal server error."))
			})
		})
	})
})

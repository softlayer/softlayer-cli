package loadbal_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Load balancer health", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *loadbal.HealthChecksCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeLBManager *testhelpers.FakeLoadBalancerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = loadbal.NewHealthChecksCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cliCommand.LoadBalancerManager = fakeLBManager

		modifyDateTest, _ := time.Parse(time.RFC3339, "2022-02-01T00:00:00Z")
		fakeLBManager.GetLoadBalancerReturns(datatypes.Network_LBaaS_LoadBalancer{
			Listeners: []datatypes.Network_LBaaS_Listener{
				datatypes.Network_LBaaS_Listener{
					DefaultPool: &datatypes.Network_LBaaS_Pool{
						HealthMonitor: &datatypes.Network_LBaaS_HealthMonitor{
							Uuid: sl.String("3f1111fe-c666-4ca4-9ded-6c66d6c6aef6"),
						},
					},
				},
			},
		}, nil)
		fakeLBManager.UpdateLBHealthMonitorsReturns(datatypes.Network_LBaaS_LoadBalancer{
			Id:      sl.Int(123456),
			Uuid:    sl.String("3f1111fe-c666-4ca4-9ded-6c66d6c6aef6"),
			Name:    sl.String("test-lb"),
			Address: sl.String("test.domain.cloud"),
			Type:    sl.Int(0),
			Datacenter: &datatypes.Location{
				LongName: sl.String("Mexico-1"),
			},
			Description:        sl.String("test lb description"),
			ProvisioningStatus: sl.String("ACTIVE"),
			OperatingStatus:    sl.String("ONLINE"),
			Listeners: []datatypes.Network_LBaaS_Listener{
				datatypes.Network_LBaaS_Listener{
					DefaultPool: &datatypes.Network_LBaaS_Pool{
						Protocol:               sl.String("HTTP"),
						ProtocolPort:           sl.Int(80),
						Uuid:                   sl.String("acbde123456789"),
						LoadBalancingAlgorithm: sl.String("ROUNDROBIN"),
					},
					Id:                 sl.Int(654321),
					Protocol:           sl.String("HTTPS"),
					ProtocolPort:       sl.Int(81),
					Uuid:               sl.String("3f1111fe-c666-4ca4-9ded-6c66d6c6aef6"),
					ConnectionLimit:    sl.Int(10),
					ClientTimeout:      sl.Int(100),
					ServerTimeout:      sl.Int(50),
					ModifyDate:         sl.Time(modifyDateTest),
					ProvisioningStatus: sl.String("ACTIVE"),
				},
			},
			Members: []datatypes.Network_LBaaS_Member{
				datatypes.Network_LBaaS_Member{
					Id:                 sl.Int(789),
					Uuid:               sl.String("3f1111fe-c666-4ca4-9ded-6c66d6c6aef6"),
					Address:            sl.String("Member address"),
					ModifyDate:         sl.Time(modifyDateTest),
					ProvisioningStatus: sl.String("ACTIVE"),
				},
			},
			HealthMonitors: []datatypes.Network_LBaaS_HealthMonitor{
				datatypes.Network_LBaaS_HealthMonitor{
					Id:                 sl.Int(987),
					Uuid:               sl.String("3f1111fe-c666-4ca4-9ded-6c66d6c6aef6"),
					MonitorType:        sl.String("HTTP"),
					Interval:           sl.Int(5),
					MaxRetries:         sl.Int(2),
					Timeout:            sl.Int(100),
					UrlPath:            sl.String("/"),
					ModifyDate:         sl.Time(modifyDateTest),
					ProvisioningStatus: sl.String("ACTIVE"),
				},
			},
			L7Pools: []datatypes.Network_LBaaS_L7Pool{
				datatypes.Network_LBaaS_L7Pool{
					Id:                     sl.Int(753),
					Uuid:                   sl.String("3f1111fe-c666-4ca4-9ded-6c66d6c6aef6"),
					Name:                   sl.String("L7Pool Name"),
					Protocol:               sl.String("HTTP"),
					LoadBalancingAlgorithm: sl.String("ROUNDROBIN"),
					ModifyDate:             sl.Time(modifyDateTest),
					ProvisioningStatus:     sl.String("ACTIVE"),
				},
			},
		}, nil)
	})
	Context("health without --lb-id", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
		})
	})
	Context("health without --health-uuid", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--health-uuid' is required"))
		})
	})
	Context("health without --health-uuid", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--health-uuid", "3f1111fe-c666-4ca4-9ded-6c66d6c6aef6")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: At least one of these flags is required :-i, --interval, -r, --retry, -t, --timeout,  -u, --url,"))
		})
	})
	Context("health with Unknow --health-uuid", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "1111", "--health-uuid", "abcde", "-i", "10")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Unable to find health check with UUID of 'abcde' in load balancer 1111."))
		})
	})
	Context("health with exist ID, --health-uuid, interval, retry, timeout and url", func() {
		It("return no error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--health-uuid", "3f1111fe-c666-4ca4-9ded-6c66d6c6aef6", "-i", "5", "-r", "2", "-t", "100", "-u", "/")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("Name            Value"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("ID              123456"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("UUID            3f1111fe-c666-4ca4-9ded-6c66d6c6aef6"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Name            test-lb"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Address         test.domain.cloud"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Type            Private to Private"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Location        Mexico-1"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Description     test lb description"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Status          ACTIVE/ONLINE"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Protocols:      ID       UUID                                   Mapping               Method       Max Connection   Timeout                     Modify                 Active"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("654321   3f1111fe-c666-4ca4-9ded-6c66d6c6aef6   HTTPS:81 -> HTTP:80   ROUNDROBIN   10               Client: 100s, Server: 50s   2022-02-01T00:00:00Z   ACTIVE"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Members:        ID    UUID                                   Address          Modify                 Active"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("789   3f1111fe-c666-4ca4-9ded-6c66d6c6aef6   Member address   2022-02-01T00:00:00Z   ACTIVE"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Health Check:   ID    UUID                                   Protocol   Interval   Retries   Timeout   URL   Modify                 Active"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("987   3f1111fe-c666-4ca4-9ded-6c66d6c6aef6   HTTP       5          2         100       /     2022-02-01T00:00:00Z   ACTIVE"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("L7 Pools:       ID    UUID                                   Name          Protocol   Method       Modify Date            ProvisioningStatus"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("753   3f1111fe-c666-4ca4-9ded-6c66d6c6aef6   L7Pool Name   HTTP       ROUNDROBIN   2022-02-01T00:00:00Z   ACTIVE"))
		})
	})
	Context("health with server fails, get load balancer fail", func() {
		BeforeEach(func() {
			fakeLBManager.GetLoadBalancerReturns(datatypes.Network_LBaaS_LoadBalancer{},
				errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--health-uuid", "3f1111fe-c666-4ca4-9ded-6c66d6c6aef6", "-i", "5")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to get load balancer: Internal server error."))
		})
	})
	Context("health with server fails, update load balancer fail", func() {
		BeforeEach(func() {
			fakeLBManager.UpdateLBHealthMonitorsReturns(datatypes.Network_LBaaS_LoadBalancer{},
				errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--health-uuid", "3f1111fe-c666-4ca4-9ded-6c66d6c6aef6", "-i", "5")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to update health check"))
			Expect(err.Error()).To(ContainSubstring("Internal server error"))
		})
	})
	Context("health with protocol like a TCP", func() {
		BeforeEach(func() {
			fakeLBManager.GetLoadBalancerReturns(datatypes.Network_LBaaS_LoadBalancer{
				Listeners: []datatypes.Network_LBaaS_Listener{
					datatypes.Network_LBaaS_Listener{
						DefaultPool: &datatypes.Network_LBaaS_Pool{
							HealthMonitor: &datatypes.Network_LBaaS_HealthMonitor{
								Uuid: sl.String("3f1111fe-c666-4ca4-9ded-6c66d6c6aef6"),
							},
							Protocol: sl.String("TCP"),
						},
					},
				},
			}, nil)
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--health-uuid", "3f1111fe-c666-4ca4-9ded-6c66d6c6aef6", "-u", "HTTP")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("--url cannot be used with TCP checks."))
		})
	})
})

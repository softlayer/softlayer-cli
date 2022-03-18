package loadbal_test

import (
	"errors"
	"strings"
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

var _ = Describe("Load balancer detail", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeLBManager *testhelpers.FakeLoadBalancerManager
		cmd           *loadbal.DetailCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewDetailCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        loadbal.LoadbalDetailMetadata().Name,
			Description: loadbal.LoadbalDetailMetadata().Description,
			Usage:       loadbal.LoadbalDetailMetadata().Usage,
			Flags:       loadbal.LoadbalDetailMetadata().Flags,
			Action:      cmd.Run,
		}
		modifyDateTest, _ := time.Parse(time.RFC3339, "2022-02-01T00:00:00Z")
		fakeLBManager.GetLoadBalancerReturns(
			datatypes.Network_LBaaS_LoadBalancer{
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
							Uuid:                   sl.String("UUID Default"),
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
	Context("detail without loadbalID", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "'--id' is required")).To(BeTrue())
		})
	})
	Context("detail with wrong loadbalID", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "--id", "abc")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "invalid value")).To(BeTrue())
		})
	})
	Context("detail with server fails", func() {
		BeforeEach(func() {
			fakeLBManager.GetLoadBalancerReturns(datatypes.Network_LBaaS_LoadBalancer{},
				errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "--id", "1234")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to get load balancer with ID 1234.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("detail with loadbal ID", func() {
		It("return loadbalancer", func() {
			err := testhelpers.RunCommand(cliCommand, "--id", "123456")
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
	Context("detail loadbal with type as Public to Private", func() {
		BeforeEach(func() {
			fakeLBManager.GetLoadBalancerReturns(datatypes.Network_LBaaS_LoadBalancer{
				Type: sl.Int(1),
			}, nil)
		})
		It("return loadbalancer", func() {
			err := testhelpers.RunCommand(cliCommand, "--id", "123456")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("Type            Public to Private"))
		})
	})
	Context("detail loadbal with type as Public to Public", func() {
		BeforeEach(func() {
			fakeLBManager.GetLoadBalancerReturns(datatypes.Network_LBaaS_LoadBalancer{
				Type: sl.Int(2),
			}, nil)
		})
		It("return loadbalancer", func() {
			err := testhelpers.RunCommand(cliCommand, "--id", "123456")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("Type            Public to Public"))
		})
	})
	Context("detail loadbal without Protocols, Members, Health Check and L7 Pools", func() {
		BeforeEach(func() {
			fakeLBManager.GetLoadBalancerReturns(datatypes.Network_LBaaS_LoadBalancer{
				Id: sl.Int(987654),
			}, nil)
		})
		It("return Not Founds", func() {
			err := testhelpers.RunCommand(cliCommand, "--id", "987654")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("Protocols:      Not Found"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Members:        Not Found"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Health Check:   Not Found"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("L7 Pools:       Not Found"))
		})
	})
})

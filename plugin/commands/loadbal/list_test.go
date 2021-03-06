package loadbal_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Load balancer list", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeLBManager *testhelpers.FakeLoadBalancerManager
		cmd           *loadbal.ListCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewListCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        loadbal.LoadbalListMetadata().Name,
			Description: loadbal.LoadbalListMetadata().Description,
			Usage:       loadbal.LoadbalListMetadata().Usage,
			Flags:       loadbal.LoadbalListMetadata().Flags,
			Action:      cmd.Run,
		}
	})

	Context("list with server fails", func() {
		BeforeEach(func() {
			fakeLBManager.GetLoadBalancersReturns(nil, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to get load balancers on your account.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("list", func() {
		BeforeEach(func() {
			fakeLBManager.GetLoadBalancersReturns(nil, nil)
		})
		It("return no load balancer", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"No load balancer was found."}))
		})
	})
	Context("list without location", func() {
		BeforeEach(func() {
			address := "address"
			desc := "desc"
			isPublic := 1
			fakeLBManager.GetLoadBalancersReturns([]datatypes.Network_LBaaS_LoadBalancer{
				datatypes.Network_LBaaS_LoadBalancer{
					Id:          sl.Int(13162),
					Address:     &address,
					Description: &desc,
					Type:        &isPublic,
				},
			}, nil)
		})
		It("return loadbalancer", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"13162"}))
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"address"}))
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Public to Private"}))
		})
	})
	Context("list with location", func() {
		BeforeEach(func() {
			address := "address"
			desc := "desc"
			isPublic := 1
			longName := "dal05"
			fakeLBManager.GetLoadBalancersReturns([]datatypes.Network_LBaaS_LoadBalancer{
				datatypes.Network_LBaaS_LoadBalancer{
					Id:          sl.Int(13162),
					Address:     &address,
					Description: &desc,
					Type:        &isPublic,
					Datacenter: &datatypes.Location{
						LongName: &longName,
					},
				},
			}, nil)
		})
		It("return loadbalancer", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"13162"}))
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"address"}))
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Public to Private"}))
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"dal05"}))
		})
	})

	Context("list with location and type as private to private", func() {
		BeforeEach(func() {
			address := "address"
			desc := "desc"
			isPublic := 0
			longName := "dal05"
			fakeLBManager.GetLoadBalancersReturns([]datatypes.Network_LBaaS_LoadBalancer{
				datatypes.Network_LBaaS_LoadBalancer{
					Id:          sl.Int(13162),
					Address:     &address,
					Description: &desc,
					Type:        &isPublic,
					Datacenter: &datatypes.Location{
						LongName: &longName,
					},
				},
			}, nil)
		})
		It("return loadbalancer", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("ID      Name   Address   Type                 Location   Create Date   Status"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("13162   -      address   Private to Private   dal05      -             -/-"))
		})
	})

	Context("list with location and type as Public to Public", func() {
		BeforeEach(func() {
			address := "address"
			desc := "desc"
			isPublic := 2
			longName := "dal05"
			fakeLBManager.GetLoadBalancersReturns([]datatypes.Network_LBaaS_LoadBalancer{
				datatypes.Network_LBaaS_LoadBalancer{
					Id:          sl.Int(13162),
					Address:     &address,
					Description: &desc,
					Type:        &isPublic,
					Datacenter: &datatypes.Location{
						LongName: &longName,
					},
				},
			}, nil)
		})
		It("return loadbalancer", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("ID      Name   Address   Type               Location   Create Date   Status"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("13162   -      address   Public to Public   dal05      -             -/-"))
		})
	})
})

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
		cmd           *loadbal.NetscalerListCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewNetscalerListCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        loadbal.LoadbalNsListMetadata().Name,
			Description: loadbal.LoadbalNsListMetadata().Description,
			Usage:       loadbal.LoadbalNsListMetadata().Usage,
			Flags:       loadbal.LoadbalNsListMetadata().Flags,
			Action:      cmd.Run,
		}

		createdDate, _ := time.Parse(time.RFC3339, "2016-12-29T00:00:00Z")
		fakeLBManager.GetADCsReturns([]datatypes.Network_Application_Delivery_Controller{
			datatypes.Network_Application_Delivery_Controller{
				Datacenter: &datatypes.Location{
					LongName: sl.String("dal01"),
				},
				Id:                           sl.Int(123),
				Name:                         sl.String("Netscaler Name"),
				Description:                  sl.String("Description Netscaler"),
				PrimaryIpAddress:             sl.String("10.10.10.10"),
				ManagementIpAddress:          sl.String("20.20.20.20"),
				OutboundPublicBandwidthUsage: sl.Float(2.0),
				CreateDate:                   sl.Time(createdDate),
			},
			datatypes.Network_Application_Delivery_Controller{
				Datacenter: &datatypes.Location{
					LongName: sl.String("dal02"),
				},
				Id:                           sl.Int(1234),
				Name:                         sl.String("Netscaler Name 2"),
				Description:                  sl.String("Description Netscaler 2"),
				PrimaryIpAddress:             sl.String("10.10.10.11"),
				ManagementIpAddress:          sl.String("20.20.20.21"),
				OutboundPublicBandwidthUsage: sl.Float(3.0),
				CreateDate:                   sl.Time(createdDate),
			},
		}, nil)
	})

	Describe("ns list", func() {
		Context("ns details, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCommand(cliCommand, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("ns list", func() {
			It("list all netscalers", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID     Location   Name               Description               IP Address    Management IP   Bandwidth   Create Date"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123    dal01      Netscaler Name     Description Netscaler     10.10.10.10   20.20.20.20     2.000000    2016-12-29T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234   dal02      Netscaler Name 2   Description Netscaler 2   10.10.10.11   20.20.20.21     3.000000    2016-12-29T00:00:00Z"))
			})

			It("list all netscalers in output json", func() {
				err := testhelpers.RunCommand(cliCommand, "--output=json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"createDate": "2016-12-29T00:00:00Z",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"datacenter": {`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"longName": "dal01"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`},`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"description": "Description Netscaler",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"id": 123,`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"managementIpAddress": "20.20.20.20",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"name": "Netscaler Name",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"outboundPublicBandwidthUsage": 2,`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"primaryIpAddress": "10.10.10.10"`))
			})
		})

		Context("errors", func() {
			It("Failed to get netscalers on your account.", func() {
				fakeLBManager.GetADCsReturns([]datatypes.Network_Application_Delivery_Controller{}, errors.New("Internal server error"))
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get netscalers on your account.Internal server error"))
			})
			It("Failed to get netscalers on your account.", func() {
				fakeLBManager.GetADCsReturns([]datatypes.Network_Application_Delivery_Controller{}, nil)
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("No netscalers was found."))
			})
		})
	})
})

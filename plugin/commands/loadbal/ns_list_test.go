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

var _ = Describe("Load balancer edit policies", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *loadbal.NetscalerListCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeLBManager *testhelpers.FakeLoadBalancerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = loadbal.NewNetscalerListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cliCommand.LoadBalancerManager = fakeLBManager

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
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("ns list", func() {
			It("list all netscalers", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID     Location   Name               Description               IP Address    Management IP   Bandwidth   Create Date"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123    dal01      Netscaler Name     Description Netscaler     10.10.10.10   20.20.20.20     2.000000    2016-12-29T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234   dal02      Netscaler Name 2   Description Netscaler 2   10.10.10.11   20.20.20.21     3.000000    2016-12-29T00:00:00Z"))
			})

			It("list all netscalers in output json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"ID": "123",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Location": "dal01",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Netscaler Name",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Description": "Description Netscaler",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"IP Address": "10.10.10.10",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Create Date": "2016-12-29T00:00:00Z"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))

			})
		})

		Context("errors", func() {
			It("Failed to get netscalers on your account.", func() {
				fakeLBManager.GetADCsReturns([]datatypes.Network_Application_Delivery_Controller{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get netscalers on your account"))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
			It("Failed to get netscalers on your account.", func() {
				fakeLBManager.GetADCsReturns([]datatypes.Network_Application_Delivery_Controller{}, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("No netscalers was found."))
			})
		})
	})
})

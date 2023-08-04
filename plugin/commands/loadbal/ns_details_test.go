package loadbal_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
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
		cliCommand    *loadbal.NetscalerDetailCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeLBManager *testhelpers.FakeLoadBalancerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = loadbal.NewNetscalerDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cliCommand.LoadBalancerManager = fakeLBManager

		expirationDate, _ := time.Parse(time.RFC3339, "2016-12-29T00:00:00Z")
		fakeLBManager.GetADCReturns(datatypes.Network_Application_Delivery_Controller{
			Id:   sl.Int(123),
			Name: sl.String("Netscaler name"),
			Datacenter: &datatypes.Location{
				LongName: sl.String("dal01"),
			},
			PrimaryIpAddress: sl.String("10.10.10.10"),
			Password: &datatypes.Software_Component_Password{
				Password: sl.String("abcde123456"),
			},
			ManagementIpAddress:   sl.String("11.11.11.11"),
			LicenseExpirationDate: sl.Time(expirationDate),
			Subnets: []datatypes.Network_Subnet{
				datatypes.Network_Subnet{
					Id:                sl.Int(456),
					NetworkIdentifier: sl.String("Network identifier"),
					Cidr:              sl.Int(789),
					SubnetType:        sl.String("Type"),
					AddressSpace:      sl.String("Addres subnet"),
				},
				datatypes.Network_Subnet{
					Id:                sl.Int(4567),
					NetworkIdentifier: sl.String("Network identifier 2"),
					Cidr:              sl.Int(7890),
					SubnetType:        sl.String("Type 2"),
					AddressSpace:      sl.String("Addres subnet 2"),
				},
			},
			NetworkVlans: []datatypes.Network_Vlan{
				datatypes.Network_Vlan{
					Id:         sl.Int(987),
					VlanNumber: sl.Int(654),
				},
				datatypes.Network_Vlan{
					Id:         sl.Int(9876),
					VlanNumber: sl.Int(6543),
				},
			},
		}, nil)
	})

	Describe("ns details", func() {
		Context("ns details, Invalid Usage", func() {
			It("ID is required", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage : This command requires one argument."))
			})
			It("ID is required", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: The netscaler ID has to be a positive integer."))
			})
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123465", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("ns detail", func() {
			It("with correct id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name                 Value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID                   123"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name                 Netscaler name"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Location             dal01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Management IP        11.11.11.11"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Root Password        abcde123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Primary IP           10.10.10.10"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("License Expiration   2016-12-29T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Subnet               ID     Subnet                      Type     Space"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("456    Network identifier/789      Type     Addres subnet"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("4567   Network identifier 2/7890   Type 2   Addres subnet 2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Vlans                ID     Number"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("987    654"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("9876   6543"))
			})
			It("with correct id in output json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--output=json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "123"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "Netscaler name"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "11.11.11.11"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "abcde123456"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "2016-12-29T00:00:00Z"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})

		Context("errors", func() {
			It("Failed to get netscaler", func() {
				fakeUI.Inputs("yes")
				fakeLBManager.GetADCReturns(datatypes.Network_Application_Delivery_Controller{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get netscaler 123 on your account"))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})
	})
})

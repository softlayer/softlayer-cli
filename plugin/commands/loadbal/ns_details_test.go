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

var _ = Describe("Load balancer net scaler detail", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeLBManager *testhelpers.FakeLoadBalancerManager
		cmd           *loadbal.NetscalerDetailCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewNetscalerDetailCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        loadbal.LoadbalNetscalerDetailMetadata().Name,
			Description: loadbal.LoadbalNetscalerDetailMetadata().Description,
			Usage:       loadbal.LoadbalNetscalerDetailMetadata().Usage,
			Flags:       loadbal.LoadbalNetscalerDetailMetadata().Flags,
			Action:      cmd.Run,
		}
	})

	Context("Return error", func() {
		It("Set command without id", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Netscaler ID is required."))
		})

		It("Set command with an incorrect id", func() {
			err := testhelpers.RunCommand(cliCommand, "abcde")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: The netscaler ID has to be a positive integer."))
		})

		It("Set command with an invalid output option", func() {
			err := testhelpers.RunCommand(cliCommand, "123465", "--output=xml")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
		})
	})

	Context("Return error", func() {
		BeforeEach(func() {
			fakeLBManager.GetADCReturns(datatypes.Network_Application_Delivery_Controller{}, errors.New("Failed to get netscaler 123465 on your account.SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123465'. (HTTP 404)"))
		})
		It("Set command with an id that doesnot exist", func() {
			err := testhelpers.RunCommand(cliCommand, "123465")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to get netscaler 123465 on your account.SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123465'. (HTTP 404)"))
		})
	})

	Context("Return no error", func() {
		BeforeEach(func() {
			licenseExpiration, _ := time.Parse(time.RFC3339, "2015-10-09T00:00:00Z")
			fakeNetScaler := datatypes.Network_Application_Delivery_Controller{
				Id:          sl.Int(123456),
				Name:        sl.String("SLADC123456-jt22"),
				Description: sl.String("Citrix NetScaler VPX 12.1 10Mbps - Standard"),
				Datacenter: &datatypes.Location{
					LongName: sl.String("Washington 7"),
				},
				PrimaryIpAddress: sl.String("10.0.0.200"),
				Password: &datatypes.Software_Component_Password{
					Password: sl.String("abc12345678"),
				},
				ManagementIpAddress:   sl.String("170.61.13.144"),
				LicenseExpirationDate: sl.Time(licenseExpiration),
				Subnets: []datatypes.Network_Subnet{
					datatypes.Network_Subnet{
						Id:                sl.Int(1234),
						NetworkIdentifier: sl.String("10.0.0.64"),
						Cidr:              sl.Int(24),
						SubnetType:        sl.String("STATIC_IP_ROUTED"),
						AddressSpace:      sl.String("PUBLIC"),
					},
				},
				NetworkVlans: []datatypes.Network_Vlan{
					datatypes.Network_Vlan{
						Id:         sl.Int(123),
						VlanNumber: sl.Int(800),
					},
				},
			}
			fakeLBManager.GetADCReturns(fakeNetScaler, nil)
		})
		It("Return no error", func() {
			err := testhelpers.RunCommand(cliCommand, "123456")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("123456"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("SLADC123456-jt22"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Citrix NetScaler VPX 12.1 10Mbps - Standard"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Washington 7"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("10.0.0.200"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("abc12345678"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("170.61.13.144"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("10.0.0.64/24"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("STATIC_IP_ROUTED"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("PUBLIC"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("123"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("800"))
		})
	})
})

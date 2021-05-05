package hardware_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware detail", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cmd                 *hardware.DetailCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		cmd = hardware.NewDetailCommand(fakeUI, fakeHardwareManager)
		cliCommand = cli.Command{
			Name:        metadata.HardwareDetailMetaData().Name,
			Description: metadata.HardwareDetailMetaData().Description,
			Usage:       metadata.HardwareDetailMetaData().Usage,
			Flags:       metadata.HardwareDetailMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("hardware detail", func() {
		Context("hardware detail without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("hardware detail with wrong hardware ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'. It must be a positive integer."))
			})
		})

		Context("hardware detail with server fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetHardwareReturns(datatypes.Hardware_Server{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get hardware server: 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("VS detail with correct VS ID ", func() {
			created, _ := time.Parse(time.RFC3339, "2017-11-08T00:00:00Z")
			BeforeEach(func() {
				fakeHardwareManager.GetHardwareReturns(
					datatypes.Hardware_Server{
						Hardware: datatypes.Hardware{
							Id:                       sl.Int(1234),
							GlobalIdentifier:         sl.String("rthtoshfkthr"),
							Hostname:                 sl.String("hw-abc"),
							Domain:                   sl.String("wilma.com"),
							FullyQualifiedDomainName: sl.String("hw-abc.wilma.com"),
							HardwareStatus: &datatypes.Hardware_Status{
								Status: sl.String("Active"),
							},
							Datacenter: &datatypes.Location{
								Name: sl.String("dal10"),
							},
							OperatingSystem: &datatypes.Software_Component_OperatingSystem{
								Software_Component: datatypes.Software_Component{
									SoftwareLicense: &datatypes.Software_License{
										SoftwareDescription: &datatypes.Software_Description{
											Name:    sl.String("CentOS"),
											Version: sl.String("6.0"),
										},
									},
									Passwords: []datatypes.Software_Component_Password{
										datatypes.Software_Component_Password{
											Username: sl.String("root"),
											Password: sl.String("password4root"),
										},
									},
								},
							},
							ProcessorPhysicalCoreAmount: sl.Uint(8),
							MemoryCapacity:              sl.Uint(32),
							PrimaryIpAddress:            sl.String("9.9.9.9"),
							PrimaryBackendIpAddress:     sl.String("1.1.1.1"),
							NetworkManagementIpAddress:  sl.String("2.2.2.2"),
							ProvisionDate:               sl.Time(created),
							BillingItem: &datatypes.Billing_Item_Hardware{
								Billing_Item: datatypes.Billing_Item{
									OrderItem: &datatypes.Billing_Order_Item{
										Order: &datatypes.Billing_Order{
											UserRecord: &datatypes.User_Customer{
												Username: sl.String("wilmawang"),
											},
										},
									},
									RecurringFee:                    sl.Float(1000.00),
									NextInvoiceTotalRecurringAmount: sl.Float(1000.00),
								},
							},
							Notes: sl.String("mynotes"),
							TagReferences: []datatypes.Tag_Reference{
								datatypes.Tag_Reference{
									Tag: &datatypes.Tag{
										Name: sl.String("tag1"),
									},
								},
								datatypes.Tag_Reference{
									Tag: &datatypes.Tag{
										Name: sl.String("tag2"),
									},
								},
							},
							NetworkVlans: []datatypes.Network_Vlan{
								datatypes.Network_Vlan{
									Id:           sl.Int(678),
									VlanNumber:   sl.Int(50),
									NetworkSpace: sl.String("PRIMARY"),
								},
							},
						},
					}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("rthtoshfkthr"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("hw-abc"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("wilma.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("hw-abc.wilma.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Active"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal10"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("CentOS"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("6.0"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("8"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("32G"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("9.9.9.9"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1.1.1.1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2.2.2.2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("wilmawang"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("mynotes"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("tag1,tag2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("678"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("50"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("PRIMARY"))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("root"))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("password4root"))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("1000.00"))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--passwords", "--price")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("rthtoshfkthr"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("hw-abc"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("wilma.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("hw-abc.wilma.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Active"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal10"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("CentOS"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("6.0"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("8"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("32G"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("9.9.9.9"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1.1.1.1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2.2.2.2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("wilmawang"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("mynotes"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("tag1,tag2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("678"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("50"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("PRIMARY"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("root"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("password4root"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1000.00"))
			})
		})
	})
})

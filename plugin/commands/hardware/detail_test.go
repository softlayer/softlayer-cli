package hardware_test

import (
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware detail", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cliCommand          *hardware.DetailCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.HardwareManager = fakeHardwareManager
	})

	Describe("hardware detail", func() {
		Context("hardware detail without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("hardware detail with wrong hardware ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'. It must be a positive integer."))
			})
		})

		Context("hardware detail with server fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetHardwareReturns(datatypes.Hardware_Server{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get hardware server: 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Failed to get the hard drives detail", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetHardwareReturns(datatypes.Hardware_Server{}, nil)
				fakeHardwareManager.GetHardDrivesReturns([]datatypes.Hardware_Component{}, errors.New("Failed to get the hard drives detail"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get the hard drives detail"))
			})
		})

		Context("Failed to get bandwidth allotment detail", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetHardwareReturns(datatypes.Hardware_Server{}, nil)
				fakeHardwareManager.GetHardDrivesReturns([]datatypes.Hardware_Component{}, nil)
				fakeHardwareManager.GetBandwidthAllotmentDetailReturns(datatypes.Network_Bandwidth_Version1_Allotment_Detail{}, errors.New("Failed to get bandwidth allotment detail"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get bandwidth allotment detail"))
			})
		})

		Context("Failed to get billing cycle bandwidth usage", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetHardwareReturns(datatypes.Hardware_Server{}, nil)
				fakeHardwareManager.GetHardDrivesReturns([]datatypes.Hardware_Component{}, nil)
				fakeHardwareManager.GetBandwidthAllotmentDetailReturns(datatypes.Network_Bandwidth_Version1_Allotment_Detail{}, nil)
				fakeHardwareManager.GetBillingCycleBandwidthUsageReturns([]datatypes.Network_Bandwidth_Usage{}, errors.New("Failed to get billing cycle bandwidth usage"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get billing cycle bandwidth usage"))
			})
		})

		Context("hardware detail with correct hardware ID ", func() {
			created, _ := time.Parse(time.RFC3339, "2017-11-08T00:00:00Z")
			firmwareCreated, _ := time.Parse(time.RFC3339, "2015-10-09T00:00:00Z")
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
									NextInvoiceChildren: []datatypes.Billing_Item{
										datatypes.Billing_Item{
											Description:                     sl.String("CentOS 7.x (64 bit)"),
											CategoryCode:                    sl.String("os"),
											NextInvoiceTotalRecurringAmount: sl.Float(0.00),
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
							LastTransaction: &datatypes.Provisioning_Version1_Transaction{
								TransactionGroup: &datatypes.Provisioning_Version1_Transaction_Group{
									Name: sl.String("Storage_EvaultProvision"),
								},
								ModifyDate: sl.Time(created),
							},
							HourlyBillingFlag: sl.Bool(true),
							ActiveComponents: []datatypes.Hardware_Component{
								datatypes.Hardware_Component{
									HardwareComponentModel: &datatypes.Hardware_Component_Model{
										HardwareGenericComponentModel: &datatypes.Hardware_Component_Model_Generic{
											HardwareComponentType: &datatypes.Hardware_Component_Type{
												KeyName: sl.String("DRIVE_CONTROLLER"),
											},
										},
										LongDescription: sl.String("LSI / DRIVE CONTROLLER / Avago MegaRAID 9361-8i / SATA/SAS - MegaRAID SAS 9361-8i / 8"),
									},
								},
							},
						},
					}, nil)
				fakeHardwareManager.GetHardwareComponentsReturns([]datatypes.Hardware_Component{
					datatypes.Hardware_Component{
						Id: sl.Int(18787137),
						HardwareComponentModel: &datatypes.Hardware_Component_Model{
							LongDescription: sl.String("Aspeed / AST2500 - Onboard / IPMI - KVM / Remote Management Count1"),
							Firmwares: []datatypes.Hardware_Component_Firmware{
								datatypes.Hardware_Component_Firmware{
									Version:    sl.String("3.10"),
									CreateDate: sl.Time(firmwareCreated),
								},
							},
							HardwareGenericComponentModel: &datatypes.Hardware_Component_Model_Generic{
								HardwareComponentType: &datatypes.Hardware_Component_Type{
									KeyName: sl.String("REMOTE_MGMT_CARD"),
								},
							},
						},
					},
				}, nil)
				fakeHardwareManager.GetHardDrivesReturns(
					[]datatypes.Hardware_Component{
						datatypes.Hardware_Component{
							HardwareComponentModel: &datatypes.Hardware_Component_Model{
								Manufacturer: sl.String("Seagate"),
								Name:         sl.String("Constellation ES"),
								HardwareGenericComponentModel: &datatypes.Hardware_Component_Model_Generic{
									Capacity: sl.Float(1000.00),
									Units:    sl.String("GB"),
								},
							},
							SerialNumber: sl.String("z1w4zqye"),
						},
					},
					nil,
				)
				fakeHardwareManager.GetBandwidthAllotmentDetailReturns(
					datatypes.Network_Bandwidth_Version1_Allotment_Detail{
						Allocation: &datatypes.Network_Bandwidth_Version1_Allocation{
							Amount: sl.Float(20000),
						},
					},
					nil,
				)
				fakeHardwareManager.GetBillingCycleBandwidthUsageReturns(
					[]datatypes.Network_Bandwidth_Usage{
						datatypes.Network_Bandwidth_Usage{
							Type: &datatypes.Network_Bandwidth_Version1_Usage_Detail_Type{
								Alias: sl.String("PUBLIC_SERVER_BW"),
							},
							AmountIn:  sl.Float(0.326090),
							AmountOut: sl.Float(0.015740),
						},
					},
					nil,
				)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
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
				Expect(fakeUI.Outputs()).To(ContainSubstring("Seagate Constellation ES"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1000.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("z1w4zqye"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Storage_EvaultProvision 2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hourly"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Public"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.326090"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.015740"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("20000"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("DRIVE_CONTROLLER"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("LSI / DRIVE CONTROLLER / Avago MegaRAID 9361-8i / SATA/SAS - MegaRAID SAS 9361-8i / 8"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--passwords", "--price", "--components")
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
				Expect(fakeUI.Outputs()).To(ContainSubstring("Seagate Constellation ES"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1000.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("z1w4zqye"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Storage_EvaultProvision 2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hourly"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Public"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.326090"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.015740"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("20000"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("DRIVE_CONTROLLER"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("LSI / DRIVE CONTROLLER / Avago MegaRAID 9361-8i / SATA/SAS - MegaRAID SAS 9361-8i / 8"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("root"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("password4root"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1000.00"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("CentOS 7.x (64 bit)"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("os"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.00"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aspeed / AST2500 - Onboard / IPMI - KVM / Remote Management Count1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("3.10"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2015-10-09T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("REMOTE_MGMT_CARD"))
			})
		})

		Context("Issue #649", func() {
			BeforeEach(func() {
				fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
				slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
				cliCommand = hardware.NewDetailCommand(slCommand)
			})
			It("return hardware detail", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1403539")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Last transaction   - -"))
			})
		})
	})
})

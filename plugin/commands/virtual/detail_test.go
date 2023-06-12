package virtual_test

import (
	"errors"
	"strings"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var created, _ = time.Parse(time.RFC3339, "2016-12-25T00:00:00Z")
var modified, _ = time.Parse(time.RFC3339, "2017-01-01T00:00:00Z")
var lastTransaction, _ = time.Parse(time.RFC3339, "2017-02-01T00:00:00Z")

var GetInstanceReturn = datatypes.Virtual_Guest{

	Id:                       sl.Int(1234),
	GlobalIdentifier:         sl.String("rthtoshfkthr"),
	Hostname:                 sl.String("vs-abc"),
	Domain:                   sl.String("wilma.com"),
	FullyQualifiedDomainName: sl.String("vs-abc.wilma.com"),
	Status:                   &datatypes.Virtual_Guest_Status{Name: sl.String("Provisioning")},
	PowerState:               &datatypes.Virtual_Guest_Power_State{Name: sl.String("PowerOn")},
	ActiveTransaction: &datatypes.Provisioning_Version1_Transaction{
		TransactionStatus: &datatypes.Provisioning_Version1_Transaction_Status{Name: sl.String("Provisioning")},
	},
	Datacenter: &datatypes.Location{Name: sl.String("dal10")},
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
	MaxCpu:                       sl.Int(8),
	MaxMemory:                    sl.Int(4096),
	PrimaryIpAddress:             sl.String("9.9.9.9"),
	PrimaryBackendIpAddress:      sl.String("1.1.1.1"),
	PrivateNetworkOnlyFlag:       sl.Bool(false),
	DedicatedAccountHostOnlyFlag: sl.Bool(false),
	CreateDate:                   sl.Time(created),
	ModifyDate:                   sl.Time(modified),
	BillingItem: &datatypes.Billing_Item_Virtual_Guest{
		Billing_Item: datatypes.Billing_Item{
			OrderItem: &datatypes.Billing_Order_Item{
				Order: &datatypes.Billing_Order{
					UserRecord: &datatypes.User_Customer{Username: sl.String("wilmawang")},
				},
				Preset: &datatypes.Product_Package_Preset{KeyName: sl.String("C1_2X2X25")},
			},
			RecurringFee:                    sl.Float(1000.00),
			NextInvoiceTotalRecurringAmount: sl.Float(1000.00),
			NextInvoiceChildren: []datatypes.Billing_Item{
				datatypes.Billing_Item{
					RecurringFee: sl.Float(1000.00),
					Description:  sl.String("CPU Cores: a suspendable product. Anticipated usage for the billing cycle is 743.9997 hours Used"),
					CategoryCode: sl.String("guest_core_usage"),
				},
			},
		},
	},
	Notes: sl.String("mynotes"),
	TagReferences: []datatypes.Tag_Reference{
		datatypes.Tag_Reference{
			Tag: &datatypes.Tag{Name: sl.String("tag1")},
		},
		datatypes.Tag_Reference{
			Tag: &datatypes.Tag{Name: sl.String("tag2")},
		},
	},
	NetworkVlans: []datatypes.Network_Vlan{
		datatypes.Network_Vlan{
			Id:           sl.Int(678),
			VlanNumber:   sl.Int(50),
			NetworkSpace: sl.String("PRIMARY"),
		},
	},
	TransientGuestFlag: sl.Bool(false),
	LastTransaction: &datatypes.Provisioning_Version1_Transaction{
		TransactionGroup: &datatypes.Provisioning_Version1_Transaction_Group{Name: sl.String("Service Setup")},
		ModifyDate:       sl.Time(lastTransaction),
	},
	HourlyBillingFlag: sl.Bool(true),
}

var BlockDeviceReturns = []datatypes.Virtual_Guest_Block_Device{
	datatypes.Virtual_Guest_Block_Device{
		DiskImage: &datatypes.Virtual_Disk_Image{
			Description: sl.String("123456789-SWAP"),
			Capacity:    sl.Int(2),
			Units:       sl.String("GB"),
		},
		MountType: sl.String("Disk"),
		Device:    sl.String("1"),
		Id:        sl.Int(111111),
	},
}

var _ = Describe("VS detail", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *virtual.DetailCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeVSManager *testhelpers.FakeVirtualServerManager
		fakeTransport *testhelpers.FakeTransportHandler
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeTransport = new(testhelpers.FakeTransportHandler)
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
	})
	Describe("VS detail", func() {
		Context("VS detail without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("VS detail with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Virtual server ID'. It must be a positive integer."))
			})
		})

		Context("VS detail with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{}, nil)
				fakeVSManager.GetLocalDisksReturns([]datatypes.Virtual_Guest_Block_Device{}, errors.New("Failed to get the local disks detail for the virtual server"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get the local disks detail for the virtual server"))
			})
		})

		Context("VS detail with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get virtual server instance: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("VS detail with correct VS ID ", func() {

			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(GetInstanceReturn, nil)
				fakeVSManager.GetLocalDisksReturns(BlockDeviceReturns, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("rthtoshfkthr"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("vs-abc.wilma.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Provisioning"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("PowerOn"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal10"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1.1.1.1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("false"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-12-25T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-01-01T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Service Setup (2017-02-01T00:00:00Z)"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hourly"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("C1_2X2X25"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("PRIMARY"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("root"))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("password4root"))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("1000.00"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--passwords", "--price")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("rthtoshfkthr"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("vs-abc"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("wilma.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("vs-abc.wilma.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Provisioning"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("PRIMARY"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("root"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("password4root"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1000.00"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("CPU Cores: a suspendable product. Anticipated usage for the billing cycle is 743.9997 hours Used"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("guest_core_usage"))
			})
		})
		Context("Github issues #252", func() {

			BeforeEach(func() {
				GetInstanceReturn.BillingItem = nil
				fakeVSManager.GetInstanceReturns(GetInstanceReturn, nil)
				fakeVSManager.GetLocalDisksReturns(BlockDeviceReturns, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--passwords", "--price")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("Price rate"))
			})
		})
		Context("Github issues #540", func() {
			var guestInstance datatypes.Virtual_Guest
			var guestBlockDevices []datatypes.Virtual_Guest_Block_Device
			options := sl.Options{Id: sl.Int(124929698)}
			BeforeEach(func() {
				errAPI := fakeTransport.DoRequest(fakeSession, "SoftLayer_Virtual_Guest", "getObject", nil, &options, &guestInstance)
				Expect(errAPI).NotTo(HaveOccurred())
				errAPI = fakeTransport.DoRequest(fakeSession, "SoftLayer_Virtual_Guest", "getBlockDevices", nil, &options, &guestBlockDevices)
				Expect(errAPI).NotTo(HaveOccurred())
				fakeVSManager.GetInstanceReturns(guestInstance, nil)
				fakeVSManager.GetLocalDisksReturns(guestBlockDevices, nil)
			})
			It("handles virtual servers with CDs mounted", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				output := fakeUI.Outputs()
				Expect(output).NotTo(ContainSubstring("Price rate"))
				Expect(output).To(ContainSubstring("Rescue   CD     3"))
			})
		})
	})
})

package hardware_test

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
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware list", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cliCommand          *hardware.ListCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.HardwareManager = fakeHardwareManager
	})

	Describe("hardware list", func() {
		Context("hardware list with wrong parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--column", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --column abc is not supported."))
			})
		})
		Context("hardware list with wrong parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --sortby abc is not supported."))
			})
		})
		Context("hardware list with server fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.ListHardwareReturns([]datatypes.Hardware_Server{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get hardware servers on your account."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})
		Context("hardware list with different --sortby", func() {
			BeforeEach(func() {
				created1, _ := time.Parse(time.RFC3339, "2017-11-08T00:00:00Z")
				created2, _ := time.Parse(time.RFC3339, "2017-11-09T00:00:00Z")
				fakeHardwareManager.ListHardwareReturns([]datatypes.Hardware_Server{
					datatypes.Hardware_Server{
						Hardware: datatypes.Hardware{
							Id:               sl.Int(1234),
							GlobalIdentifier: sl.String("rthtoshfkthr"),
							Hostname:         sl.String("hw-abc"),
							Domain:           sl.String("wilma.com"),
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
								},
							},
							ProcessorPhysicalCoreAmount: sl.Uint(8),
							MemoryCapacity:              sl.Uint(32),
							PrimaryIpAddress:            sl.String("9.9.9.9"),
							PrimaryBackendIpAddress:     sl.String("1.1.1.1"),
							NetworkManagementIpAddress:  sl.String("2.2.2.2"),
							ProvisionDate:               sl.Time(created1),
							BillingItem: &datatypes.Billing_Item_Hardware{
								Billing_Item: datatypes.Billing_Item{
									OrderItem: &datatypes.Billing_Order_Item{
										Order: &datatypes.Billing_Order{
											UserRecord: &datatypes.User_Customer{
												Username: sl.String("wilmawang"),
											},
										},
									},
								},
							},
						},
					},
					datatypes.Hardware_Server{
						Hardware: datatypes.Hardware{
							Id:               sl.Int(4321),
							GlobalIdentifier: sl.String("toshfkthdddr"),
							Hostname:         sl.String("hw-cbs"),
							Domain:           sl.String("ibm.com"),
							HardwareStatus: &datatypes.Hardware_Status{
								Status: sl.String("Shutdown"),
							},
							Datacenter: &datatypes.Location{
								Name: sl.String("tok02"),
							},
							OperatingSystem: &datatypes.Software_Component_OperatingSystem{
								Software_Component: datatypes.Software_Component{
									SoftwareLicense: &datatypes.Software_License{
										SoftwareDescription: &datatypes.Software_Description{
											Name:    sl.String("Alpine"),
											Version: sl.String("1.0"),
										},
									},
								},
							},
							ProcessorPhysicalCoreAmount: sl.Uint(4),
							MemoryCapacity:              sl.Uint(16),
							PrimaryIpAddress:            sl.String("8.9.9.9"),
							PrimaryBackendIpAddress:     sl.String("2.1.1.1"),
							NetworkManagementIpAddress:  sl.String("5.2.2.2"),
							ProvisionDate:               sl.Time(created2),
							BillingItem: &datatypes.Billing_Item_Hardware{
								Billing_Item: datatypes.Billing_Item{
									OrderItem: &datatypes.Billing_Order_Item{
										Order: &datatypes.Billing_Order{
											UserRecord: &datatypes.User_Customer{
												Username: sl.String("AlexWang"),
											},
										},
									},
								},
							},
						},
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("1234"))
				Expect(results[2]).To(ContainSubstring("4321"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "hostname")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("hw-abc"))
				Expect(results[2]).To(ContainSubstring("hw-cbs"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "domain")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("ibm.com"))
				Expect(results[2]).To(ContainSubstring("wilma.com"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "datacenter")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("dal10"))
				Expect(results[2]).To(ContainSubstring("tok02"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "cpu", "--column", "cpu")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("4"))
				Expect(results[2]).To(ContainSubstring("8"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "memory", "--column", "memory")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("16"))
				Expect(results[2]).To(ContainSubstring("32"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "public_ip")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("8.9.9.9"))
				Expect(results[2]).To(ContainSubstring("9.9.9.9"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "private_ip")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("1.1.1.1"))
				Expect(results[2]).To(ContainSubstring("2.1.1.1"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "ipmi_ip", "--column", "ipmi_ip")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("2.2.2.2"))
				Expect(results[2]).To(ContainSubstring("5.2.2.2"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "created", "--column", "created")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("2017-11-08T00:00:00Z"))
				Expect(results[2]).To(ContainSubstring("2017-11-09T00:00:00Z"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "created_by", "--column", "created_by")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("AlexWang"))
				Expect(results[2]).To(ContainSubstring("wilmawang"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "os", "--column", "os")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("Alpine"))
				Expect(results[2]).To(ContainSubstring("CentOS"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "status", "--column", "status")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("Active"))
				Expect(results[2]).To(ContainSubstring("Shutdown"))
			})
		})
	})
})

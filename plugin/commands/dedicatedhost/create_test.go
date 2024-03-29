package dedicatedhost_test

import (
	"errors"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dedicatedhost"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Dedicated host create", func() {
	var (
		fakeUI                   *terminal.FakeUI
		cliCommand               *dedicatedhost.CreateCommand
		fakeSession              *session.Session
		slCommand                *metadata.SoftlayerCommand
		FakeDedicatedhostManager *testhelpers.FakeDedicatedHostManager
		FakeNetworkManager       *testhelpers.FakeNetworkManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = dedicatedhost.NewCreateCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		FakeDedicatedhostManager = new(testhelpers.FakeDedicatedHostManager)
		cliCommand.DedicatedHostManager = FakeDedicatedhostManager
		FakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cliCommand.NetworkManager = FakeNetworkManager
	})

	Describe("Dedicatedhost create", func() {
		Context("Create host without hostname", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--hostname' is required"))
			})
		})
		Context("Create host without domain", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "test")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--domain' is required"))
			})
		})
		Context("Create host without datacenter", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "test", "--domain", "softlayer.com")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--datacenter' is required"))
			})
		})
		Context("Create host without vlan", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "test", "--domain", "softlayer.com", "--datacenter", "dal09")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--vlan-private' is required"))
			})
		})
		Context("Create host with wrong billing", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "test", "--domain", "softlayer.com", "--datacenter", "dal09", "--billing", "dbd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [--billing] has to be either hourly or monthly."))
			})
		})
		Context("Create host with get vlan fails", func() {
			BeforeEach(func() {
				FakeNetworkManager.GetVlanReturns(datatypes.Network_Vlan{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "test", "--domain", "softlayer.com", "--datacenter", "dal09", "--billing", "hourly", "--vlan-private", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get vlan 123."))
			})
		})

		Context("Create host without -f and not continue", func() {
			BeforeEach(func() {
				FakeNetworkManager.GetVlanReturns(datatypes.Network_Vlan{}, nil)
			})
			It("return error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "test", "--domain", "softlayer.com", "--datacenter", "dal09", "--billing", "hourly", "--vlan-private", "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This action will incur charges on your account. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
			})
		})

		Context("Verify the vlan order with fail because the vlan is in a different location", func() {
			BeforeEach(func() {
				FakeNetworkManager.GetVlanReturns(datatypes.Network_Vlan{
					Id: sl.Int(123),
					PrimaryRouter: &datatypes.Hardware_Router{
						Hardware_Switch: datatypes.Hardware_Switch{
							Hardware: datatypes.Hardware{
								Id:       sl.Int(1115295),
								Hostname: sl.String("bcr01a.wdc07"),
								Datacenter: &datatypes.Location{
									Name: sl.String("wdc07"),
								},
							},
						},
					},
				}, nil)
				FakeDedicatedhostManager.GenerateOrderTemplateReturns(datatypes.Container_Product_Order_Virtual_DedicatedHost{
					Container_Product_Order: datatypes.Container_Product_Order{
						Hardware: []datatypes.Hardware{
							datatypes.Hardware{
								Domain:   sl.String("test.com"),
								Hostname: sl.String("test"),
								PrimaryBackendNetworkComponent: &datatypes.Network_Component{
									Router: &datatypes.Hardware{
										Id: sl.Int(1234567),
									},
								},
							},
						},
						Location:  sl.String("AMSTERDAM"),
						PackageId: sl.Int(813),
						Prices: []datatypes.Product_Item_Price{
							datatypes.Product_Item_Price{
								Id: sl.Int(200269),
							},
						},
						UseHourlyPricing: sl.Bool(true),
					},
				}, nil)
				FakeDedicatedhostManager.VerifyInstanceCreationReturns(datatypes.Container_Product_Order{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "test", "--domain", "softlayer.com", "--datacenter", "dal09", "--billing", "hourly", "--vlan-private", "123", "--test")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("The vlan is located at: wdc07, Please add a valid private vlan according the datacenter selected."))
			})
		})

		Context("Generate host with verify host fails", func() {
			BeforeEach(func() {
				FakeNetworkManager.GetVlanReturns(datatypes.Network_Vlan{
					Id: sl.Int(123),
					PrimaryRouter: &datatypes.Hardware_Router{
						Hardware_Switch: datatypes.Hardware_Switch{
							Hardware: datatypes.Hardware{
								Id:       sl.Int(1115295),
								Hostname: sl.String("bcr01a.dal09"),
								Datacenter: &datatypes.Location{
									Name: sl.String("dal09"),
								},
							},
						},
					},
				}, nil)
				FakeDedicatedhostManager.GenerateOrderTemplateReturns(datatypes.Container_Product_Order_Virtual_DedicatedHost{
					Container_Product_Order: datatypes.Container_Product_Order{
						Hardware: []datatypes.Hardware{
							datatypes.Hardware{
								Domain:   sl.String("test.com"),
								Hostname: sl.String("test"),
								PrimaryBackendNetworkComponent: &datatypes.Network_Component{
									Router: &datatypes.Hardware{
										Id: sl.Int(1234567),
									},
								},
							},
						},
						Location:  sl.String("AMSTERDAM"),
						PackageId: sl.Int(813),
						Prices: []datatypes.Product_Item_Price{
							datatypes.Product_Item_Price{
								Id: sl.Int(200269),
							},
						},
						UseHourlyPricing: sl.Bool(true),
					},
				}, nil)
				FakeDedicatedhostManager.VerifyInstanceCreationReturns(datatypes.Container_Product_Order{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "test", "--domain", "softlayer.com", "--datacenter", "dal09", "--billing", "hourly", "--vlan-private", "123", "--test")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to verify virtual server creation.\n"))
			})
		})

		Context("Generate host with order host fails", func() {
			BeforeEach(func() {
				FakeNetworkManager.GetVlanReturns(datatypes.Network_Vlan{
					Id: sl.Int(123),
					PrimaryRouter: &datatypes.Hardware_Router{
						Hardware_Switch: datatypes.Hardware_Switch{
							Hardware: datatypes.Hardware{
								Id:       sl.Int(1115295),
								Hostname: sl.String("bcr01a.dal09"),
								Datacenter: &datatypes.Location{
									Name: sl.String("dal09"),
								},
							},
						},
					},
				}, nil)
				FakeDedicatedhostManager.GenerateOrderTemplateReturns(datatypes.Container_Product_Order_Virtual_DedicatedHost{
					Container_Product_Order: datatypes.Container_Product_Order{
						Hardware: []datatypes.Hardware{
							datatypes.Hardware{
								Domain:   sl.String("test.com"),
								Hostname: sl.String("test"),
								PrimaryBackendNetworkComponent: &datatypes.Network_Component{
									Router: &datatypes.Hardware{
										Id: sl.Int(1234567),
									},
								},
							},
						},
						Location:  sl.String("AMSTERDAM"),
						PackageId: sl.Int(813),
						Prices: []datatypes.Product_Item_Price{
							datatypes.Product_Item_Price{
								Id: sl.Int(200269),
							},
						},
						UseHourlyPricing: sl.Bool(true),
					},
				}, nil)
				FakeDedicatedhostManager.OrderInstanceReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "test", "--domain", "softlayer.com", "--datacenter", "dal09", "--billing", "hourly", "--vlan-private", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to Order the dedicatedhost.\n"))
			})
		})

		Context("Verify host create with succeed", func() {
			BeforeEach(func() {
				FakeNetworkManager.GetVlanReturns(datatypes.Network_Vlan{
					Id: sl.Int(123),
					PrimaryRouter: &datatypes.Hardware_Router{
						Hardware_Switch: datatypes.Hardware_Switch{
							Hardware: datatypes.Hardware{
								Id:       sl.Int(1115295),
								Hostname: sl.String("bcr01a.dal09"),
								Datacenter: &datatypes.Location{
									Name: sl.String("dal09"),
								}},
						},
					},
				}, nil)
				FakeDedicatedhostManager.GenerateOrderTemplateReturns(datatypes.Container_Product_Order_Virtual_DedicatedHost{
					Container_Product_Order: datatypes.Container_Product_Order{
						Hardware: []datatypes.Hardware{
							datatypes.Hardware{
								Domain:   sl.String("test.com"),
								Hostname: sl.String("test"),
								PrimaryBackendNetworkComponent: &datatypes.Network_Component{
									Router: &datatypes.Hardware{
										Id: sl.Int(1234567),
									},
								},
							},
						},
						Location:  sl.String("AMSTERDAM"),
						PackageId: sl.Int(813),
						Prices: []datatypes.Product_Item_Price{
							datatypes.Product_Item_Price{
								Id: sl.Int(200269),
							},
						},
						UseHourlyPricing: sl.Bool(true),
					},
				}, nil)
				FakeDedicatedhostManager.VerifyInstanceCreationReturns(datatypes.Container_Product_Order{
					Hardware: []datatypes.Hardware{
						datatypes.Hardware{
							Domain:   sl.String("test.com"),
							Hostname: sl.String("test"),
							PrimaryBackendNetworkComponent: &datatypes.Network_Component{
								Router: &datatypes.Hardware{
									Id: sl.Int(1234567),
								},
							},
						},
					},
					Location:  sl.String("AMSTERDAM"),
					PackageId: sl.Int(813),
					Prices: []datatypes.Product_Item_Price{
						datatypes.Product_Item_Price{
							Id: sl.Int(200269),
						},
					},
					UseHourlyPricing: sl.Bool(true),
				}, nil)
			})
			It("return order", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "test", "--domain", "softlayer.com", "--datacenter", "dal09", "--billing", "hourly", "--vlan-private", "123", "--size", "56_CORES_X_242_RAM_X_1_4_TB", "--test")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstring("The order is correct."))
			})
		})

		Context("Order host with succeed", func() {
			BeforeEach(func() {
				FakeNetworkManager.GetVlanReturns(datatypes.Network_Vlan{
					Id: sl.Int(123),
					PrimaryRouter: &datatypes.Hardware_Router{
						Hardware_Switch: datatypes.Hardware_Switch{
							Hardware: datatypes.Hardware{
								Id:       sl.Int(1115295),
								Hostname: sl.String("bcr01a.dal09"),
								Datacenter: &datatypes.Location{
									Name: sl.String("dal09"),
								}},
						},
					},
				}, nil)
				FakeDedicatedhostManager.GenerateOrderTemplateReturns(datatypes.Container_Product_Order_Virtual_DedicatedHost{
					Container_Product_Order: datatypes.Container_Product_Order{
						Hardware: []datatypes.Hardware{
							datatypes.Hardware{
								Domain:   sl.String("test.com"),
								Hostname: sl.String("test"),
								PrimaryBackendNetworkComponent: &datatypes.Network_Component{
									Router: &datatypes.Hardware{
										Id: sl.Int(1234567),
									},
								},
							},
						},
						Location:  sl.String("AMSTERDAM"),
						PackageId: sl.Int(813),
						Prices: []datatypes.Product_Item_Price{
							datatypes.Product_Item_Price{
								Id: sl.Int(200269),
							},
						},
						UseHourlyPricing: sl.Bool(true),
					},
				}, nil)
				FakeDedicatedhostManager.OrderInstanceReturns(datatypes.Container_Product_Order_Receipt{
					OrderId: sl.Int(345678),
				}, nil)
			})
			It("return order", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "test", "--domain", "softlayer.com", "--datacenter", "dal09", "--billing", "hourly", "--vlan-private", "123", "--size", "56_CORES_X_242_RAM_X_1_4_TB")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstring("The order 345678 was placed."))
			})
			It("return order", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "test", "--domain", "softlayer.com", "--datacenter", "dal09", "--billing", "hourly", "--vlan-private", "123", "--force")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The order 345678 was placed."))
			})
		})
	})
})

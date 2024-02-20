package loadbal_test

import (
	"errors"
	"strings"

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

var _ = Describe("Load balancer create options", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *loadbal.OptionsCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeLBManager      *testhelpers.FakeLoadBalancerManager
		fakeNetworkManager *testhelpers.FakeNetworkManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = loadbal.NewOptionsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cliCommand.LoadBalancerManager = fakeLBManager
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cliCommand.NetworkManager = fakeNetworkManager

		fakeLBManager.CreateLoadBalancerOptionsReturns([]datatypes.Product_Package{
			datatypes.Product_Package{
				Regions: []datatypes.Location_Region{
					datatypes.Location_Region{
						Keyname: sl.String("REGION_KEY_NAME_1"),
						Location: &datatypes.Location_Region_Location{
							Location: &datatypes.Location{
								Name: sl.String("loc01"),
							},
						},
					},
					datatypes.Location_Region{
						Keyname: sl.String("REGION_KEY_NAME_2"),
						Location: &datatypes.Location_Region_Location{
							Location: &datatypes.Location{
								Name: sl.String("loc02"),
								Groups: []datatypes.Location_Group{
									datatypes.Location_Group{
										Id: sl.Int(123456),
									},
								},
							},
						},
					},
				},
				Items: []datatypes.Product_Item{
					datatypes.Product_Item{
						KeyName: sl.String("KEY_NAME_PRICE_1"),
						Prices: []datatypes.Product_Item_Price{
							datatypes.Product_Item_Price{
								LocationGroupId:    sl.Int(123456),
								HourlyRecurringFee: sl.Float(1.2),
							},
						},
					},
					datatypes.Product_Item{
						KeyName: sl.String("KEY_NAME_PRICE_2"),
						Prices: []datatypes.Product_Item_Price{
							datatypes.Product_Item_Price{
								LocationGroupId:    sl.Int(1234567),
								HourlyRecurringFee: sl.Float(1.2),
							},
							datatypes.Product_Item_Price{
								HourlyRecurringFee: sl.Float(2.2),
							},
						},
					},
				},
			},
		}, nil)

		fakeNetworkManager.ListSubnetsReturns([]datatypes.Network_Subnet{
			datatypes.Network_Subnet{
				SubnetType:        sl.String("PRIMARY"),
				NetworkIdentifier: sl.String("10.10.10.10"),
				Cidr:              sl.Int(12),
				NetworkVlan: &datatypes.Network_Vlan{
					VlanNumber: sl.Int(456),
				},
				PodName: sl.String("test.pod01"),
				Id:      sl.Int(789),
			},
			datatypes.Network_Subnet{
				SubnetType:        sl.String("ADDITIONAL_PRIMARY"),
				NetworkIdentifier: sl.String("20.20.20.20"),
				Cidr:              sl.Int(13),
				NetworkVlan: &datatypes.Network_Vlan{
					VlanNumber: sl.Int(457),
				},
				PodName: sl.String("test.pod02"),
				Id:      sl.Int(781),
			},
			datatypes.Network_Subnet{
				SubnetType:        sl.String("NOT_EXIST"),
				NetworkIdentifier: sl.String("30.30.30.30"),
				Cidr:              sl.Int(14),
				NetworkVlan: &datatypes.Network_Vlan{
					VlanNumber: sl.Int(458),
				},
				PodName: sl.String("test.pod03"),
				Id:      sl.Int(782),
			},
		}, nil)
	})

	Context("create options returns error", func() {
		BeforeEach(func() {
			fakeLBManager.CreateLoadBalancerOptionsReturns(nil, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to get load balancer product packages.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("create options", func() {
		It("return no error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("Datacenter          keyName"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("REGION_KEY_NAME_1   loc01"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("REGION_KEY_NAME_2   loc02"))
		})
	})
	Context("create options with flag -d", func() {
		It("return no error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-d", "loc02")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("Private Subnets"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Key Name"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("test.pod01.456"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("test.pod02.457"))
		})
	})
	Context("create options with 0 subnets", func() {
		BeforeEach(func() {
			fakeNetworkManager.ListSubnetsReturns([]datatypes.Network_Subnet{}, nil)
		})
		It("return no error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-d", "loc02")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("Private Subnets"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Not Found"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("KEY_NAME_PRICE_1"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("2.200000"))
		})
	})
	Context("create options return subnet error", func() {
		BeforeEach(func() {
			fakeNetworkManager.ListSubnetsReturns([]datatypes.Network_Subnet{}, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-d", "loc02")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("Prices:           Private Subnets"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Private Subnets   Failed to get subnets.Internal server error"))
		})
	})
	Context("create options with double Region", func() {
		BeforeEach(func() {
			fakeLBManager.CreateLoadBalancerOptionsReturns([]datatypes.Product_Package{
				datatypes.Product_Package{
					Regions: []datatypes.Location_Region{
						datatypes.Location_Region{
							Keyname: sl.String("REGION_KEY_NAME_1"),
							Location: &datatypes.Location_Region_Location{
								Location: &datatypes.Location{
									Name: sl.String("loc02"),
								},
							},
						},
						datatypes.Location_Region{
							Keyname: sl.String("REGION_KEY_NAME_2"),
							Location: &datatypes.Location_Region_Location{
								Location: &datatypes.Location{
									Name: sl.String("loc02"),
									Groups: []datatypes.Location_Group{
										datatypes.Location_Group{
											Id: sl.Int(123456),
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
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-d", "loc02")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("-----------------------------"))
		})
	})
})

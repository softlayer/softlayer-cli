package block_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Volume options", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		FakeNetworkManager *testhelpers.FakeNetworkManager
		cliCommand         *block.VolumeOptionsCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		FakeNetworkManager = new(testhelpers.FakeNetworkManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewVolumeOptionsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
		cliCommand.NetworkManager = FakeNetworkManager
	})

	Describe("Volume options", func() {
		Context("Volume options with server API call fails", func() {
			It("return error GetRegions", func() {
				FakeStorageManager.GetRegionsReturns([]datatypes.Location_Region{}, errors.New("Internal Server Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get all datacenters by PackageId."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
			It("return error ListItems", func() {
				FakeStorageManager.ListItemsReturns([]datatypes.Product_Item{}, errors.New("Internal Server Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get all storages packages."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
			It("return error GetOsType", func() {
				FakeStorageManager.GetOsTypeReturns([]datatypes.Network_Storage_Iscsi_OS_Type{}, errors.New("Internal Server Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get all os types."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
			It("return error GetAllDatacenters", func() {
				FakeStorageManager.GetAllDatacentersReturns([]datatypes.Location{}, errors.New("Internal Server Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get all datacenters."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
			It("return error GetClosingPods", func() {
				FakeNetworkManager.GetClosingPodsReturns([]datatypes.Network_Pod{}, errors.New("Internal Server Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get Pods."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Volume options", func() {
			BeforeEach(func() {
				FakeStorageManager.GetRegionsReturns([]datatypes.Location_Region{
					{
						Description: sl.String("AMS03 - Amsterdam"),
						Location: &datatypes.Location_Region_Location{
							LocationId: sl.Int(321),
							Location: &datatypes.Location{
								Name:     sl.String("ams01"),
								LongName: sl.String("AMS03 - Amsterdam"),
							}},
					},
					{
						Description: sl.String("MON01 - Montreal"),
						Location: &datatypes.Location_Region_Location{
							LocationId: sl.Int(322),
							Location: &datatypes.Location{
								Name:     sl.String("mon01"),
								LongName: sl.String("MON01 - Montreal"),
							}},
					},
				}, nil)

				FakeStorageManager.ListItemsReturns([]datatypes.Product_Item{
					{
						Prices: []datatypes.Product_Item_Price{
							{CapacityRestrictionType: sl.String("IOPS"),
								CapacityRestrictionMinimum: sl.String("100"),
								HourlyRecurringFee:         sl.Float(0.002),
								RecurringFee:               sl.Float(0.3),
							},
						},
						Id:                sl.Int(123),
						Description:       sl.String("20 - 39 GBs"),
						KeyName:           sl.String("20_39_GBS"),
						LocationConflicts: []datatypes.Product_Item_Resource_Conflict{},
					},
					{
						Prices: []datatypes.Product_Item_Price{
							{CapacityRestrictionType: sl.String("STORAGE_TIER_LEVEL"),
								PricingLocationGroup: &datatypes.Location_Group_Pricing{
									Location_Group: datatypes.Location_Group{Name: sl.String("ams01")},
								},
								CapacityRestrictionMinimum: sl.String("100"),
								HourlyRecurringFee:         sl.Float(0.002),
								RecurringFee:               sl.Float(0.3),
							},
						},
						Id:          sl.Int(124),
						Description: sl.String("20 GB Storage Space"),
						KeyName:     sl.String("20_GB_PERFORMANCE_STORAGE_SPACE"),
						LocationConflicts: []datatypes.Product_Item_Resource_Conflict{
							{
								ResourceTableId: sl.Int(321),
							},
						},
					},
				}, nil)

				FakeStorageManager.GetAllDatacentersReturns([]datatypes.Location{
					{Name: sl.String("ams01"), Id: sl.Int(321)},
					{Name: sl.String("dal06"), Id: sl.Int(324)},
				}, nil)

				FakeStorageManager.GetOsTypeReturns([]datatypes.Network_Storage_Iscsi_OS_Type{
					{
						Name:        sl.String("Linux"),
						KeyName:     sl.String("LINUX"),
						Description: sl.String("Use if your host operating system is Linux.")},
				}, nil)

				FakeNetworkManager.GetClosingPodsReturns([]datatypes.Network_Pod{
					{
						DatacenterLongName: sl.String("MON01 - Montreal"),
						Name:               sl.String("mon01.pod02"),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("AMS03 - Amsterdam"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ams01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("closing soon: mon01.pod02"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("20 - 39 GBs"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("20_39_GBS"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("20 GB Storage Space"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("20_GB_PERFORMANCE_STORAGE_SPACE"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Available Snapshot Size (GB)"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("10000-12000"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500,600,700,1000,2000,4000"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("LINUX"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Use if your host operating system is Linux."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Location Conflicts"))
			})

			It("return no error with flag --prices", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--prices")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Prices"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Tier"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.25"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hourly/Monthly"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.002000/0.300000"))
			})
		})
	})
})

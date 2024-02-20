package managers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Storage Utils", func() {
	var (
		fakeSLSession  *session.Session
		storageManager managers.StorageManager
		productPackage datatypes.Product_Package
		err            error
	)
	sizes := []int{20, 40, 80, 100, 250, 1000, 2000, 4000, 8000, 12000}
	snapshotSize := []int{5, 10, 20, 40, 80, 100, 250, 1000, 2000, 4000}
	tiers := []float64{0.25, 2, 4, 10}
	iops := map[int][]int{
		20:    []int{100, 1000},
		40:    []int{100, 2000},
		80:    []int{100, 4000},
		100:   []int{100, 6000},
		250:   []int{100, 6000},
		500:   []int{100, 6000},
		1000:  []int{100, 6000},
		2000:  []int{200, 6000},
		4000:  []int{300, 6000},
		8000:  []int{500, 6000},
		12000: []int{1000, 6000}}

	BeforeEach(func() {
		filenames := []string{"getAllObjects_saas", }
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
		storageManager = managers.NewStorageManager(fakeSLSession)
		_ = storageManager
	})

	Describe("FindSaasEnduranceSpacePrice", func() {
		It("", func() {
			productPackage, err = managers.GetPackage(services.GetProductPackageService(fakeSLSession), managers.SaaS_Category)
			Expect(err).ToNot(HaveOccurred())
			for _, t := range tiers {
				for _, s := range sizes {
					//4TB max volume limit for 10 IOPS
					if s == 8000 && t == 10 {
						continue
					}
					if s == 12000 && t == 10 {
						continue
					}
					price, err := managers.FindSaasEnduranceSpacePrice(productPackage, s, t)
					Expect(err).ToNot(HaveOccurred())
					Expect(price).NotTo(BeNil())
				}
			}
		})
	})

	Describe("FindSaasEnduranceTierPrice", func() {
		It("", func() {
			productPackage, err = managers.GetPackage(services.GetProductPackageService(fakeSLSession), managers.SaaS_Category)
			Expect(err).ToNot(HaveOccurred())
			for _, t := range tiers {
				price, err := managers.FindSaasEnduranceTierPrice(productPackage, t)
				Expect(err).ToNot(HaveOccurred())
				Expect(price).NotTo(BeNil())
			}
		})
	})

	Describe("FindSaasSnapshotSpacePrice", func() {
		It("endurance", func() {
			productPackage, err = managers.GetPackage(services.GetProductPackageService(fakeSLSession), managers.SaaS_Category)
			Expect(err).ToNot(HaveOccurred())
			for _, t := range tiers {
				for _, s := range snapshotSize {
					price, err := managers.FindSaasSnapshotSpacePrice(productPackage, s, t, 0)
					Expect(err).ToNot(HaveOccurred())
					Expect(price).NotTo(BeNil())
				}
			}
		})
		It("performance", func() {
			productPackage, err = managers.GetPackage(services.GetProductPackageService(fakeSLSession), managers.SaaS_Category)
			Expect(err).ToNot(HaveOccurred())
			for _, s := range snapshotSize {
				if s == 5 || s == 10 {
					continue
				}
				iopsRange := iops[s]
				for i := iopsRange[0]; i < iopsRange[1]; i = i + 100 {
					price, err := managers.FindSaasSnapshotSpacePrice(productPackage, s, 0, i)
					Expect(err).ToNot(HaveOccurred())
					Expect(price).NotTo(BeNil())
				}
			}
		})
	})

	Describe("FindSaasPerformanceSpacePrice", func() {
		It("", func() {
			productPackage, err = managers.GetPackage(services.GetProductPackageService(fakeSLSession), managers.SaaS_Category)
			Expect(err).ToNot(HaveOccurred())
			for _, s := range sizes {
				price, err := managers.FindSaasPerformanceSpacePrice(productPackage, s)
				Expect(err).ToNot(HaveOccurred())
				Expect(price).NotTo(BeNil())

			}
		})
	})

	Describe("FindSaasPerformanceIopsPrice", func() {
		It("", func() {
			productPackage, err = managers.GetPackage(services.GetProductPackageService(fakeSLSession), managers.SaaS_Category)
			Expect(err).ToNot(HaveOccurred())
			for _, s := range sizes {
				iopsRange := iops[s]
				for i := iopsRange[0]; i < iopsRange[1]; i = i + 100 {
					price, err := managers.FindSaasPerformanceIopsPrice(productPackage, s, i)
					Expect(err).ToNot(HaveOccurred())
					Expect(price).NotTo(BeNil())
				}
			}
		})
	})

	Describe("FindSaasReplicationPrice", func() {
		It("endurance", func() {
			productPackage, err = managers.GetPackage(services.GetProductPackageService(fakeSLSession), managers.SaaS_Category)
			Expect(err).ToNot(HaveOccurred())
			for _, t := range tiers {
				price, err := managers.FindSaasReplicationPrice(productPackage, t, 0)
				Expect(err).ToNot(HaveOccurred())
				Expect(price).NotTo(BeNil())

			}
		})
		It("performance", func() {
			productPackage, err = managers.GetPackage(services.GetProductPackageService(fakeSLSession), managers.SaaS_Category)
			Expect(err).ToNot(HaveOccurred())
			for _, s := range snapshotSize {
				if s == 5 || s == 10 {
					continue
				}
				iopsRange := iops[s]
				for i := iopsRange[0]; i < iopsRange[1]; i = i + 100 {
					price, err := managers.FindSaasReplicationPrice(productPackage, 0, i)
					Expect(err).ToNot(HaveOccurred())
					Expect(price).NotTo(BeNil())
				}
			}
		})
	})

	Describe("Test PrepareDuplicateOrderObject", func() {
		BeforeEach(func() {
			productPackage, _ = managers.GetPackage(services.GetProductPackageService(fakeSLSession), managers.SaaS_Category)
		})
		It("Test Empty Config", func() {
			configuration := managers.DuplicateOrderConfig{
				VolumeType: "block",
			}
			originalVolume, get_volume_err := storageManager.GetVolumeDetails("block", 12345, "")
			Expect(get_volume_err).ToNot(HaveOccurred())
			result, err := managers.PrepareDuplicateOrderObject(productPackage, originalVolume, configuration)	
			Expect(err).ToNot(HaveOccurred())
			Expect(*result.VolumeSize).To(Equal(20))
		})
		It("Test 10g sized duplicateSnapshotSize", func() {
			configuration := managers.DuplicateOrderConfig{
				VolumeType: "block",
				DuplicateSnapshotSize: 10,
			}
			originalVolume, get_volume_err := storageManager.GetVolumeDetails("block", 12345, "")
			Expect(get_volume_err).ToNot(HaveOccurred())
			result, err := managers.PrepareDuplicateOrderObject(productPackage, originalVolume, configuration)	
			Expect(err).ToNot(HaveOccurred())
			Expect(*result.VolumeSize).To(Equal(20))
			snapshot_found := false
			for _, price := range result.Container_Product_Order.Prices {
				for _, category := range price.Categories {
					if *category.CategoryCode == "storage_snapshot_space" {
						snapshot_found = true
					}
				}
			}
			Expect(snapshot_found).To(Equal(true))
		})
		It("Test Default sized duplicateSnapshotSize", func() {
			configuration := managers.DuplicateOrderConfig{
				VolumeType: "block",
				DuplicateSnapshotSize: -1,
			}
			originalVolume, get_volume_err := storageManager.GetVolumeDetails("block", 12345, "")
			Expect(get_volume_err).ToNot(HaveOccurred())
			result, err := managers.PrepareDuplicateOrderObject(productPackage, originalVolume, configuration)	
			Expect(err).ToNot(HaveOccurred())
			Expect(*result.VolumeSize).To(Equal(20))
			snapshot_found := false
			expected_price := datatypes.Product_Item_Price{}
			for _, price := range result.Container_Product_Order.Prices {
				for _, category := range price.Categories {
					if *category.CategoryCode == "storage_snapshot_space" {
						snapshot_found = true
						expected_price = price
					}
				}
			}
			Expect(snapshot_found).To(Equal(true))
			// Make sure we got the right snapshot space price.
			Expect(*expected_price.CapacityRestrictionMaximum).To(Equal("300"))
		})
		It("Test Zero sized duplicateSnapshotSize", func() {
			configuration := managers.DuplicateOrderConfig{
				VolumeType: "block",
				DuplicateSnapshotSize: 0,
			}
			originalVolume, get_volume_err := storageManager.GetVolumeDetails("block", 12345, "")
			Expect(get_volume_err).ToNot(HaveOccurred())
			result, err := managers.PrepareDuplicateOrderObject(productPackage, originalVolume, configuration)	
			Expect(err).ToNot(HaveOccurred())
			Expect(*result.VolumeSize).To(Equal(20))
			snapshot_found := false
			for _, price := range result.Container_Product_Order.Prices {
				for _, category := range price.Categories {
					if *category.CategoryCode == "storage_snapshot_space" {
						snapshot_found = true
					}
				}
			}
			// SNAPSHOT should NOT be in here
			Expect(snapshot_found).To(Equal(false))
		})
	})
	Describe("Issues3190", func() {
		endurance_key := "STORAGE_SPACE_FOR_2_IOPS_PER_GB_NEW_AND_IMPROVED"
		performance_key := "1_100_GBS_NEW_AND_IMPROVED"
		maxCapacity := "100"
		minCapacity := "1"
		priceId := 99
		category := managers.Space_Category
		itemCategory := datatypes.Product_Item_Category{
			CategoryCode: &category,
		}
		productPackage := datatypes.Product_Package{
			Items: []datatypes.Product_Item{
				datatypes.Product_Item{
					KeyName: &endurance_key,
					CapacityMaximum: &maxCapacity,
					CapacityMinimum: &minCapacity,
					ItemCategory: &itemCategory,
					Prices: []datatypes.Product_Item_Price{
						datatypes.Product_Item_Price{
							Id: &priceId,
							LocationGroupId: nil,
							Categories: []datatypes.Product_Item_Category{
								itemCategory,
							},
						},
					},
				},
			},
		}
		Context("FindSaasEnduranceSpacePrice", func() {
			It("Checks for special item keynames", func() {
				price, err := managers.FindSaasEnduranceSpacePrice(productPackage, 55, 2)
				Expect(err).ToNot(HaveOccurred())
				Expect(*price.Id).To(Equal(99))
			})
		})
		Context("FindSaasPerformanceSpacePrice", func() {
			It("Checks for special item keynames", func() {
				productPackage.Items[0].KeyName = &performance_key
				price, err := managers.FindSaasPerformanceSpacePrice(productPackage, 55)
				Expect(err).ToNot(HaveOccurred())
				Expect(*price.Id).To(Equal(99))
			})
		})
	})
})

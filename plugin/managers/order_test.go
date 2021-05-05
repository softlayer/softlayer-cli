package managers_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Order", func() {
	var (
		fakeSLSession *session.Session
		OrderManager  managers.OrderManager
	)
	BeforeEach(func() {
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		OrderManager = managers.NewOrderManager(fakeSLSession)
	})

	Describe("ListCategories", func() {
		BeforeEach(func() {
			filenames := []string{
				"SoftLayer_Product_Package_getConfiguration",
				"SoftLayer_Product_Package_getAllObjects",
			}
			fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
			OrderManager = managers.NewOrderManager(fakeSLSession)
		})

		Context("ListCategories given a packageKeyname", func() {
			It("Return Categories", func() {
				categories, err := OrderManager.ListCategories("packageKeyname")
				Expect(err).ToNot(HaveOccurred())
				Expect(*categories[0].ItemCategory.Name).To(Equal("Server"))
				Expect(*categories[0].ItemCategory.CategoryCode).To(Equal("server"))
			})
		})
	})

	Describe("ListItems", func() {
		Context("ListItems given a ListItems, keyword and category", func() {
			It("Return no error", func() {
				Items, err := OrderManager.ListItems("ListItems", "keyword", "category")
				Expect(err).ToNot(HaveOccurred())
				//bug current the returned list is empty
				Expect(*Items[0].KeyName).To(Equal("CDN_25_GB_STORAGE"))

			})
		})
	})

	Describe("ListPackage", func() {
		Context("ListPackage given a keyword, packageType", func() {
			It("Return no error", func() {
				Packages, err := OrderManager.ListPackage("keyword", "packageType")
				Expect(err).ToNot(HaveOccurred())
				//TODO current the returned list is empty
				Expect(*Packages[0].Name).To(Equal("Additional Products"))
				Expect(*Packages[0].KeyName).To(Equal("ADDITIONAL_PRODUCTS"))
			})
		})
	})

	Describe("PackageLocation", func() {
		Context("PackageLocation given packageKeyname", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Product_Package_getRegionss",
					"SoftLayer_Product_Package_getAllObjects",
					"SoftLayer_Product_Package_getRegions",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				OrderManager = managers.NewOrderManager(fakeSLSession)
			})
			It("Return no error", func() {
				PackageLocation, err := OrderManager.PackageLocation("packageKeyname")
				Expect(err).ToNot(HaveOccurred())
				for _, region := range PackageLocation {
					for _, datacenter := range region.Locations {
						Expect(datacenter.Location.Id).To(Equal(nil))
						Expect(datacenter.Location.Name).To(Equal(nil))
						Expect(region.Description).To(Equal(nil))
						Expect(region.Keyname).To(Equal(nil))
					}
				}

			})
		})
	})

	Describe("GenerateOrder", func() {
		BeforeEach(func() {
			filenames := []string{
				"SoftLayer_Location_Datacenter_getDatacenters",
				"SoftLayer_Product_Package_getAllObjects",
				"SoftLayer_Product_Package_getAccountRestrictedActivePresets",
				"SoftLayer_Product_Package_getActivePresets",
				"SoftLayer_Product_Package_getItems",
				"SoftLayer_Product_Order_verifyOrder",
			}
			fakeSLSession = testhelpers.NewFakeSoftlayerLocationSession(filenames)
			OrderManager = managers.NewOrderManager(fakeSLSession)
		})
		Context("VerifyPlaceOrder given packageKeyname, location , itemKeynames , complexType , hourly , presetKeyname , extras", func() {
			It("Return no error", func() {
				var i interface{}
				i = &datatypes.Container_Product_Order_Network_Subnet{}
				order, err := OrderManager.GenerateOrder("packageKeyname", "location", []string{"CDN_25_GB_STORAGE"}, "Softlayer_Container_Product_Order_Network_Subnet", false, "presetKeyname", i, 0)
				Expect(err).ToNot(HaveOccurred())
				orderValue := reflect.ValueOf(order).Elem()
				Expect(orderValue.FieldByName("PackageId").Elem().Interface().(int)).To(Equal(865))
				Expect(orderValue.FieldByName("Location").Elem().Interface().(string)).To(Equal("265592"))
				Expect(orderValue.FieldByName("Quantity").Elem().Interface().(int)).To(Equal(1))
				Expect(orderValue.FieldByName("UseHourlyPricing").Elem().Interface().(bool)).To(Equal(false))
				Expect(orderValue.FieldByName("PresetId").Elem().Interface().(int)).To(Equal(785))
				Expect(orderValue.FieldByName("ComplexType").Elem().Interface().(string)).To(Equal("Softlayer_Container_Product_Order_Network_Subnet"))
				Expect(*orderValue.FieldByName("Prices").Interface().([]datatypes.Product_Item_Price)[0].Id).To(Equal(230623))

			})
		})
	})

	Describe("VerifyPlaceOrder", func() {
		BeforeEach(func() {
			filenames := []string{
				"SoftLayer_Location_Datacenter_getDatacenters",
				"SoftLayer_Product_Package_getAllObjects",
				"SoftLayer_Product_Package_getAccountRestrictedActivePresets",
				"SoftLayer_Product_Package_getActivePresets",
				"SoftLayer_Product_Package_getItems",
				"SoftLayer_Product_Order_verifyOrder",
			}
			fakeSLSession = testhelpers.NewFakeSoftlayerLocationSession(filenames)
			OrderManager = managers.NewOrderManager(fakeSLSession)
		})
		Context("VerifyPlaceOrder given packageKeyname, location , itemKeynames , complexType , hourly , presetKeyname , extras", func() {
			It("Return no error", func() {
				_, err := OrderManager.VerifyPlaceOrder("packageKeyname", "location", []string{"CDN_25_GB_STORAGE"}, "complexType", false, "presetKeyname", &datatypes.Container_Product_Order{}, 0)
				Expect(err).ToNot(HaveOccurred())

			})
		})
	})

	Describe("PlaceOrder", func() {
		BeforeEach(func() {
			filenames := []string{
				"SoftLayer_Location_Datacenter_getDatacenters",
				"SoftLayer_Product_Package_getAllObjects",
				"SoftLayer_Product_Package_getAccountRestrictedActivePresets",
				"SoftLayer_Product_Package_getActivePresets",
				"SoftLayer_Product_Package_getItems",
				"SoftLayer_Product_Order_placeOrder",
			}
			fakeSLSession = testhelpers.NewFakeSoftlayerLocationSession(filenames)
			OrderManager = managers.NewOrderManager(fakeSLSession)
		})
		Context("PlaceOrder given volume id", func() {
			It("Return no error", func() {
				orderPlace, err := OrderManager.PlaceOrder("packageKeyname", "location", []string{"CDN_25_GB_STORAGE"}, "complexType", false, "presetKeyname", &datatypes.Container_Product_Order{}, 0)
				Expect(err).ToNot(HaveOccurred())
				Expect(*orderPlace.OrderId).To(Equal(11493593))
				Expect(*orderPlace.PlacedOrder.Status).To(Equal("PENDING_AUTO_APPROVAL"))

			})
		})
	})

	Describe("PlaceQuote", func() {
		BeforeEach(func() {
			filenames := []string{
				"SoftLayer_Location_Datacenter_getDatacenters",
				"SoftLayer_Product_Package_getAllObjects",
				"SoftLayer_Product_Package_getAccountRestrictedActivePresets",
				"SoftLayer_Product_Package_getActivePresets",
				"SoftLayer_Product_Package_getItems",
				"SoftLayer_Product_Order_placeQuote",
			}
			fakeSLSession = testhelpers.NewFakeSoftlayerLocationSession(filenames)
			OrderManager = managers.NewOrderManager(fakeSLSession)
		})

		Context("PlaceQuote under current account", func() {
			It("Return no error", func() {
				placeQuote, err := OrderManager.PlaceQuote("packageKeyname", "tok02", []string{"CDN_25_GB_STORAGE"}, "complexType", "name", "presetKeyname", &datatypes.Container_Product_Order{}, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(*placeQuote.Quote.Id).To(Equal(2523413))
			})
		})
	})

	Describe("ListPreset", func() {
		Context("ListPreset given volume id", func() {
			It("Return the volume details and no error", func() {
				presets, err := OrderManager.ListPreset("packageKeyname", "keyword")
				Expect(err).ToNot(HaveOccurred())
				Expect(*presets[0].Name).To(Equal("DSilver 4110 96GB 1X960GB SSD SED NoRAID"))
				Expect(*presets[0].KeyName).To(Equal("DSILVER_4110_96GB_1X960GB_SSD_SED_NORAID"))

			})
		})
	})

	Describe("GetPresetbyKey", func() {
		BeforeEach(func() {
			filenames := []string{
				"SoftLayer_Product_Package_getAllObjects",
				"SoftLayer_Product_Package_getAccountRestrictedActivePresets",
				"SoftLayer_Product_Package_getActivePresets",
			}
			fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
			OrderManager = managers.NewOrderManager(fakeSLSession)
		})
		Context("GetPresetbyKey", func() {
			It("Return the order receipt and no error", func() {
				Preset, err := OrderManager.GetPresetbyKey("packageKeyname", "presetKeyname")
				Expect(err).ToNot(HaveOccurred())
				Expect(*Preset.Id).To(Equal(785))
			})
		})
	})

	Describe("GetPresetPrices", func() {
		BeforeEach(func() {
			filenames := []string{
				"SoftLayer_Product_Package_Preset_getObject",
			}
			fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
			OrderManager = managers.NewOrderManager(fakeSLSession)
		})
		Context("GetPresetbyKey", func() {
			It("Return the order receipt and no error", func() {
				presetItems, err := OrderManager.GetPresetPrices(0)
				Expect(err).ToNot(HaveOccurred())
				Expect(*presetItems.Id).To(Equal(645))
				Expect(*presetItems.PackageId).To(Equal(1035))
			})
		})
	})

	Describe("GetPackageByKey", func() {
		Context("get result from service json file", func() {
			It("Return no error", func() {
				Preset, err := OrderManager.GetPackageByKey("packagenamem", "mask")
				Expect(err).ToNot(HaveOccurred())
				Expect(*Preset.Id).To(Equal(865))
			})
		})
	})
	Describe("GetLocation", func() {
		Context("get result from service json file", func() {
			It("with location id", func() {
				Location, err := OrderManager.GetLocation("11")
				Expect(err).ToNot(HaveOccurred())
				Expect(Location).To(Equal("11"))
			})
			It("with wrong location name", func() {
				Location, err := OrderManager.GetLocation("testerror")
				Expect(err).To(HaveOccurred())
				Expect(Location).To(Equal(""))
			})
			It("with location name, but service json return a list of location.", func() {
				_, err := OrderManager.GetLocation("dal10")
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Describe("GetPriceIdList", func(){
		Items := []string{
			"PRIVATE_NETWORK_VLAN",
			"DOMAIN_INFO_1_YEAR",
		}
		BeforeEach(func() {
			fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
			OrderManager = managers.NewOrderManager(fakeSLSession)
		})
		Context("Price Id selection", func() {
			It("Basic Price selection", func(){
				prices, err := OrderManager.GetPriceIdList("Test", Items, 0)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(prices)).To(Equal(2))
				Expect(prices[0]).To(Equal(2727099)) // Term Price check
				Expect(prices[1]).To(Equal(26835))  // Normal Price (nil TermLength)
			})
			It("Capacity restriction checking", func(){
				Items[1] = "CDN_25_GB_STORAGE"
				prices, err := OrderManager.GetPriceIdList("Test", Items, 25)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(prices)).To(Equal(2))
				Expect(prices[1]).To(Equal(230623))  // Normal Price (nil TermLength)
			})
			It("Error: No matching Keyname", func(){
				Items[1] = "FAKE"
				_, err := OrderManager.GetPriceIdList("Test", Items, 25)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

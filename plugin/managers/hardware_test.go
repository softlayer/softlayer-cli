package managers_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("HardwareServerManager", func() {
	var (
		fakeSLSession   *session.Session
		fakeHandler     *testhelpers.FakeTransportHandler
		hardwareManager managers.HardwareServerManager
		productPackage  datatypes.Product_Package
		datacenter      datatypes.Location_Region
	)

	BeforeEach(func() {
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		fakeHandler = testhelpers.GetSessionHandler(fakeSLSession)
		hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
	})

	AfterEach(func() {
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})
	Describe("Cancel hardware", func() {
		Context("Cancel hardware with billing id not found", func() {
			BeforeEach(func() {
				filenames := []string{"getObject_missingBillingItem"}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
			})
			It("It returns error", func() {
				err := hardwareManager.CancelHardware(123, "abd", "", true)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("No billing item found for hardware"))
			})
		})
		Context("Cancel hardware succeed", func() {
			BeforeEach(func() {
				filenames := []string{"cancelItem_hardware"}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
			})
			It("It does not returns error", func() {
				err := hardwareManager.CancelHardware(123, "abd", "", true)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("List hardware", func() {
		Context("List hardware", func() {
			It("it returns hardware", func() {
				hws, err := hardwareManager.ListHardware(nil, 0, 0, "", "", "", 0, "", "", "", 0, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(hws)).To(Equal(2))
				apiCalls := fakeHandler.ApiCallLogs
				Expect(len(apiCalls)).To(Equal(1))
				Expect(apiCalls[0].Service).To(Equal("SoftLayer_Account"))
				Expect(apiCalls[0].Method).To(Equal("getHardware"))
			})
		})
		Context("List Hardware all options", func() {
			It("Returns a hardware list", func() {
				hws, err := hardwareManager.ListHardware(
					[]string{"tag1"}, 1, 2, "hostnametest", "testdomain", "dctest", 10, "1.2.3.4", "5.6.7.8", "testuser", 55, "mask[id]")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(hws)).To(Equal(2))
				apiCalls := fakeHandler.ApiCallLogs
				Expect(len(apiCalls)).To(Equal(1))
				Expect(apiCalls[0].Service).To(Equal("SoftLayer_Account"))
				slOptions := apiCalls[0].Options
				// Check to make sure all object filters get set properly.
				Expect(slOptions.Filter).To(ContainSubstring(`"id":{"operation":"orderBy","options":[{"name":"sort","value":["DESC"]}]}`))
				Expect(slOptions.Filter).To(ContainSubstring(`id":{"operation":55}`))
				Expect(slOptions.Filter).To(ContainSubstring(`"userRecord":{"username":{"operation":"testuser"}}`))
				Expect(slOptions.Filter).To(ContainSubstring(`"datacenter":{"name":{"operation":"dctest"}}`))
				Expect(slOptions.Filter).To(ContainSubstring(`"domain":{"operation":"testdomain"}`))
				Expect(slOptions.Filter).To(ContainSubstring(`"hostname":{"operation":"hostnametest"}`))
				Expect(slOptions.Filter).To(ContainSubstring(`"memoryCapacity":{"operation":2}`))
				Expect(slOptions.Filter).To(ContainSubstring(`"networkComponents":{"maxSpeed":{"operation":10}}`))
				Expect(slOptions.Filter).To(ContainSubstring(`"primaryBackendIpAddress":{"operation":"5.6.7.8"}`))
				Expect(slOptions.Filter).To(ContainSubstring(`"primaryIpAddress":{"operation":"1.2.3.4"}`))
				Expect(slOptions.Filter).To(ContainSubstring(`"processorPhysicalCoreAmount":{"operation":1}`))
				Expect(slOptions.Filter).To(ContainSubstring(`"tagReferences":{"tag":{"name":{"operation":"in","options":[{"name":"data","value":["tag1"]}]}}`))

			})
		})
	})

	Describe("Get hardware", func() {
		Context("get hardware", func() {
			It("it returns hardware", func() {
				hw, err := hardwareManager.GetHardware(218027, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(*hw.Id).To(Equal(218027))
			})
		})
	})

	Describe("Reload OS", func() {
		Context("reload OS", func() {
			It("it returns nil", func() {
				err := hardwareManager.Reload(123456, "", nil, false, false)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("Rescue OS", func() {
		Context("Rescue OS", func() {
			It("it returns nil", func() {
				err := hardwareManager.Rescure(123456)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("GetCancellationReasons", func() {
		Context("", func() {
			It("it returns map", func() {
				reasons := hardwareManager.GetCancellationReasons()
				Expect(reasons).To(Equal(map[string]string{
					"unneeded":        "No longer needed",
					"closing":         "Business closing down",
					"cost":            "Server / Upgrade Costs",
					"migrate_larger":  "Migrating to larger server",
					"migrate_smaller": "Migrating to smaller server",
					"datacenter":      "Migrating to a different SoftLayer datacenter",
					"performance":     "Network performance / latency",
					"support":         "Support response / timing",
					"sales":           "Sales process / upgrades",
					"moving":          "Moving to competitor",
				}))
			})
		})
	})

	Describe("GetPackage", func() {
		Context("not found Package", func() {
			BeforeEach(func() {
				filenames := []string{"getAllObjects_hardwarenotfound"}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
			})
			It("it returns zero hardware package", func() {
				productPackage, err := hardwareManager.GetPackage()
				Expect(err).To(HaveOccurred())
				Expect(productPackage.Id).To(BeNil())
				Expect(err.Error()).To(ContainSubstring("Ordering package is not found"))
			})
		})

		Context("get Package", func() {
			BeforeEach(func() {
				filenames := []string{"getAllObjects_hardware"}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
			})
			It("it returns hardware package", func() {
				productPackage, err := hardwareManager.GetPackage()
				Expect(err).NotTo(HaveOccurred())
				Expect(productPackage).NotTo(BeNil())
			})
		})
	})

	Describe("Edit", func() {
		Context("Edit all succeed", func() {
			It("it returns nil", func() {
				succeeds, msgs := hardwareManager.Edit(123456, "test-userdata", "test-hostname", "test-domain", "test-notes", "test-tags", 100, 100)
				Expect(succeeds).To(Equal([]bool{true, true, true, true, true, true, true}))
				Expect(msgs).NotTo(BeNil())
			})
		})
		Context("Edit set metadata fails", func() {
			It("it returns 1 error", func() {
				fakeHandler.AddApiError("SoftLayer_Hardware_Server", "setUserMetadata", 500, "BAD")
				succeeds, msgs := hardwareManager.Edit(123456, "test-userdata", "test-hostname", "test-domain", "test-notes", "test-tags", 100, 100)
				Expect(succeeds).To(Equal([]bool{false, true, true, true, true, true, true}))
				Expect(msgs).NotTo(BeNil())
				Expect(msgs[0]).To(ContainSubstring("BAD: BAD (HTTP 500)"))
			})
		})
		Context("Edit set tag fails", func() {
			It("it returns 1 error", func() {
				fakeHandler.AddApiError("SoftLayer_Hardware_Server", "setTags", 500, "BAD")
				succeeds, msgs := hardwareManager.Edit(123456, "test-userdata", "test-hostname", "test-domain", "test-notes", "test-tags", 100, 100)
				Expect(succeeds).To(Equal([]bool{true, false, true, true, true, true, true}))
				Expect(msgs).NotTo(BeNil())
				Expect(msgs[1]).To(ContainSubstring("BAD: BAD (HTTP 500)"))
			})
		})
		Context("Edit set hostname fails", func() {
			It("it returns 1 error", func() {
				fakeHandler.AddApiError("SoftLayer_Hardware_Server", "editObject", 500, "BAD")
				succeeds, msgs := hardwareManager.Edit(123456, "test-userdata", "test-hostname", "test-domain", "test-notes", "test-tags", 100, 100)
				Expect(succeeds).To(Equal([]bool{true, true, false, true, true}))
				Expect(msgs).NotTo(BeNil())
				Expect(msgs[2]).To(ContainSubstring("BAD: BAD (HTTP 500)"))
			})
		})

		Context("Edit set public port speed fails", func() {
			It("it returns 1 error", func() {
				fakeHandler.AddApiError("SoftLayer_Hardware_Server", "setPublicNetworkInterfaceSpeed", 500, "BAD")
				succeeds, msgs := hardwareManager.Edit(123456, "test-userdata", "test-hostname", "test-domain", "test-notes", "test-tags", 100, 100)
				Expect(succeeds).To(Equal([]bool{true, true, true, true, true, false, true}))
				Expect(msgs).NotTo(BeNil())
				Expect(msgs[5]).To(ContainSubstring("BAD: BAD (HTTP 500)"))
			})
		})

		Context("Edit set private port speed fails", func() {
			It("it returns 1 error", func() {
				fakeHandler.AddApiError("SoftLayer_Hardware_Server", "setPrivateNetworkInterfaceSpeed", 500, "BAD")
				succeeds, msgs := hardwareManager.Edit(123456, "test-userdata", "test-hostname", "test-domain", "test-notes", "test-tags", 100, 100)
				Expect(succeeds).To(Equal([]bool{true, true, true, true, true, true, false}))
				Expect(msgs).NotTo(BeNil())
				Expect(msgs[6]).To(ContainSubstring("BAD: BAD (HTTP 500)"))
			})
		})
	})

	Describe("UpdateFirmware", func() {
		Context("UpdateFirmware succeed", func() {
			It("it returns nil", func() {
				err := hardwareManager.UpdateFirmware(123456, true, true, true, true)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("UpdateFirmware fails", func() {
			It("it returns error", func() {
				fakeHandler.AddApiError("SoftLayer_Hardware_Server", "createFirmwareUpdateTransaction", 500, "FAILED")
				err := hardwareManager.UpdateFirmware(123456, true, true, true, true)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("FAILED: FAILED (HTTP 500)"))
			})
		})
	})

	Describe("GetCreateOptions", func() {
		BeforeEach(func() {
			filenames := []string{"getAllObjects_hardware"}
			fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
			hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
		})
		It("it returns valid options", func() {
			productPackage, err := hardwareManager.GetPackage()
			Expect(err).NotTo(HaveOccurred())
			result := hardwareManager.GetCreateOptions(productPackage)
			Expect(result[managers.KEY_LOCATIONS]["dal06"]).To(Equal("Dallas 6"))
			Expect(result[managers.KEY_SIZES]["D2620V4_64GB_2X1TB_SATA_RAID_1"]).To(Equal("Dual Xeon 2620v4, 64GB Ram, 2x1TB SATA disks, RAID1"))
			Expect(result[managers.KEY_OS]["UBUNTU_16_64"]).To(Equal("Ubuntu 16.04-64"))
			Expect(result[managers.KEY_PORT_SPEED]["100"]).To(Equal("100 Mbps Public & Private Network Uplinks"))
			Expect(result[managers.KEY_EXTRAS]["64_BLOCK_STATIC_PUBLIC_IPV6_ADDRESSES"]).To(Equal("/64 Block Static Public IPv6 Addresses"))
		})
	})

	Describe("GetDefaultPriceId", func() {
		Context("GetDefaultPriceId", func() {
			BeforeEach(func() {
				filenames := []string{"getAllObjects_hardware"}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
				productPackage, _ = hardwareManager.GetPackage()
				datacenter, _ = managers.GetLocation(productPackage, "dal10")
			})
			It("it returns int", func() {
				priceId, err := hardwareManager.GetDefaultPriceId(productPackage.Items, managers.DEFAULT_CATEGORIES[0], true, datacenter)
				Expect(err).NotTo(HaveOccurred())
				Expect(priceId).To(Equal(34807))
			})
			It("it returns int", func() {
				priceId, err := hardwareManager.GetDefaultPriceId(productPackage.Items, managers.DEFAULT_CATEGORIES[1], false, datacenter)
				Expect(err).NotTo(HaveOccurred())
				Expect(priceId).To(Equal(33483))
			})
			It("it returns int", func() {
				priceId, err := hardwareManager.GetDefaultPriceId(productPackage.Items, managers.DEFAULT_CATEGORIES[2], false, datacenter)
				Expect(err).NotTo(HaveOccurred())
				Expect(priceId).To(Equal(25014))
			})
		})
	})

	Describe("GetOSPriceId", func() {
		Context("GetOSPriceId", func() {
			BeforeEach(func() {
				filenames := []string{"getAllObjects_hardware"}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
				productPackage, _ = hardwareManager.GetPackage()
				datacenter, _ = managers.GetLocation(productPackage, "dal10")
			})
			It("it returns int", func() {
				priceId, err := hardwareManager.GetOSPriceId(productPackage.Items, "UBUNTU_16_64", datacenter)
				Expect(err).NotTo(HaveOccurred())
				Expect(priceId).To(Equal(175789))
			})
			It("it returns int", func() {
				priceId, err := hardwareManager.GetOSPriceId(productPackage.Items, "REDHAT_7_64", datacenter)
				Expect(err).NotTo(HaveOccurred())
				Expect(priceId).To(Equal(49073))
			})
			It("it returns int", func() {
				priceId, err := hardwareManager.GetOSPriceId(productPackage.Items, "WIN_2016-STD_64", datacenter)
				Expect(err).NotTo(HaveOccurred())
				Expect(priceId).To(Equal(179921))
			})
		})
	})

	Describe("GetBandwidthPriceId", func() {
		Context("GetBandwidthPriceId", func() {
			BeforeEach(func() {
				filenames := []string{"getAllObjects_hardware"}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
				productPackage, _ = hardwareManager.GetPackage()
				datacenter, _ = managers.GetLocation(productPackage, "dal10")
			})
			It("it returns int", func() {
				priceId, err := hardwareManager.GetBandwidthPriceId(productPackage.Items, true, true, datacenter)
				Expect(err).NotTo(HaveOccurred())
				Expect(priceId).To(Equal(34183))
			})
			It("it returns int", func() {
				priceId, err := hardwareManager.GetBandwidthPriceId(productPackage.Items, false, false, datacenter)
				Expect(err).NotTo(HaveOccurred())
				Expect(priceId).To(Equal(50233))
			})
		})
	})

	Describe("GetPortSpeedPriceId", func() {
		BeforeEach(func() {
			filenames := []string{"getAllObjects_hardware"}
			fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
			hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
			productPackage, _ = hardwareManager.GetPackage()
			datacenter, _ = managers.GetLocation(productPackage, "dal10")
		})
		It("it returns int", func() {
			priceId, err := hardwareManager.GetPortSpeedPriceId(productPackage.Items, 100, false, datacenter)
			Expect(err).NotTo(HaveOccurred())
			Expect(priceId).To(Equal(26737))
		})
		It("it returns int", func() {
			priceId, err := hardwareManager.GetPortSpeedPriceId(productPackage.Items, 100, true, datacenter)
			Expect(err).NotTo(HaveOccurred())
			Expect(priceId).To(Equal(23787))
		})
	})

	Describe("GetExtraPriceId", func() {
		Context("GetExtraPriceId", func() {
			BeforeEach(func() {
				filenames := []string{"getAllObjects_hardware"}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
				productPackage, _ = hardwareManager.GetPackage()
				datacenter, _ = managers.GetLocation(productPackage, "dal10")
			})
			It("it returns int", func() {
				priceId, err := hardwareManager.GetExtraPriceId(productPackage.Items, "1_IPV6_ADDRESS", true, datacenter)
				Expect(err).NotTo(HaveOccurred())
				Expect(priceId).To(Equal(29403))
			})
			It("it returns int", func() {
				priceId, err := hardwareManager.GetExtraPriceId(productPackage.Items, "64_BLOCK_STATIC_PUBLIC_IPV6_ADDRESSES", true, datacenter)
				Expect(err).NotTo(HaveOccurred())
				Expect(priceId).To(Equal(26340))
			})
			It("it returns int", func() {
				priceId, err := hardwareManager.GetExtraPriceId(productPackage.Items, "8_PUBLIC_IP_ADDRESSES", true, datacenter)
				Expect(err).NotTo(HaveOccurred())
				Expect(priceId).To(Equal(28968))
			})
		})
	})

	Describe("GenerateCreateTemplate", func() {
		Context("GenerateCreateTemplate", func() {
			BeforeEach(func() {
				filenames := []string{"getAllObjects_hardware"}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
				productPackage, _ = hardwareManager.GetPackage()
				datacenter, _ = managers.GetLocation(productPackage, "dal10")
			})
			It("it returns template", func() {
				params := map[string]interface{}{
					"size":           "D2620V4_128GB_2X800GB_SSD_RAID_1_K80_GPU2",
					"hostname":       "ibmcloud-cli",
					"domain":         "ibm.com",
					"osName":         "UBUNTU_16_64",
					"datacenter":     "dal10",
					"postInstallURL": "https://xxx/install.sh",
					"portSpeed":      1000,
					"billing":        "hourly",
					"noPublic":       false,
				}
				order, err := hardwareManager.GenerateCreateTemplate(productPackage, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(*order.Hardware[0].Hostname).To(Equal("ibmcloud-cli"))
				Expect(*order.Hardware[0].Domain).To(Equal("ibm.com"))
				Expect(order.ProvisionScripts[0]).To(Equal("https://xxx/install.sh"))
				Expect(*order.UseHourlyPricing).To(BeTrue())
				Expect(len(order.Prices)).To(Equal(6))
			})
			It("it returns template", func() {
				params := map[string]interface{}{
					"size":           "D2620V4_128GB_2X800GB_SSD_RAID_1_K80_GPU2",
					"hostname":       "ibmcloud-cli",
					"domain":         "ibm.com",
					"osName":         "UBUNTU_16_64",
					"datacenter":     "dal10",
					"postInstallURL": "https://xxx/install.sh",
					"portSpeed":      1000,
					"billing":        "hourly",
					"noPublic":       false,
					"sshKeys":        []int{123, 234},
					"extras":         []string{"4_PUBLIC_IP_ADDRESSES"},
				}
				order, err := hardwareManager.GenerateCreateTemplate(productPackage, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(*order.Hardware[0].Hostname).To(Equal("ibmcloud-cli"))
				Expect(*order.Hardware[0].Domain).To(Equal("ibm.com"))
				Expect(order.ProvisionScripts[0]).To(Equal("https://xxx/install.sh"))
				Expect(*order.UseHourlyPricing).To(BeTrue())
				Expect(len(order.Prices)).To(Equal(7))
			})
		})
	})

	Describe("Place order", func() {
		Context("Place order", func() {
			It("it returns receipt", func() {
				orderTemplate := datatypes.Container_Product_Order{}
				orderReceipt, err := hardwareManager.PlaceOrder(orderTemplate)
				Expect(err).NotTo(HaveOccurred())
				Expect(orderReceipt).NotTo(BeNil())
			})
		})
	})

	Describe("Verify order", func() {
		Context("verify order", func() {
			It("it returns order", func() {
				orderTemplate := datatypes.Container_Product_Order{}
				order, err := hardwareManager.VerifyOrder(orderTemplate)
				Expect(err).NotTo(HaveOccurred())
				Expect(order).NotTo(BeNil())
			})
		})
	})

	Describe("Toggle IPMI", func() {
		Context("Enable IPMI", func() {
			It("should return success", func() {
				err := hardwareManager.ToggleIPMI(123456, true)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("When there is an error", func() {
			It("should return error", func() {
				fakeHandler.AddApiError("SoftLayer_Hardware_Server", "toggleManagementInterface", 500, "IPMI ERROR")
				err := hardwareManager.ToggleIPMI(123456, false)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("IPMI ERROR: IPMI ERROR (HTTP 500)"))
			})
		})
	})
	Describe("GetBandwidthData Tests", func() {
		var (
			startTime time.Time
			endTime   time.Time
		)
		BeforeEach(func() {
			startTime, _ = time.Parse("2006-01-02", "2021-01-01")
			endTime, _ = time.Parse("2006-01-02", "2021-02-01")
		})
		Context("Test Happy Path", func() {
			It("Tests API is called properly", func() {
				data, err := hardwareManager.GetBandwidthData(12345, startTime, endTime, 300)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(data)).To(Equal(12))
				Expect(*data[0].Type).To(Equal("cpu0"))
			})
		})
	})
})

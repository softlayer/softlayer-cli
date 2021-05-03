package managers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("HardwareServerManager", func() {
	var (
		fakeSLSession   *session.Session
		hardwareManager managers.HardwareServerManager
		productPackage  datatypes.Product_Package
		datacenter      datatypes.Location_Region
	)

	BeforeEach(func() {
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
	})

	Describe("Cancel hardware", func() {
		Context("Cancel hardware with billing id not found", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Hardware_Server_getObject_missingBillingItem",
				}
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
				filenames := []string{
					"SoftLayer_Hardware_Server_getObject",
					"SoftLayer_Billing_Item_cancelItem_hardware",
				}
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
			BeforeEach(func() {
				fakeSLSession = testhelpers.NewFakeSoftlayerPagnationSession(nil)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
			})
			It("it returns hardware", func() {
				hws, err := hardwareManager.ListHardware(nil, 0, 0, "", "", "", 0, "", "", "", 0, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(hws)).To(Equal(2))
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
				filenames := []string{
					"SoftLayer_Product_Package_getAllObjects_hardwarenotfound",
				}
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
				filenames := []string{
					"SoftLayer_Product_Package_getAllObjects_hardware",
				}
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
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Hardware_Server_setUserMetadata",
					"SoftLayer_Hardware_Server_setTags",
					"SoftLayer_Hardware_Server_editObject",
					"SoftLayer_Hardware_Server_setPublicNetworkInterfaceSpeed",
					"SoftLayer_Hardware_Server_setPrivateNetworkInterfaceSpeed",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
			})
			It("it returns nil", func() {
				succeeds, msgs := hardwareManager.Edit(123456, "test-userdata", "test-hostname", "test-domain", "test-notes", "test-tags", 100, 100)
				Expect(succeeds).To(Equal([]bool{true, true, true, true, true, true, true}))
				Expect(msgs).NotTo(BeNil())
			})
		})
		Context("Edit set metadata fails", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Hardware_Server_setUserMetadata_error",
					"SoftLayer_Hardware_Server_setTags",
					"SoftLayer_Hardware_Server_editObject",
					"SoftLayer_Hardware_Server_setPublicNetworkInterfaceSpeed",
					"SoftLayer_Hardware_Server_setPrivateNetworkInterfaceSpeed",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
			})
			It("it returns 1 error", func() {
				succeeds, msgs := hardwareManager.Edit(123456, "test-userdata", "test-hostname", "test-domain", "test-notes", "test-tags", 100, 100)
				Expect(succeeds).To(Equal([]bool{false, true, true, true, true, true, true}))
				Expect(msgs).NotTo(BeNil())
			})
		})
		Context("Edit set tag fails", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Hardware_Server_setTags_error",
					"SoftLayer_Hardware_Server_setUserMetadata",
					"SoftLayer_Hardware_Server_editObject",
					"SoftLayer_Hardware_Server_setPublicNetworkInterfaceSpeed",
					"SoftLayer_Hardware_Server_setPrivateNetworkInterfaceSpeed",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
			})
			It("it returns 1 error", func() {
				succeeds, msgs := hardwareManager.Edit(123456, "test-userdata", "test-hostname", "test-domain", "test-notes", "test-tags", 100, 100)
				Expect(succeeds).To(Equal([]bool{true, false, true, true, true, true, true}))
				Expect(msgs).NotTo(BeNil())
			})
		})
		Context("Edit set hostname fails", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Hardware_Server_setTags",
					"SoftLayer_Hardware_Server_setUserMetadata",
					"SoftLayer_Hardware_Server_editObject_error",
					"SoftLayer_Hardware_Server_setPublicNetworkInterfaceSpeed",
					"SoftLayer_Hardware_Server_setPrivateNetworkInterfaceSpeed",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
			})
			It("it returns 1 error", func() {
				succeeds, msgs := hardwareManager.Edit(123456, "test-userdata", "test-hostname", "test-domain", "test-notes", "test-tags", 100, 100)
				Expect(succeeds).To(Equal([]bool{true, true, false, true, true}))
				Expect(msgs).NotTo(BeNil())
			})
		})

		Context("Edit set public port speed fails", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Hardware_Server_setTags",
					"SoftLayer_Hardware_Server_setUserMetadata",
					"SoftLayer_Hardware_Server_editObject",
					"SoftLayer_Hardware_Server_setPublicNetworkInterfaceSpeed_error",
					"SoftLayer_Hardware_Server_setPrivateNetworkInterfaceSpeed",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
			})
			It("it returns 1 error", func() {
				succeeds, msgs := hardwareManager.Edit(123456, "test-userdata", "test-hostname", "test-domain", "test-notes", "test-tags", 100, 100)
				Expect(succeeds).To(Equal([]bool{true, true, true, true, true, false, true}))
				Expect(msgs).NotTo(BeNil())
			})
		})

		Context("Edit set private port speed fails", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Hardware_Server_setTags",
					"SoftLayer_Hardware_Server_setUserMetadata",
					"SoftLayer_Hardware_Server_editObject",
					"SoftLayer_Hardware_Server_setPublicNetworkInterfaceSpeed",
					"SoftLayer_Hardware_Server_setPrivateNetworkInterfaceSpeed_error",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
			})
			It("it returns 1 error", func() {
				succeeds, msgs := hardwareManager.Edit(123456, "test-userdata", "test-hostname", "test-domain", "test-notes", "test-tags", 100, 100)
				Expect(succeeds).To(Equal([]bool{true, true, true, true, true, true, false}))
				Expect(msgs).NotTo(BeNil())
			})
		})
	})

	Describe("UpdateFirmware", func() {
		Context("UpdateFirmware succeed", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Hardware_Server_createFirmwareUpdateTransaction",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
			})
			It("it returns nil", func() {
				err := hardwareManager.UpdateFirmware(123456, true, true, true, true)
				Expect(err).NotTo(HaveOccurred())
			})
			It("it returns nil", func() {
				err := hardwareManager.UpdateFirmware(123456, true, false, false, false)
				Expect(err).NotTo(HaveOccurred())
			})
			It("it returns nil", func() {
				err := hardwareManager.UpdateFirmware(123456, false, true, false, false)
				Expect(err).NotTo(HaveOccurred())
			})
			It("it returns nil", func() {
				err := hardwareManager.UpdateFirmware(123456, false, false, true, false)
				Expect(err).NotTo(HaveOccurred())
			})
			It("it returns nil", func() {
				err := hardwareManager.UpdateFirmware(123456, false, false, false, true)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("UpdateFirmware fails", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Hardware_Server_createFirmwareUpdateTransaction_error",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)
			})
			It("it returns error", func() {
				err := hardwareManager.UpdateFirmware(123456, true, true, true, true)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetCreateOptions", func() {
		BeforeEach(func() {
			filenames := []string{
				"SoftLayer_Product_Package_getAllObjects_hardware",
			}
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
				filenames := []string{
					"SoftLayer_Product_Package_getAllObjects_hardware",
				}
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
				filenames := []string{
					"SoftLayer_Product_Package_getAllObjects_hardware",
				}
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
				filenames := []string{
					"SoftLayer_Product_Package_getAllObjects_hardware",
				}
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
			filenames := []string{
				"SoftLayer_Product_Package_getAllObjects_hardware",
			}
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
				filenames := []string{
					"SoftLayer_Product_Package_getAllObjects_hardware",
				}
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
				filenames := []string{
					"SoftLayer_Product_Package_getAllObjects_hardware",
				}
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
		Context("When enable IPMI", func() {
			BeforeEach(func() {
				filenames := []string{
					"hardware/toggle_ipmi_enabled",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)

			})

			It("should return success", func() {
				err := hardwareManager.ToggleIPMI(123456, true)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("When disable IPMI", func() {
			BeforeEach(func() {
				filenames := []string{
					"hardware/toggle_ipmi_disabled",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)

			})

			It("should return success", func() {
				err := hardwareManager.ToggleIPMI(123456, false)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("When there is an error", func() {
			BeforeEach(func() {
				filenames := []string{
					"hardware/toggle_ipmi_error",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				hardwareManager = managers.NewHardwareServerManager(fakeSLSession)

			})
			It("should return error", func() {
				err := hardwareManager.ToggleIPMI(123456, false)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})

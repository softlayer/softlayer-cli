package managers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("FirewallManager", func() {
	var (
		fakeSLSession *session.Session
		fwManager     managers.FirewallManager
	)

	BeforeEach(func() {
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		fwManager = managers.NewFirewallManager(fakeSLSession)
	})

	Describe("AddVlanFirewall", func() {
		Context("AddVlanFirewall given vlan id and with HA=false", func() {
			BeforeEach(func() {
				filenames := []string{
					"getItems_dedicatedFirewallNonHA",
					"placeOrder_firewallNonHA",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				fwManager = managers.NewFirewallManager(fakeSLSession)
			})
			It("Return order receipt and no error", func() {
				orderReceipt, err := fwManager.AddVlanFirewall(1455489, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(orderReceipt.OrderId).NotTo(Equal(nil))
				Expect(*orderReceipt.PlacedOrder.Items[0].CategoryCode).To(Equal("vlan_firewall"))
				Expect(*orderReceipt.PlacedOrder.Items[0].Description).To(Equal("Hardware Firewall (Dedicated)"))
			})
		})
		Context("AddVlanFirewall given vlan id and with HA=true", func() {
			BeforeEach(func() {
				filenames := []string{
					"getItems_dedicatedFirewallHA",
					"placeOrder_firewallHA",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				fwManager = managers.NewFirewallManager(fakeSLSession)
			})
			It("Return order receipt and no error", func() {
				orderReceipt, err := fwManager.AddVlanFirewall(1455489, true)
				Expect(err).ToNot(HaveOccurred())
				Expect(orderReceipt.OrderId).NotTo(Equal(nil))
				Expect(*orderReceipt.PlacedOrder.Items[0].CategoryCode).To(Equal("vlan_firewall"))
				Expect(*orderReceipt.PlacedOrder.Items[0].Description).To(Equal("Hardware Firewall (High Availability)"))
			})
		})
	})

	Describe("AddStandardFirewall", func() {
		Context("AddStandardFirewall given server id and isVirtual=true", func() {
			BeforeEach(func() {
				filenames := []string{
					"getItems_100MFirewall",
					"placeOrder_VSFirewall",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				fwManager = managers.NewFirewallManager(fakeSLSession)
			})
			It("Return order receipt and no error", func() {
				orderReceipt, err := fwManager.AddStandardFirewall(25868261, true)
				Expect(err).ToNot(HaveOccurred())
				Expect(orderReceipt.OrderId).NotTo(Equal(nil))
				Expect(len(orderReceipt.OrderDetails.VirtualGuests) > 0).To(BeTrue())
				Expect(*orderReceipt.PlacedOrder.Items[0].CategoryCode).To(Equal("firewall"))
			})
		})
		//TODO
		// Context("AddStandardFirewall given server id and isVirtual=false", func() {
		// 	BeforeEach(func() {
		// 		filenames := []string{
		// 			"getItems_1000MFirewall",
		// 			"SoftLayer_Product_Order_placeOrder_HWFirewall",
		// 		}
		// 		fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
		// 		fwManager = managers.NewFirewallManager(fakeSLSession)
		// 	})
		// 	It("Return order receipt and no error", func() {
		// 		orderReceipt, err := fwManager.AddStandardFirewall(25868261, false)
		// 		Expect(err).ToNot(HaveOccurred())
		// 		Expect(orderReceipt.OrderId).NotTo(Equal(nil))
		// 		Expect(len(orderReceipt.OrderDetails.Hardware) > 0).To(BeTrue())
		// 		Expect(*orderReceipt.PlacedOrder.Items[0].CategoryCode).To(Equal("firewall"))
		// 	})
		// })
	})

	Describe("GetFirewalls", func() {
		Context("GetFirewalls", func() {
			BeforeEach(func() {
				filenames := []string{"getNetworkVlans_firewall",}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				fwManager = managers.NewFirewallManager(fakeSLSession)
			})
			It("Return a list of vlans with firewalls and no error", func() {
				fws, err := fwManager.GetFirewalls()
				Expect(err).ToNot(HaveOccurred())
				Expect(len(fws) > 0).To(BeTrue())
			})
		})
	})

	Describe("GetFirewallBillingItem", func() {
		Context("GetFirewallBillingItem given firewall id and dedicated=true", func() {
			It("Return the billing item and no error", func() {
				billingItem, err := fwManager.GetFirewallBillingItem(1455551, true)
				Expect(err).ToNot(HaveOccurred())
				Expect(*billingItem.CategoryCode).To(Equal("network_vlan"))
			})
		})
		Context("GetFirewallBillingItem given firewall id and dedicated=false", func() {
			It("Return the billing item and no error", func() {
				billingItem, err := fwManager.GetFirewallBillingItem(1455551, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(*billingItem.CategoryCode).To(Equal("network_vlan"))
			})
		})
	})

	//TODO
	// Describe("GetStandardFirewallRules", func() {
	// 	Context("GetStandardFirewallRules given firewall id", func() {
	// 		It("Return a list of firewall rules and no error", func() {
	// 		})
	// 	})
	// })

	// Describe("GetDedicatedFirewallRules", func() {
	// 	Context("GetDedicatedFirewallRules given firewall id", func() {
	// 		It("Return a list of firewall rules and no error", func() {

	// 		})
	// 	})
	// })

	Describe("CancelFirewall", func() {
		Context("CancelFirewall given firewall id and dedicated=true", func() {
			It("Return no error", func() {
				err := fwManager.CancelFirewall(123, true)
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Context("CancelFirewall given firewall id and dedicated=false", func() {
			It("Return no error", func() {
				err := fwManager.CancelFirewall(123, false)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("GetStandardPackage", func() {
		Context("GetStandardPackage given server id and virtual=true", func() {
			BeforeEach(func() {
				filenames := []string{"getItems_100MFirewall",}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				fwManager = managers.NewFirewallManager(fakeSLSession)
			})
			It("Return a product item and no error", func() {
				items, err := fwManager.GetStandardPackage(25804753, true)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(items)).To(Equal(1))
				Expect(int(*items[0].Capacity)).To(Equal(100))
			})
		})
		Context("GetStandardPackage given server id and virtual=false", func() {
			BeforeEach(func() {
				filenames := []string{"getItems_1000MFirewall",}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				fwManager = managers.NewFirewallManager(fakeSLSession)
			})
			It("Return a product item and no error", func() {
				items, err := fwManager.GetStandardPackage(25804753, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(items)).To(Equal(1))
				Expect(int(*items[0].Capacity)).To(Equal(1000))
			})
		})
	})

	Describe("GetDedicatedPackage", func() {
		Context("GetDedicatedPackage with HA=false", func() {
			BeforeEach(func() {
				filenames := []string{"getItems_dedicatedFirewallNonHA",}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				fwManager = managers.NewFirewallManager(fakeSLSession)
			})
			It("Return a product item and no error", func() {
				items, err := fwManager.GetDedicatedPackage(false)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(items)).To(Equal(1))
				Expect(*items[0].KeyName).To(Equal("HARDWARE_FIREWALL_DEDICATED"))
				Expect(*items[0].ItemCategory.CategoryCode).To(Equal("vlan_firewall"))
				found := false
				for _, price := range items[0].Prices {
					if price.LocationGroupId == nil {
						found = true
						break
					}
				}
				//need to find a price whose locationGroupId is nil
				Expect(found).To(BeTrue())
			})
		})
		Context("GetDedicatedPackage with HA=true", func() {
			BeforeEach(func() {
				filenames := []string{"getItems_dedicatedFirewallHA",}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				fwManager = managers.NewFirewallManager(fakeSLSession)
			})
			It("Return a product item and no error", func() {
				items, err := fwManager.GetDedicatedPackage(true)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(items)).To(Equal(1))
				Expect(*items[0].KeyName).To(Equal("HARDWARE_FIREWALL_HIGH_AVAILABILITY"))
				Expect(*items[0].ItemCategory.CategoryCode).To(Equal("vlan_firewall"))
				found := false
				for _, price := range items[0].Prices {
					if price.LocationGroupId == nil {
						found = true
						break
					}
				}
				//need to find a price whose locationGroupId is nil
				Expect(found).To(BeTrue())
			})
		})
	})

	Describe("GetFirewallPortSpeed", func() {
		Context("GetFirewallPortSpeed given the server id and isVirtual=true", func() {
			It("Return port spped and no error", func() {
				portSpeed, err := fwManager.GetFirewallPortSpeed(25804753, true)
				Expect(err).ToNot(HaveOccurred())
				Expect(portSpeed).To(Equal(100))
			})
		})
		Context("GetFirewallPortSpeed given the server id and isVirtual=false", func() {
			It("Return port spped and no error", func() {
				portSpeed, err := fwManager.GetFirewallPortSpeed(25804753, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(portSpeed).To(Equal(1000))
			})
		})
	})
})

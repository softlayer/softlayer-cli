package managers_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("NetworkManager", func() {
	var (
		fakeSLSession  *session.Session
		networkManager managers.NetworkManager
	)
	BeforeEach(func() {
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		networkManager = managers.NewNetworkManager(fakeSLSession)
	})

	Describe("Get detail about a vlan", func() {
		Context("Get detail about a vlan given its ID", func() {
			It("It returns the details of the vlan", func() {
				vlan, err := networkManager.GetVlan(1262125, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(*vlan.Id).To(Equal(1262125))
				Expect(*vlan.VlanNumber).To(Equal(956))
				Expect(*vlan.PrimaryRouter.DatacenterName).To(Equal("Dallas 9"))
				Expect(len(vlan.Subnets)).To(Equal(2))
				Expect(len(vlan.Hardware)).To(Equal(1))
				Expect(len(vlan.VirtualGuests)).To(Equal(0))
			})
		})
	})

	Describe("Get all vlans", func() {
		BeforeEach(func() {
			fakeSLSession = testhelpers.NewFakeSoftlayerPagnationSession(nil)
			networkManager = managers.NewNetworkManager(fakeSLSession)
		})
		Context("Get all vlans under current account", func() {
			It("It returns a list of vlans", func() {
				vlans, err := networkManager.ListVlans("", 0, "", 0, "")
				Expect(err).ToNot(HaveOccurred())
				for _, vlan := range vlans {
					Expect(*vlan.Id).ShouldNot(BeNil())
					//Expect(*vlan.Name).ShouldNot(BeNil()) //some vlan has no name
					Expect(*vlan.VlanNumber).ShouldNot(BeNil())
					Expect(*vlan.NetworkSpace).Should(Or(Equal("PUBLIC"), Equal("PRIVATE")))
					Expect(*vlan.PrimaryRouter.Datacenter.Name).ShouldNot(BeNil())
				}
			})
		})
	})

	Describe("Get all subnets", func() {
		BeforeEach(func() {
			fakeSLSession = testhelpers.NewFakeSoftlayerPagnationSession(nil)
			networkManager = managers.NewNetworkManager(fakeSLSession)
		})
		Context("Get all subnets under current account", func() {
			It("It returns a list of subnets", func() {
				subnets, err := networkManager.ListSubnets("", "", 0, "", "", 0, "")
				Expect(err).ToNot(HaveOccurred())
				for _, subnet := range subnets {
					Expect(*subnet.Id).ShouldNot(BeNil())
					Expect(*subnet.NetworkIdentifier).ShouldNot(BeNil())
					Expect(*subnet.SubnetType).Should(Or(Equal("PRIMARY_6"), Equal("STATIC_IP_ROUTED"), Equal("PRIMARY"), Equal("ADDITIONAL_PRIMARY"), Equal("SECONDARY_ON_VLAN"), Equal("STATIC_IP_ROUTED_6"), Equal("SUBNET_ON_VLAN")))
					if subnet.NetworkVlan != nil {
						Expect(*subnet.NetworkVlan.NetworkSpace).Should(Or(Equal("PUBLIC"), Equal("PRIVATE")))
					} else {
						//fmt.Println("no vlan", *subnet.Id) //id=510674 's vlan is empty
					}
					Expect(*subnet.Datacenter.Name).ShouldNot(BeNil())
					Expect(len(subnet.IpAddresses) > 0).Should(BeTrue())
				}
			})
		})
	})

	Describe("Get detail about a subnet", func() {
		Context("Get detail about a subnet given its ID", func() {
			It("It returns the details of the vlan", func() {
				subnet, err := networkManager.GetSubnet(510674, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(*subnet.Id).To(Equal(510674))
				Expect(*subnet.NetworkVlanId).To(Equal(193419))
				Expect(*subnet.Version).To(Equal(4))
				Expect(*subnet.NetworkIdentifier).To(Equal("10.40.92.68"))
			})
		})
	})

	Describe("Get all global ips", func() {
		Context("Get all global ips under current account", func() {
			It("It returns a list of global ips ", func() {
				globalips, err := networkManager.ListGlobalIPs(0, 0)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(globalips) > 0).Should(BeTrue())
				for _, ip := range globalips {
					Expect(*ip.Id).ShouldNot(BeNil())
					Expect(*ip.IpAddress.IpAddress).To(Equal(*ip.IpAddress.Subnet.NetworkIdentifier))
					Expect(*ip.IpAddress.Subnet.Version).Should(Or(Equal(4), Equal(6)))
				}
			})
		})
		Context("Get all v4 global ips under current account", func() {
			It("It returns a list of global ips ", func() {
				globalips, err := networkManager.ListGlobalIPs(4, 0)
				Expect(err).ToNot(HaveOccurred())
				for _, ip := range globalips {
					Expect(*ip.Id).ShouldNot(BeNil())
					Expect(*ip.IpAddress.IpAddress).To(Equal(*ip.IpAddress.Subnet.NetworkIdentifier))
					Expect(*ip.IpAddress.Subnet.Version).Should(Equal(4))
				}
			})
		})
		Context("Get all v6 global ips under current account", func() {
			It("It returns a list of global ips ", func() {
				globalips, err := networkManager.ListGlobalIPs(6, 0)
				Expect(err).ToNot(HaveOccurred())
				for _, ip := range globalips {
					Expect(*ip.Id).ShouldNot(BeNil())
					Expect(*ip.Id).ShouldNot(BeNil())
					Expect(*ip.IpAddress.IpAddress).To(Equal(*ip.IpAddress.Subnet.NetworkIdentifier))
					Expect(*ip.IpAddress.Subnet.Version).Should(Equal(6))
				}
			})
		})
	})

	Describe("Route a global ip to a target route", func() {
		Context("Route a global ip to a target route", func() {
			It("It returns no error", func() {
				_, err := networkManager.AssignGlobalIP(510674, "")
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Unroute a global ip from a target route", func() {
		Context("Unroute a global ip from a target route", func() {
			It("It returns no error", func() {
				_, err := networkManager.UnassignGlobalIP(510674)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Get detail about ipaddress", func() {
		Context("Get detail about a ipaddress given its IP address", func() {
			It("It returns the details of the vlan", func() {
				ip, err := networkManager.IPLookup("10.40.207.172")
				Expect(err).ToNot(HaveOccurred())
				Expect(*ip.Id).To(Equal(10776596))
				Expect(*ip.SubnetId).To(Equal(514990))
				Expect(*ip.IpAddress).To(Equal("10.40.207.172"))
			})
		})
	})

	Describe("Cancel a global ip", func() {
		Context("Cancel a global ip given its ID", func() {
			It("It returns no error", func() {
				err := networkManager.CancelGlobalIP(510674)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Cancel a subnet", func() {
		Context("Cancel a subnet given its ID", func() {
			It("It returns no error", func() {
				err := networkManager.CancelSubnet(510674)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Add a subnet", func() {
		Context("Add a subnet", func() {
			It("It returns an order", func() {
				order, err := networkManager.AddSubnet("public", 4, 1307, 4, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(*order.OrderId).ShouldNot(BeNil())
				Expect(*order.OrderDate).ShouldNot(BeNil())
			})
		})
		Context("Add a subnet-test order", func() {
			It("It returns no error", func() {
				order, err := networkManager.AddSubnet("public", 4, 1307, 4, true)
				Expect(err).ToNot(HaveOccurred())
				Expect(order.OrderId).Should(BeNil())
			})
		})
	})

	Describe("Add a globalIP", func() {
		Context("Add a globalIP", func() {
			It("It returns an order", func() {
				order, err := networkManager.AddGlobalIP(4, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(*order.OrderId).ShouldNot(BeNil())
				Expect(*order.OrderDate).ShouldNot(BeNil())
			})
		})
		Context("Add a globalIP order", func() {
			It("It returns no error", func() {
				order, err := networkManager.AddGlobalIP(6, true)
				Expect(err).ToNot(HaveOccurred())
				Expect(order.OrderId).Should(BeNil())
			})
		})
	})

	Describe("Edit Vlan", func() {
		Context("Edit vlan's name", func() {
			It("It returns no error ", func() {
				err := networkManager.EditVlan(0, "vlan-rename")
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Cancel Vlan", func() {
		Context("Cancel vlan by vlan ID", func() {
			It("It returns no error ", func() {
				err := networkManager.CancelVLAN(0)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("List all data centers", func() {
		Context("List all data centers", func() {
			It("It returns no error ", func() {
				datacenters, err := networkManager.ListDatacenters()
				Expect(err).ToNot(HaveOccurred())
				for key, value := range datacenters {
					Expect(key > 0).Should(BeTrue())
					Expect(value != "").Should(BeTrue())
				}
			})
		})
	})

	Describe("List all routers in a datacenter", func() {
		Context("List all routers in a datacenter", func() {
			It("It returns no error ", func() {
				routers, err := networkManager.ListRouters(123, "")
				Expect(err).ToNot(HaveOccurred())
				for _, r := range routers {
					Expect(r != "").Should(BeTrue())
					//verify the name of router can be split to hostname and datacenter short name
					result := strings.Split(r, ".")
					Expect(result[0] != "").Should(BeTrue())
					Expect(result[1] != "").Should(BeTrue())
					//verify the router hostname either start with bcr(private vlan type) or fcr(public vlan type)
					Expect(strings.HasPrefix(result[0], "bcr") || strings.HasPrefix(result[0], "fcr")).Should(BeTrue())
				}
			})
		})
	})
})

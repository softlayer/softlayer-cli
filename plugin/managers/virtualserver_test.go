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

var _ = Describe("VirtualServerManager", func() {
	var (
		fakeSLSession *session.Session
		vsManager     managers.VirtualServerManager
	)

	BeforeEach(func() {
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		vsManager = managers.NewVirtualServerManager(fakeSLSession)
	})

	Describe("Cancel instance", func() {
		Context("Cancel instance given its ID", func() {
			It("It returns no error", func() {
				err := vsManager.CancelInstance(1234567)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Migrate instance", func() {
		Context("Migrate instance given its ID", func() {
			It("It returns Provisioning Version1 Transaction instance", func() {
				vsMigrate, err := vsManager.MigrateInstance(1234567)
				Expect(err).ToNot(HaveOccurred())
				Expect(err).ToNot(HaveOccurred())
				Expect(*vsMigrate.Id).Should(Equal(1234567))
			})
		})
	})

	Describe("Create instance", func() {
		Context("Create instance give a template instance", func() {
			It("It returns created virtual guest instance", func() {
				template, _ := vsManager.GetInstance(25804753, "")
				vs, err := vsManager.CreateInstance(&template)
				Expect(err).ToNot(HaveOccurred())
				Expect(err).ToNot(HaveOccurred())
				Expect(*vs.Hostname).Should(Equal("wilma2"))
				Expect(*vs.Domain).Should(Equal("wilma.org"))
				Expect(*vs.Datacenter.Name).Should(Equal("par01"))
				Expect(*vs.StartCpus).Should(Equal(2))
				Expect(*vs.MaxMemory).Should(Equal(2048))
				Expect(*vs.HourlyBillingFlag).Should(BeTrue())
				Expect(*vs.DedicatedAccountHostOnlyFlag).Should(BeFalse())
				Expect(*vs.PrivateNetworkOnlyFlag).Should(BeFalse())
				Expect(*vs.PostInstallScriptUri).Should(Equal("http://www.mycompany/scipt1"))
				Expect(*vs.OperatingSystemReferenceCode).Should(Equal("CENTOS_7_64"))
				Expect(*vs.LocalDiskFlag).Should(BeTrue())
				Expect(len(vs.NetworkComponents) > 0).Should(BeTrue())
				for _, network := range vs.NetworkComponents {
					Expect(*network.MaxSpeed).Should(Equal(10))
				}
				Expect(*vs.PrimaryNetworkComponent.NetworkVlan.Id).Should(Equal(1421723))
				Expect(*vs.PrimaryBackendNetworkComponent.NetworkVlan.Id).Should(Equal(1421725))
			})
		})
	})

	Describe("Generate creation template", func() {
		Context("Generate creation template give a parameter map", func() {
			vs := new(datatypes.Virtual_Guest)
			var err error
			params := make(map[string]interface{})
			params["hostname"] = "wilma2"
			params["domain"] = "wilma.org"
			params["cpu"] = 2
			params["memory"] = 2048
			params["datacenter"] = "par01"
			params["os"] = "CENTOS_7_64"
			params["billing"] = true
			params["dedicated"] = false
			params["private"] = false
			params["san"] = false
			params["i"] = "http://www.mycompany/scipt1"
			params["disks"] = []int{25}
			params["network"] = 10
			params["vlan-public"] = 1421723
			params["vlan-private"] = 1421725
			It("It returns virtual guest template", func() {
				vs, err = vsManager.GenerateInstanceCreationTemplate(vs, params)
				Expect(err).ToNot(HaveOccurred())
				Expect(*vs.Hostname).Should(Equal("wilma2"))
				Expect(*vs.Domain).Should(Equal("wilma.org"))
				Expect(*vs.Datacenter.Name).Should(Equal("par01"))
				Expect(*vs.StartCpus).Should(Equal(2))
				Expect(*vs.MaxMemory).Should(Equal(2048))
				Expect(*vs.HourlyBillingFlag).Should(BeTrue())
				Expect(*vs.DedicatedAccountHostOnlyFlag).Should(BeFalse())
				Expect(*vs.PrivateNetworkOnlyFlag).Should(BeFalse())
				Expect(*vs.PostInstallScriptUri).Should(Equal("http://www.mycompany/scipt1"))
				Expect(*vs.OperatingSystemReferenceCode).Should(Equal("CENTOS_7_64"))
				Expect(*vs.LocalDiskFlag).Should(BeTrue())
				Expect(len(vs.NetworkComponents) > 0).Should(BeTrue())
				for _, network := range vs.NetworkComponents {
					Expect(*network.MaxSpeed).Should(Equal(10))
				}
				Expect(*vs.PrimaryNetworkComponent.NetworkVlan.Id).Should(Equal(1421723))
				Expect(*vs.PrimaryBackendNetworkComponent.NetworkVlan.Id).Should(Equal(1421725))
			})
		})
	})

	Describe("Verify instance creation", func() {
		Context("Verify instance creation given a template virtual guest", func() {
			It("It returns order", func() {
				template, _ := vsManager.GetInstance(25804753, "")
				order, err := vsManager.VerifyInstanceCreation(template)
				Expect(err).ToNot(HaveOccurred())
				Expect(*order.ComplexType).To(Equal("SoftLayer_Container_Product_Order_Virtual_Guest"))
				Expect(*order.Quantity).To(Equal(1))
				Expect(*order.UseHourlyPricing).To(Equal(true))
			})
		})
	})

	Describe("Get instance", func() {
		Context("get instance given its ID", func() {
			It("It return the virtual guest instance", func() {
				vs, err := vsManager.GetInstance(25804753, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(*vs.Id).To(Equal(25804753))
			})
		})
	})

	Describe("Get liked instance", func() {
		Context("Get liked instance given template vs and liked vs ID", func() {
			It("It returns an instance that has the same configuration settings", func() {
				template := datatypes.Virtual_Guest{}
				vs, err := vsManager.GetLikedInstance(&template, 25804753)
				Expect(err).ToNot(HaveOccurred())
				Expect(vs.Id).Should(BeNil())
				Expect(*vs.Hostname).Should(Equal("wilma2"))
				Expect(*vs.Domain).Should(Equal("wilma.org"))
				Expect(*vs.Datacenter.Name).Should(Equal("par01"))
				Expect(*vs.StartCpus).Should(Equal(2))
				Expect(*vs.MaxMemory).Should(Equal(2048))
				Expect(*vs.HourlyBillingFlag).Should(BeTrue())
				Expect(*vs.DedicatedAccountHostOnlyFlag).Should(BeFalse())
				Expect(*vs.PrivateNetworkOnlyFlag).Should(BeFalse())
				Expect(*vs.PostInstallScriptUri).ShouldNot(BeNil())
				Expect(len(vs.UserData) > 0).Should(BeTrue())
				Expect(len(vs.NetworkComponents) > 0).Should(BeTrue())
				Expect(*vs.OperatingSystemReferenceCode).Should(Equal("CENTOS_7_64"))
				Expect(*vs.LocalDiskFlag).Should(BeTrue())
			})
		})
	})

	Describe("Capture instance to an image", func() {
		Context("Capture instance to an image", func() {
			It("It returns no err and a transaction", func() {
				txn, err := vsManager.CaptureImage(25804753, "wilmaimage", "imagenote", true)
				Expect(err).ToNot(HaveOccurred())
				Expect(*txn.GuestId).To(Equal(25804753))
			})
		})
	})

	Describe("Pause instance", func() {
		Context("Pause instance given its ID", func() {
			It("It returns no error", func() {
				err := vsManager.PauseInstance(123456)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Power on instance", func() {
		Context("Power on instance given its ID", func() {
			It("It returns no error", func() {
				err := vsManager.PowerOnInstance(123456)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Power off instance", func() {
		Context("Power off instance given its ID", func() {
			It("It returns no error", func() {
				err := vsManager.PowerOffInstance(123456, false, false)
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Context("Power off softly instance given its ID", func() {
			It("It returns no error", func() {
				err := vsManager.PowerOffInstance(123456, true, false)
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Context("Power off hardly instance given its ID", func() {
			It("It returns no error", func() {
				err := vsManager.PowerOffInstance(123456, false, true)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Reboot instance", func() {
		Context("Reboot instance given its ID", func() {
			It("It returns no error", func() {
				err := vsManager.RebootInstance(123456, false, false)
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Context("Reboot softly instance given its ID", func() {
			It("It returns no error", func() {
				err := vsManager.RebootInstance(123456, true, false)
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Context("Reboot hardly instance given its ID", func() {
			It("It returns no error", func() {
				err := vsManager.RebootInstance(123456, false, true)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Reload instance", func() {
		Context("Reload instance given its ID", func() {
			It("It returns no error", func() {
				err := vsManager.ReloadInstance(123456, "http://www.mycompany/scripts/12345.sh", []int{123}, 234567)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Resume instance", func() {
		Context("Resume instance given its ID", func() {
			It("It returns no error", func() {
				err := vsManager.ResumeInstance(123456)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Rescue instance", func() {
		Context("Rescue instance given its ID", func() {
			It("It returns no error", func() {
				err := vsManager.RescueInstance(123456)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Instance is ready", func() {
		Context("Check the instance if it is ready for use", func() {
			It("It returns it is ready", func() {
				ready, msg, err := vsManager.InstanceIsReady(123456, time.Now())
				Expect(err).ToNot(HaveOccurred())
				Expect(ready).To(Equal(true))
				Expect(msg).To(Equal(""))
			})
		})
	})

	Describe("Set user metadata for instance", func() {
		Context("Set user metadata for instance given its ID and a string slice", func() {
			It("It returns no error", func() {
				err := vsManager.SetUserMetadata(123456, []string{"mydata"})
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Set tags for instance", func() {
		Context("Set user metadata for instance given its ID and a string of tags", func() {
			It("It returns no error", func() {
				err := vsManager.SetTags(123456, "mytags")
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Set network speed for instance", func() {
		Context("Set public network speed for instance  given its ID", func() {
			It("It returns no error", func() {
				err := vsManager.SetNetworkPortSpeed(123456, true, 1000)
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Context("Set private network speed for instance  given its ID", func() {
			It("It returns no error", func() {
				err := vsManager.SetNetworkPortSpeed(123456, false, 1000)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Edit instance", func() {
		Context("Edit instance to update its hostname, domain, userdata, tag and network speed", func() {
			It("It returns a list of bool and messages", func() {
				speed := 1000
				succeesses, msgs := vsManager.EditInstance(123456, "wilma", "mycompany.com", "mydata", "mytag", &speed, &speed)
				for index, success := range succeesses {
					Expect(success).To(Equal(true))
					Expect(msgs[index]).ShouldNot(BeNil())
				}
			})
		})
	})

	Describe("List virtual guest instance under current acount", func() {
		BeforeEach(func() {
			fakeSLSession = testhelpers.NewFakeSoftlayerPagnationSession(nil)
			vsManager = managers.NewVirtualServerManager(fakeSLSession)
		})
		Context("List all virtual guest instance under current acount", func() {
			It("It returns a list of virtual guest instances", func() {
				vss, err := vsManager.ListInstances(false, false, "", "", "", "", "", "", 0, 0, 0, 0, nil, "")
				Expect(err).ToNot(HaveOccurred())
				for _, vs := range vss {
					Expect(*vs.Account.Id).To(Equal(278444))
					Expect(*vs.Id).ShouldNot(BeNil())
				}
			})
		})
		Context("List hourly-billed virtual guest instance under current acount", func() {
			It("It returns a list of hourly-billed virtual guest instances", func() {
				vss, err := vsManager.ListInstances(true, false, "", "", "", "", "", "", 0, 0, 0, 0, nil, "")
				Expect(err).ToNot(HaveOccurred())
				for _, vs := range vss {
					Expect(*vs.Account.Id).To(Equal(278444))
					Expect(*vs.HourlyBillingFlag).To(Equal(true))
				}
			})
		})
		Context("List monthly-billed virtual guest instance under current acount", func() {
			It("It returns a list of monthly-billed virtual guest instances", func() {
				vss, err := vsManager.ListInstances(false, true, "", "", "", "", "", "", 0, 0, 0, 0, nil, "")
				Expect(err).ToNot(HaveOccurred())
				for _, vs := range vss {
					Expect(*vs.Account.Id).To(Equal(278444))
					Expect(*vs.HourlyBillingFlag).To(Equal(false))
				}
			})
		})
	})

})

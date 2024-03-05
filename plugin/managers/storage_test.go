package managers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("StorageManager", func() {
	var (
		fakeSLSession  *session.Session
		fakeHandler    *testhelpers.FakeTransportHandler
		StorageManager managers.StorageManager
	)
	BeforeEach(func() {
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		fakeHandler = testhelpers.GetSessionHandler(fakeSLSession)
		StorageManager = managers.NewStorageManager(fakeSLSession)
	})
	AfterEach(func() {
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})

	Describe("GetBlockVolumeAccessList", func() {
		Context("GetBlockVolumeAccessList given a volume id", func() {
			It("Return the volume with allowed access hosts", func() {
				volume, err := StorageManager.GetVolumeAccessList(17336531)
				Expect(err).ToNot(HaveOccurred())
				for _, vs := range volume.AllowedVirtualGuests {
					Expect(vs.Id).NotTo(Equal(nil))
					Expect(vs.FullyQualifiedDomainName).NotTo(Equal(nil))
					Expect(vs.AllowedHost.Id).NotTo(Equal(nil))
					Expect(vs.AllowedHost.Name).NotTo(Equal(nil))
					Expect(vs.AllowedHost.Credential.Username).NotTo(Equal(nil))
					Expect(vs.AllowedHost.Credential.Password).NotTo(Equal(nil))
				}
			})
		})
	})

	Describe("AuthorizeHostToVolume", func() {
		Context("AuthorizeHostToVolume given a volume id and a list of hosts", func() {
			It("Return no error", func() {
				hosts, err := StorageManager.AuthorizeHostToVolume(17336531, nil, []int{25868261}, nil, nil)
				Expect(err).ToNot(HaveOccurred())
				for _, host := range hosts {
					//bug current the returned list is empty
					Expect(host.Id).NotTo(Equal(nil))
				}
			})
		})
	})

	Describe("DeauthorizeHostToVolume", func() {
		Context("DeauthorizeHostToVolume given a volume id and a list of hosts", func() {
			It("Return no error", func() {
				hosts, err := StorageManager.DeauthorizeHostToVolume(17336531, nil, []int{25868261}, nil, nil)
				Expect(err).ToNot(HaveOccurred())
				for _, host := range hosts {
					//TODO current the returned list is empty
					Expect(host.Id).NotTo(Equal(nil))
				}
			})
		})
	})

	Describe("OrderReplicantVolume", func() {
		Context("OrderReplicantVolume given volume id and parameters", func() {
			BeforeEach(func() {
				filenames := []string{
					"getAllObjects_saas",
					"placeOrder_endurance",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				StorageManager = managers.NewStorageManager(fakeSLSession)
			})
			It("Return no error", func() {
				order, err := StorageManager.OrderReplicantVolume("block", 17336531, "DAILY", "tok02", 4, 0, "LINUX")
				Expect(err).ToNot(HaveOccurred())
				Expect(order.OrderId).NotTo(Equal(nil))
			})
		})
	})

	Describe("FaileOverToReplicant", func() {
		Context("FaileOverToReplicant given volume id and replicant id", func() {
			It("Return no error", func() {
				err := StorageManager.FailOverToReplicant(17336531, 999)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("FaileBackFromReplicant", func() {
		Context("FaileBackFromReplicant given volume id", func() {
			It("Return no error", func() {
				err := StorageManager.FailBackFromReplicant(17336531)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("ListBlockVolumes", func() {
		Context("ListBlockVolumes under current account", func() {
			It("Block Happy Path", func() {
				volumes, err := StorageManager.ListVolumes("block", "", "", "", "", 0, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(len(volumes)).Should(BeNumerically(">", 0))
				for _, volume := range volumes {
					Expect(volume.Id).NotTo(Equal(nil))
					Expect(*volume.StorageType.KeyName).To(Equal("ENDURANCE_BLOCK_STORAGE"))
				}
				apiCalls := fakeHandler.ApiCallLogs
				Expect(len(apiCalls)).To(Equal(1))
				Expect(apiCalls[0]).To(MatchFields(IgnoreExtras, Fields{
					"Service": Equal("SoftLayer_Account"),
					"Method":  Equal("getIscsiNetworkStorage"),
					"Options": PointTo(MatchFields(IgnoreExtras, Fields{"Limit": PointTo(Equal(50))})),
				}))
			})
			It("File Happy Path", func() {
				volumes, err := StorageManager.ListVolumes("file", "", "", "", "", 0, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(len(volumes)).Should(BeNumerically(">", 0))
				for _, volume := range volumes {
					Expect(volume.Id).NotTo(Equal(nil))
					Expect(*volume.StorageType.KeyName).To(Equal("ENDURANCE_FILE_STORAGE"))
				}
				apiCalls := fakeHandler.ApiCallLogs
				Expect(len(apiCalls)).To(Equal(1))
				Expect(apiCalls[0]).To(MatchFields(IgnoreExtras, Fields{
					"Service": Equal("SoftLayer_Account"),
					"Method":  Equal("getNasNetworkStorage"),
					"Options": PointTo(MatchFields(IgnoreExtras, Fields{"Limit": PointTo(Equal(50))})),
				}))
			})
		})
		Context("Issue822 - Special case for ListVolumes with a datacenter filter", func() {
			It("Block: No Result Limit", func() {
				_, err := StorageManager.ListVolumes("block", "dal10", "", "", "", 0, "")
				Expect(err).ToNot(HaveOccurred())
				apiCalls := fakeHandler.ApiCallLogs
				Expect(len(apiCalls)).To(Equal(1))
				// See https://pkg.go.dev/github.com/onsi/gomega/gstruct for this stuff
				// fmt.Printf("APICALL: %+v", apiCalls[0].Options)
				Expect(apiCalls[0]).To(MatchFields(IgnoreExtras, Fields{
					"Service": Equal("SoftLayer_Account"),
					"Method":  Equal("getIscsiNetworkStorage"),
					"Options": PointTo(MatchFields(IgnoreExtras, Fields{"Limit": BeNil()})),
				}))
			})
			It("File: No Result Limit", func() {
				_, err := StorageManager.ListVolumes("file", "dal10", "", "", "", 0, "")
				Expect(err).ToNot(HaveOccurred())
				apiCalls := fakeHandler.ApiCallLogs
				Expect(len(apiCalls)).To(Equal(1))
				// See https://pkg.go.dev/github.com/onsi/gomega/gstruct for this stuff
				// fmt.Printf("APICALL: %+v", apiCalls[0].Options)
				Expect(apiCalls[0]).To(MatchFields(IgnoreExtras, Fields{
					"Service": Equal("SoftLayer_Account"),
					"Method":  Equal("getNasNetworkStorage"),
					"Options": PointTo(MatchFields(IgnoreExtras, Fields{"Limit": BeNil()})),
				}))
			})
		})
	})

	Describe("GetBlockVolumeDetails", func() {
		Context("GetBlockVolumeDetails given volume id", func() {
			It("Return the volume details and no error", func() {
				volume, err := StorageManager.GetVolumeDetails("block", 17336531, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(volume.Id).NotTo(Equal(nil))
				Expect(*volume.CapacityGb).To(Equal(20))
				Expect(*volume.SnapshotCapacityGb).To(Equal("20"))
				Expect(*volume.Username).To(Equal("IBM01SEL278444-16"))
				Expect(*volume.StorageType.KeyName).To(Equal("ENDURANCE_BLOCK_STORAGE"))
				Expect(*volume.StorageTierLevel).To(Equal("WRITEHEAVY_TIER"))
				Expect(*volume.ReplicationStatus).To(Equal("FAILBACK_COMPLETED"))
			})
		})
	})

	Describe("OrderBlockVolume", func() {
		Context("OrderBlockVolume given type=endurance", func() {
			BeforeEach(func() {
				filenames := []string{
					"getAllObjects_saas",
					"placeOrder_endurance",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				StorageManager = managers.NewStorageManager(fakeSLSession)
			})
			It("Return the order receipt and no error", func() {
				orderReceipt, err := StorageManager.OrderVolume("block", "tok02", "endurance", "LINUX", 20, 4, 0, 20, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(orderReceipt.OrderId).NotTo(Equal(nil))
			})
		})
		Context("OrderBlockVolume given type=performance", func() {
			BeforeEach(func() {
				filenames := []string{
					"getAllObjects_saas",
					"placeOrder_performance",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				StorageManager = managers.NewStorageManager(fakeSLSession)
			})
			It("Return the order receipt and no error", func() {
				orderReceipt, err := StorageManager.OrderVolume("block", "tok02", "performance", "LINUX", 1000, 0, 1000, 500, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(orderReceipt.OrderId).NotTo(Equal(nil))
			})
		})
	})

	Describe("CancelBlockVolume", func() {
		Context("CancelBlockVolume given volume id", func() {
			It("Return no error", func() {
				err := StorageManager.CancelVolume("block", 17336531, "", true)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("GetBlockVolumeSnapshotList", func() {
		Context("GetBlockVolumeSnapshotList given volume id", func() {
			It("Return a list of snapshots and no error", func() {
				snapshots, err := StorageManager.GetVolumeSnapshotList(17336531)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(snapshots) > 0).To(Equal(true))
				for _, snapshot := range snapshots {
					Expect(snapshot.Id).NotTo(Equal(nil))
					Expect(snapshot.SnapshotCreationTimestamp).NotTo(Equal(nil))
					Expect(snapshot.SnapshotSizeBytes).NotTo(Equal(nil))
					Expect(*snapshot.StorageType.KeyName).To(Equal("SNAPSHOT"))
				}
			})
		})
	})

	Describe("DeleteSnapshot", func() {
		Context("DeleteSnapshot given snapshot id", func() {
			It("Return no error", func() {
				err := StorageManager.DeleteSnapshot(17360371)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("CreateSnapshot", func() {
		Context("CreateSnapshot given snapshot id", func() {
			It("Return the snapshot and no error", func() {
				snapshot, err := StorageManager.CreateSnapshot(17360371, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(snapshot.Id).NotTo(Equal(nil))
				Expect(snapshot.Username).NotTo(Equal(nil))
				Expect(*snapshot.CapacityGb).To(Equal(20))
				Expect(snapshot.SnapshotCreationTimestamp).NotTo(Equal(nil))
				Expect(*snapshot.NasType).To(Equal("SNAPSHOT"))
				Expect(snapshot.ServiceResourceName).NotTo(Equal(nil))
				Expect(snapshot.ServiceResourceBackendIpAddress).NotTo(Equal(nil))
			})
		})
	})

	Describe("CancelSnapshotSpace", func() {
		Context("CancelSnapshotSpace given volume id", func() {
			It("Return no error", func() {
				err := StorageManager.CancelSnapshotSpace("block", 17360371, "", true)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("EnableSnapshot", func() {
		Context("EnableSnapshot given volume id and scheduletype and time", func() {
			It("Return no error", func() {
				err := StorageManager.EnableSnapshot(17360371, "DAILY", 2, 0, 0, "SUNDAY")
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("DisableSnapshot", func() {
		Context("DisableSnapshot given volume id and scheduletype", func() {
			It("Return no error", func() {
				err := StorageManager.DisableSnapshots(17360371, "DAILY")
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("RestoreFromSnapshot", func() {
		Context("RestoreFromSnapshot given volume id and snapshot id", func() {
			It("Return no error", func() {
				err := StorageManager.RestoreFromSnapshot(17360371, 17360367)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("OrderSnapshotSpace", func() {
		Context("OrderSnapshotSpace given volume id and capacity and tier", func() {
			BeforeEach(func() {
				filenames := []string{
					"getAllObjects_saas",
					"placeOrder_endurance",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				StorageManager = managers.NewStorageManager(fakeSLSession)
			})
			It("Return order receipt and no error", func() {
				orderReceipt, err := StorageManager.OrderSnapshotSpace("block", 17360371, 20, 4, 0, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(orderReceipt.OrderId).NotTo(Equal(nil))
			})
		})
	})
	Describe("GetVolumeCountLimits", func() {
		Context("GetVolumeCountLimits test", func() {
			It("Return no error", func() {
				_, err := StorageManager.GetVolumeCountLimits()
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
	Describe("VolumeConvert", func() {
		Context("VolumeConvert test", func() {
			It("Return no error", func() {
				err := StorageManager.VolumeConvert(1234)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
	Describe("VolumeRefresh", func() {
		Context("VolumeRefresh test", func() {
			It("Return no error", func() {
				err := StorageManager.VolumeRefresh(1234, 4567, false)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
	Describe("Disaster Recovery Failover", func() {
		Context("DisasterRecoveryFailover test", func() {
			It("Return no error", func() {
				err := StorageManager.DisasterRecoveryFailover(1234, 4567)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})

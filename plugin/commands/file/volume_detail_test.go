package file_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/file"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Volume detail", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *file.VolumeDetailCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
		FakeStorageManager *testhelpers.FakeStorageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "file")
		cliCommand = file.NewVolumeDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("Volume detail", func() {
		Context("Volume detail without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("Volume detail with wrong volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
			})
		})

		Context("Volume detail with correct volume id but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.GetVolumeDetailsReturns(datatypes.Network_Storage{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get details of volume 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Volume detail with correct volume id", func() {
			BeforeEach(func() {
				FakeStorageManager.GetVolumeDetailsReturns(datatypes.Network_Storage{
					Id:               sl.Int(1234),
					Username:         sl.String("myvolume"),
					CapacityGb:       sl.Int(1000),
					Iops:             sl.String("400"),
					StorageTierLevel: sl.String("HEAVY_WRITE"),
					StorageType: &datatypes.Network_Storage_Type{
						KeyName: sl.String("performance"),
					},
					ServiceResource: &datatypes.Network_Service_Resource{
						Datacenter: &datatypes.Location{
							Name: sl.String("tok02"),
						},
					},
					ServiceResourceBackendIpAddress: sl.String("9.9.9.9"),
					SnapshotCapacityGb:              sl.String("500"),
					ParentVolume: &datatypes.Network_Storage{
						SnapshotSizeBytes: sl.String("10000"),
					},
					ActiveTransactionCount: sl.Uint(uint(1)),
					ActiveTransactions: []datatypes.Provisioning_Version1_Transaction{
						datatypes.Provisioning_Version1_Transaction{
							TransactionStatus: &datatypes.Provisioning_Version1_Transaction_Status{
								FriendlyName: sl.String("Restarting"),
							},
						},
					},
					ReplicationPartnerCount: sl.Uint(uint(1)),
					ReplicationStatus:       sl.String("replication finished"),
					ReplicationPartners: []datatypes.Network_Storage{
						datatypes.Network_Storage{
							Id:                              sl.Int(5678),
							Username:                        sl.String("myreplicant"),
							ServiceResourceBackendIpAddress: sl.String("9.9.9.8"),
							ServiceResource: &datatypes.Network_Service_Resource{
								Datacenter: &datatypes.Location{
									Name: sl.String("dal10"),
								},
							},
							ReplicationSchedule: &datatypes.Network_Storage_Schedule{
								Type: &datatypes.Network_Storage_Schedule_Type{
									Keyname: sl.String("DAILY"),
								},
							},
						},
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("myvolume"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("tok02"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("9.9.9.9"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Restarting"))
			})
		})
	})
})

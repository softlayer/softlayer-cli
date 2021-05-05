package file_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/file"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Volume detail", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *file.VolumeDetailCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = file.NewVolumeDetailCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        metadata.FileVolumeDetailMetaData().Name,
			Description: metadata.FileVolumeDetailMetaData().Description,
			Usage:       metadata.FileVolumeDetailMetaData().Usage,
			Flags:       metadata.FileVolumeDetailMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Volume detail", func() {
		Context("Volume detail without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Volume detail with wrong volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Volume ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Volume detail with correct volume id but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.GetVolumeDetailsReturns(datatypes.Network_Storage{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get details of volume 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
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
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"myvolume"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1000"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"400"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"HEAVY_WRITE"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"tok02"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"performance"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.9"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"500"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Restarting"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"replication finished"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"5678"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"myreplicant"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.8"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"dal10"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"DAILY"}))
			})
		})
	})
})

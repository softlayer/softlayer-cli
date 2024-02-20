package block_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("block object-storage-detail", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cliCommand         *block.ObjectStorageDetailCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewObjectStorageDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("block object-storage-detail", func() {
		Context("Return error", func() {
			It("Object storage detail without ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})

			It("Object storage detail with wrong ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Storage ID'. It must be a positive integer."))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format"))
			})
		})

		Context("Return error", func() {
			It("Failed get NAS Network Storages", func() {
				FakeStorageManager.GetNetworkStorageDetailReturns(datatypes.Network_Storage{}, errors.New("Failed to get details of storage"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get details of storage"))
			})

			It("Failed get NAS Network Storages", func() {
				FakeStorageManager.GetBucketsReturns([]datatypes.Container_Network_Storage_Hub_ObjectStorage_Bucket{}, errors.New("Failed to get bucket of storage"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get bucket of storage"))
			})
		})

		Context("Objet storage detail correct use", func() {
			BeforeEach(func() {
				FakeStorageManager.GetNetworkStorageDetailReturns(datatypes.Network_Storage{
					Id:       sl.Int(123),
					Username: sl.String("StorageName"),
					ServiceResource: &datatypes.Network_Service_Resource{
						Id:   sl.Int(123),
						Name: sl.String("ResourceName"),
						Type: &datatypes.Network_Service_Resource_Type{
							Type: sl.String("CLEVERSAFE_SVC_API"),
						},
						Datacenter: &datatypes.Location{
							Name: sl.String("sjc03"),
						},
					},
					StorageType: &datatypes.Network_Storage_Type{
						KeyName: sl.String("OBJECT_STORAGE_STANDARD"),
					},
				},
					nil)
				FakeStorageManager.GetBucketsReturns([]datatypes.Container_Network_Storage_Hub_ObjectStorage_Bucket{
					datatypes.Container_Network_Storage_Hub_ObjectStorage_Bucket{
						BytesUsed: sl.Int(6543211234),
						Name:      sl.String("BucketName"),
					},
				}, nil)
			})
			It("Return object storage detail", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("123"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("StorageName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ResourceName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("CLEVERSAFE_SVC_API"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("sjc03"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("OBJECT_STORAGE_STANDARD"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("6.09G"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("BucketName"))
			})

			It("Return object storage detail in json format", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--output=json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("["))
				Expect(fakeUI.Outputs()).To(ContainSubstring("]"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("{"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("}"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("StorageName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ResourceName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("CLEVERSAFE_SVC_API"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("sjc03"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("OBJECT_STORAGE_STANDARD"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("6.09G"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("BucketName"))
			})
		})
	})
})

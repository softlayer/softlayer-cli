package block_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("block object-storage-permission", func() {
	var (
		fakeUI                   *terminal.FakeUI
		FakeStorageManager       *testhelpers.FakeStorageManager
		FakeObjectStorageManager *testhelpers.FakeObjectStorageManager
		cliCommand               *block.ObjectStoragePermissionCommand
		fakeSession              *session.Session
		slCommand                *metadata.SoftlayerStorageCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		FakeObjectStorageManager = new(testhelpers.FakeObjectStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewObjectStoragePermissionCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
		cliCommand.ObjectStorageManager = FakeObjectStorageManager
	})

	Describe("block object-storage-permission", func() {
		Context("Return error", func() {
			It("Object storage permission without ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})

			It("Object storage permission with wrong ID", func() {
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
				FakeStorageManager.GetNetworkMessageDeliveryAccountsReturns(datatypes.Network_Storage_Hub_Cleversafe_Account{}, errors.New("Failed to get permissions."))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get permissions."))
			})

			It("Failed get NAS Network Storages", func() {
				FakeObjectStorageManager.GetEndpointsReturns([]datatypes.Container_Network_Storage_Hub_ObjectStorage_Endpoint{}, errors.New("Failed to get endPoints."))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get endPoints."))
			})
		})

		Context("Objet storage permission correct use", func() {
			BeforeEach(func() {
				FakeStorageManager.GetNetworkMessageDeliveryAccountsReturns(datatypes.Network_Storage_Hub_Cleversafe_Account{
					Uuid: sl.String("abc123"),
					Credentials: []datatypes.Network_Storage_Credential{
						datatypes.Network_Storage_Credential{
							Id:       sl.Int(12345),
							Username: sl.String("credential"),
							Password: sl.String("abc321"),
							Type: &datatypes.Network_Storage_Credential_Type{
								Description: sl.String("Description credential"),
							},
						},
					},
				}, nil)
				FakeObjectStorageManager.GetEndpointsReturns([]datatypes.Container_Network_Storage_Hub_ObjectStorage_Endpoint{
					datatypes.Container_Network_Storage_Hub_ObjectStorage_Endpoint{
						Region:   sl.String("us-geo"),
						Location: sl.String("Dallas"),
						Type:     sl.String("public"),
						Url:      sl.String("s3.us.cloud-object-storage"),
					},
				}, nil)
			})
			It("Return object storage permission", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("UUID"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("abc123"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Credentials"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Access Key ID"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("credential"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Secret Access Key"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("abc321"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Description"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Description credential"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("EndPoint URL´s"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Region"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("us-geo"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Location"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Dallas"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Type"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("public"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("URL"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("s3.us.cloud-object-storage"))
			})

			It("Return object storage permission in json format", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--output=json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("["))
				Expect(fakeUI.Outputs()).To(ContainSubstring("]"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("{"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("}"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("UUID"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("abc123"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Credentials"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Access Key ID"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("EndPoint URL´s"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Region"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("us-geo"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Location"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Dallas"))
			})
		})
	})
})

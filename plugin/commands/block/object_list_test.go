package block_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("block object-list", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.ObjectListCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewObjectListCommand(fakeUI, fakeStorageManager)
		cliCommand = cli.Command{
			Name:        block.BlockObjectListMetaData().Name,
			Description: block.BlockObjectListMetaData().Description,
			Usage:       block.BlockObjectListMetaData().Usage,
			Flags:       block.BlockObjectListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("block object-list", func() {

		Context("Return error", func() {

			It("Set invalid output", func() {
				err := testhelpers.RunCommand(cliCommand, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeStorageManager.GetHubNetworkStorageReturns([]datatypes.Network_Storage{}, errors.New("Failed to get Cloud Object Storages."))
			})
			It("Failed get Cloud Object Storages", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get Cloud Object Storages."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2017-11-08T00:00:00Z")
				fakerCloudObjectStorages := []datatypes.Network_Storage{
					datatypes.Network_Storage{
						Id:       sl.Int(123456),
						Username: sl.String("SLOSC123456-1"),
						StorageType: &datatypes.Network_Storage_Type{
							Description: sl.String("Object Storage Standard Account"),
							KeyName:     sl.String("OBJECT_STORAGE_STANDARD"),
						},
						BillingItem: &datatypes.Billing_Item{
							CreateDate: sl.Time(created),
						},
					},
				}
				fakeStorageManager.GetHubNetworkStorageReturns(fakerCloudObjectStorages, nil)
			})
			It("Return no error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SLOSC123456-1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Object Storage Standard Account"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("OBJECT_STORAGE_STANDARD"))
			})
		})
	})
})

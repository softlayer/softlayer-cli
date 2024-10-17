package file_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/file"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var orderReceipt = datatypes.Container_Product_Order_Receipt{
	OrderId: sl.Int(998877),
	PlacedOrder: &datatypes.Billing_Order{
		Items: []datatypes.Billing_Order_Item{
			datatypes.Billing_Order_Item{
				Description: sl.String("Test Item 1"),
			},
			datatypes.Billing_Order_Item{
				Description: sl.String("Another Test Item"),
			},
		},
	},
}
var _ = Describe("Volume Modify", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cliCommand         *file.VolumeModifyCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "file")
		cliCommand = file.NewVolumeModifyCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
		FakeStorageManager.GetVolumeIdReturns(1234, nil)
	})

	Describe("sl file volume-modify", func() {
		Context("No Id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("Bad Id", func() {
			It("return error", func() {
				FakeStorageManager.GetVolumeIdReturns(0, errors.New("BAD Volume ID"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "Abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("BAD Volume ID"))
			})
		})
		Context("Happy Path", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderModifiedVolumeReturns(orderReceipt, nil)
			})
			It("Success", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--new-size", "500", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 998877 was placed successfully!."))
			})
		})
	})
})

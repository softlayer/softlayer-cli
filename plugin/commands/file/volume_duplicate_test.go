package file_test

import (
	"errors"

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

var FakeOrderReceipt = datatypes.Container_Product_Order_Receipt{
	OrderId: sl.Int(555),
	PlacedOrder: &datatypes.Billing_Order{
		Id: sl.Int(4444),
		Items: []datatypes.Billing_Order_Item{
			datatypes.Billing_Order_Item{
				Description: sl.String("A Test Item"),
			},
		},
	},
}

var _ = Describe("Volume duplicate", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *file.VolumeDuplicateCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
		FakeStorageManager *testhelpers.FakeStorageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "file")
		cliCommand = file.NewVolumeDuplicateCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
		FakeStorageManager.GetVolumeIdReturns(1234, nil)
	})

	Describe("Volume duplicate", func() {
		Context("Volume duplicate without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("Bad volume id", func() {
			It("return error", func() {
				FakeStorageManager.GetVolumeIdReturns(0, errors.New("BAD Volume ID"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "ZZZ")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("BAD Volume ID"))
			})
		})
		Context("Volume duplicate with 0 DuplicateSnapshotSize", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderDuplicateVolumeReturns(FakeOrderReceipt, nil)
			})
			It("No snapshot size ordered", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--duplicate-snapshot-size", "0", "-f")
				Expect(err).NotTo(HaveOccurred())
				results := fakeUI.Outputs()
				calledWith := FakeStorageManager.OrderDuplicateVolumeArgsForCall(0)
				Expect(calledWith.DuplicateSnapshotSize).To(Equal(0))
				Expect(results).To(ContainSubstring("Order 555 was placed"))
			})
		})
		Context("Volume duplicate without DuplicateSnapshotSize", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderDuplicateVolumeReturns(FakeOrderReceipt, nil)
			})
			It("No snapshot size ordered", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
				Expect(err).NotTo(HaveOccurred())
				results := fakeUI.Outputs()
				calledWith := FakeStorageManager.OrderDuplicateVolumeArgsForCall(0)
				Expect(calledWith.DuplicateSnapshotSize).To(Equal(-1))
				Expect(results).To(ContainSubstring("Order 555 was placed"))
			})
		})
		Context("Ordering Error", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderDuplicateVolumeReturns(FakeOrderReceipt, errors.New("SoftLayer_Exception_ApiError"))
			})
			It("Print Error Output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ApiError"))
			})
		})
	})
})

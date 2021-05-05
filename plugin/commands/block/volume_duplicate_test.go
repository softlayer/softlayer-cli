package block_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
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
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.VolumeDuplicateCommand
		cliCommand         cli.Command
		context            plugin.PluginContext
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		context = plugin.InitPluginContext("softlayer")
		cmd = block.NewVolumeDuplicateCommand(fakeUI, FakeStorageManager, context)
		cliCommand = cli.Command{
			Name:        metadata.BlockVolumeDuplicateMetaData().Name,
			Description: metadata.BlockVolumeDuplicateMetaData().Description,
			Usage:       metadata.BlockVolumeDuplicateMetaData().Usage,
			Flags:       metadata.BlockVolumeDuplicateMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Volume duplicate", func() {
		Context("Volume duplicate without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Bad volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "ZZZ")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Volume ID"))
			})
		})
		Context("Volume duplicate with 0 DuplicateSnapshotSize", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderDuplicateVolumeReturns(FakeOrderReceipt, nil)
			})
			It("No snapshot size ordered", func() {
				err := testhelpers.RunCommand(cliCommand, "12345", "--duplicate-snapshot-size", "0", "-f")
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
				err := testhelpers.RunCommand(cliCommand, "12345", "-f")
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
				err := testhelpers.RunCommand(cliCommand, "12345", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ApiError"))
			})
		})
	})
})

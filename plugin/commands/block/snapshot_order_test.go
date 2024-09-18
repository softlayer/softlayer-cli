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

var _ = Describe("Block Snapshot order", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *block.SnapshotOrderCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
		FakeStorageManager *testhelpers.FakeStorageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewSnapshotOrderCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("Snapshot order", func() {
		Context("Bad Usage", func() {
			It("No Volume ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
			It("Bad Volume ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc", "-s=100")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
			})
			It("No --size", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`required flag(s) "size" not set`))
			})
			It("Bad Tier", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "100", "-t", "0.3")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-t|--tier] is optional, options are: 0.25,2,4,10."))
			})
			It("No confirmation", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "100", "-t", "0.25")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This action will incur charges on your account. Continue?"))
			})
		})

		Context("Snapshot order with -f and continue", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderSnapshotSpaceReturns(datatypes.Container_Product_Order_Receipt{
					OrderId: sl.Int(123456),
					PlacedOrder: &datatypes.Billing_Order{
						Items: []datatypes.Billing_Order_Item{
							datatypes.Billing_Order_Item{
								Description: sl.String("Item1 description"),
							},
							datatypes.Billing_Order_Item{
								Description: sl.String("Item2 description"),
							},
						},
					},
				}, nil)
			})
			It("Normal Order Happy Path", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "100", "-t", "0.25", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
			})
			It("Upgrade Order Happy Path", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "4567", "-s", "1000", "-t", "10", "-u", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				storage_type, volumeId, size, tier, iops, upgrade := FakeStorageManager.OrderSnapshotSpaceArgsForCall(0)
				Expect(storage_type).To(Equal("block"))
				Expect(volumeId).To(Equal(4567))
				Expect(size).To(Equal(1000))
				Expect(tier).To(Equal(10.0))
				Expect(iops).To(Equal(0))
				Expect(upgrade).To(BeTrue())
			})
		})

		Context("Snapshot order with correct parameters but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderSnapshotSpaceReturns(
					datatypes.Container_Product_Order_Receipt{}, errors.New("Internal Server Error"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "100", "-t", "0.25", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to order snapshot space for volume 1234.Please verify your options and try again."))
			})
		})
	})
})

var _ = Describe("File Snapshot order", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *block.SnapshotOrderCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
		FakeStorageManager *testhelpers.FakeStorageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "file")
		cliCommand = block.NewSnapshotOrderCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("Snapshot order", func() {
		Context("Snapshot order with -f and continue", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderSnapshotSpaceReturns(datatypes.Container_Product_Order_Receipt{
					OrderId: sl.Int(123456),
					PlacedOrder: &datatypes.Billing_Order{
						Items: []datatypes.Billing_Order_Item{
							datatypes.Billing_Order_Item{
								Description: sl.String("Item1 description"),
							},
							datatypes.Billing_Order_Item{
								Description: sl.String("Item2 description"),
							},
						},
					},
				}, nil)
			})
			It("Normal Order Happy Path", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "100", "-t", "0.25", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				storage_type, volumeId, size, tier, iops, upgrade := FakeStorageManager.OrderSnapshotSpaceArgsForCall(0)
				Expect(storage_type).To(Equal("file"))
				Expect(volumeId).To(Equal(1234))
				Expect(size).To(Equal(100))
				Expect(tier).To(Equal(0.25))
				Expect(iops).To(Equal(0))
				Expect(upgrade).To(BeFalse())
			})
			It("Upgrade Order Happy Path", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "4567", "-s", "1000", "-t", "10", "-u", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				storage_type, volumeId, size, tier, iops, upgrade := FakeStorageManager.OrderSnapshotSpaceArgsForCall(0)
				Expect(storage_type).To(Equal("file"))
				Expect(volumeId).To(Equal(4567))
				Expect(size).To(Equal(1000))
				Expect(tier).To(Equal(10.0))
				Expect(iops).To(Equal(0))
				Expect(upgrade).To(BeTrue())
			})
		})

		Context("Snapshot order with correct parameters but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderSnapshotSpaceReturns(
					datatypes.Container_Product_Order_Receipt{}, errors.New("Internal Server Error"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-s", "100", "-t", "0.25", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to order snapshot space for volume 1234.Please verify your options and try again."))
			})
		})
	})
})

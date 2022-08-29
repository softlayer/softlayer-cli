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

var _ = Describe("Volume order", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *file.VolumeOrderCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
		FakeStorageManager *testhelpers.FakeStorageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "file")
		cliCommand = file.NewVolumeOrderCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("Volume order", func() {
		Context("Volume order without -t", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -t|--storage-type is required, must be either performance or endurance."))
				Expect(err.Error()).To(ContainSubstring("sl file volume-options' to check available options."))
			})
		})

		Context("Volume order with wrong -t", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -t|--storage-type is required, must be either performance or endurance."))
			})
		})

		Context("Volume order without -s", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "performance")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -s|--size is required, must be a positive integer."))
				Expect(err.Error()).To(ContainSubstring("sl file volume-options' to check available options."))
			})
		})

		Context("Volume order without -d", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "performance", "-s", "1000")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -d|--datacenter is required."))
			})
		})

		Context("Volume order without iops for performance", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "performance", "-s", "1000", "-d", "tok02")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -i|--iops is required with performance volume."))
			})
		})

		Context("Volume order with wrong iops for performance", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "performance", "-s", "1000", "-d", "tok02", "-i", "345")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -i|--iops must be a multiple of 100."))
			})
		})

		Context("Volume order with wrong iops for performance", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "performance", "-s", "1000", "-d", "tok02", "-i", "20")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -i|--iops must be between 100 and 6000, inclusive."))
			})
		})

		Context("Volume order with wrong iops for performance", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "performance", "-s", "1000", "-d", "tok02", "-i", "10000")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -i|--iops must be between 100 and 6000, inclusive."))
			})
		})

		Context("Volume order with correct parameters for performance but server fails", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderVolumeReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "performance", "-s", "1000", "-d", "tok02", "-i", "6000", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to order file volume.Please verify your options and try again."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Volume order with correct parameters for performance but not continue", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "performance", "-s", "1000", "-d", "tok02", "-i", "6000")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This action will incur charges on your account. Continue?"))
			})
		})

		Context("Volume order with correct parameters for performance", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderVolumeReturns(datatypes.Container_Product_Order_Receipt{
					OrderId: sl.Int(12345678),
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
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "performance", "-s", "1000", "-d", "tok02", "--iops", "4000", "-f")
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 12345678 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("sl file volume-list --order 12345678'"))
			})
			It("return no with monthly billing option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "endurance", "-s", "1000", "-d", "tok02", "-e", "0.25", "-f", "-b", "monthly")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 12345678 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("sl file volume-list --order 12345678'"))
			})
			It("return no with monthly billing option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "endurance", "-s", "1000", "-d", "tok02", "-e", "0.25", "-f", "-b", "hourly")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 12345678 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("sl file volume-list --order 12345678'"))
			})
		})

		Context("Volume order without tier for endurance", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "endurance", "-s", "1000", "-d", "tok02")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -e|--tier is required with endurance volume in IOPS/GB, options are: 0.25, 2, 4, 10."))
			})
		})

		Context("Volume order with wrong tier for endurance", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "endurance", "-s", "1000", "-d", "tok02", "-e", "0.5")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -e|--tier is required with endurance volume in IOPS/GB, options are: 0.25, 2, 4, 10."))
				Expect(err.Error()).To(ContainSubstring("sl file volume-options' to check available options."))
			})
		})

		Context("Volume order with correct parameters for endurance but not continue", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "endurance", "-s", "1000", "-d", "tok02", "-e", "0.25")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This action will incur charges on your account. Continue?"))
			})
		})

		Context("Volume order with correct parameters for endurance but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderVolumeReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "endurance", "-s", "1000", "-d", "tok02", "-e", "0.25", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to order file volume.Please verify your options and try again."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})
		Context("Volume order with wrong billing option", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "endurance", "-s", "1000", "-d", "tok02", "-e", "0.25", "-f", "-b", "worngoption")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -b|--billing can only be either hourly or monthly."))
				Expect(err.Error()).To(ContainSubstring("sl file volume-options' to check available options."))
			})

		})
		Context("Volume order with correct parameters for endurance", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderVolumeReturns(datatypes.Container_Product_Order_Receipt{
					OrderId: sl.Int(12345678),
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
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "endurance", "-s", "1000", "-d", "tok02", "-e", "0.25", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 12345678 was placed."))
			})
		})
	})
})

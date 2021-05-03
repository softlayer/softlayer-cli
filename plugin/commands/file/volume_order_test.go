package file_test

import (
	"errors"
	"fmt"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/file"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Volume order", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *file.VolumeOrderCommand
		cliCommand         cli.Command
		context            plugin.PluginContext
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		context = plugin.InitPluginContext("softlayer")
		cmd = file.NewVolumeOrderCommand(fakeUI, FakeStorageManager, context)
		cliCommand = cli.Command{
			Name:        metadata.FileVolumeOrderMetaData().Name,
			Description: metadata.FileVolumeOrderMetaData().Description,
			Usage:       metadata.FileVolumeOrderMetaData().Usage,
			Flags:       metadata.FileVolumeOrderMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Volume order", func() {
		Context("Volume order without -t", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: -t|--storage-type is required, must be either performance or endurance.")).To(BeTrue())
				Expect(err.Error()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl file volume-options' to check available options.", cmd.Context.CLIName())}))
			})
		})

		Context("Volume order with wrong -t", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: -t|--storage-type is required, must be either performance or endurance.")).To(BeTrue())
				Expect(err.Error()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl file volume-options' to check available options.", cmd.Context.CLIName())}))
			})
		})

		Context("Volume order without -s", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "performance")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: -s|--size is required, must be a positive integer.")).To(BeTrue())
				Expect(err.Error()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl file volume-options' to check available options.", cmd.Context.CLIName())}))
			})
		})

		Context("Volume order without -d", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "performance", "-s", "1000")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: -d|--datacenter is required.")).To(BeTrue())
				Expect(err.Error()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl file volume-options' to check available options.", cmd.Context.CLIName())}))
			})
		})

		Context("Volume order without iops for performance", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "performance", "-s", "1000", "-d", "tok02")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: -i|--iops is required with performance volume.")).To(BeTrue())
				Expect(err.Error()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl file volume-options' to check available options.", cmd.Context.CLIName())}))
			})
		})

		Context("Volume order with wrong iops for performance", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "performance", "-s", "1000", "-d", "tok02", "-i", "345")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: -i|--iops must be a multiple of 100.")).To(BeTrue())
				Expect(err.Error()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl file volume-options' to check available options.", cmd.Context.CLIName())}))
			})
		})

		Context("Volume order with wrong iops for performance", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "performance", "-s", "1000", "-d", "tok02", "-i", "20")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: -i|--iops must be between 100 and 6000, inclusive.")).To(BeTrue())
				Expect(err.Error()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl file volume-options' to check available options.", cmd.Context.CLIName())}))
			})
		})

		Context("Volume order with wrong iops for performance", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "performance", "-s", "1000", "-d", "tok02", "-i", "10000")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: -i|--iops must be between 100 and 6000, inclusive.")).To(BeTrue())
				Expect(err.Error()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl file volume-options' to check available options.", cmd.Context.CLIName())}))
			})
		})

		Context("Volume order with correct parameters for performance but server fails", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderVolumeReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "performance", "-s", "1000", "-d", "tok02", "-i", "6000", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to order file volume.Please verify your options and try again.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Volume order with correct parameters for performance but not continue", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "-t", "performance", "-s", "1000", "-d", "tok02", "-i", "6000")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This action will incur charges on your account. Continue?"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted."}))
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
				err := testhelpers.RunCommand(cliCommand, "-t", "performance", "-s", "1000", "-d", "tok02", "--iops", "4000", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Order 12345678 was placed."}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Item1 description"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Item2 description"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{fmt.Sprintf("You may run '%s sl file volume-list --order 12345678' to find this file volume after it is ready.", cmd.Context.CLIName())}))
			})
			It("return no with monthly billing option", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "endurance", "-s", "1000", "-d", "tok02", "-e", "0.25", "-f", "-b", "monthly")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Order 12345678 was placed."}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Item1 description"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Item2 description"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{fmt.Sprintf("You may run '%s sl file volume-list --order 12345678' to find this file volume after it is ready.", cmd.Context.CLIName())}))
			})
			It("return no with monthly billing option", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "endurance", "-s", "1000", "-d", "tok02", "-e", "0.25", "-f", "-b", "hourly")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Order 12345678 was placed."}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Item1 description"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Item2 description"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{fmt.Sprintf("You may run '%s sl file volume-list --order 12345678' to find this file volume after it is ready.", cmd.Context.CLIName())}))
			})
		})

		Context("Volume order without tier for endurance", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "endurance", "-s", "1000", "-d", "tok02")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: -e|--tier is required with endurance volume in IOPS/GB, options are: 0.25, 2, 4, 10.")).To(BeTrue())
				Expect(err.Error()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl file volume-options' to check available options.", cmd.Context.CLIName())}))
			})
		})

		Context("Volume order with wrong tier for endurance", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "endurance", "-s", "1000", "-d", "tok02", "-e", "0.5")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: -e|--tier is required with endurance volume in IOPS/GB, options are: 0.25, 2, 4, 10.")).To(BeTrue())
				Expect(err.Error()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl file volume-options' to check available options.", cmd.Context.CLIName())}))
			})
		})

		Context("Volume order with correct parameters for endurance but not continue", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "-t", "endurance", "-s", "1000", "-d", "tok02", "-e", "0.25")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This action will incur charges on your account. Continue?"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted."}))
			})
		})

		Context("Volume order with correct parameters for endurance but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderVolumeReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "endurance", "-s", "1000", "-d", "tok02", "-e", "0.25", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to order file volume.Please verify your options and try again.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
		Context("Volume order with wrong billing option", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-t", "endurance", "-s", "1000", "-d", "tok02", "-e", "0.25", "-f", "-b", "worngoption")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstrings([]string{"Incorrect Usage: -b|--billing can only be either hourly or monthly."}))
				Expect(err.Error()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl file volume-options' to check available options.", cmd.Context.CLIName())}))
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
				err := testhelpers.RunCommand(cliCommand, "-t", "endurance", "-s", "1000", "-d", "tok02", "-e", "0.25", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Order 12345678 was placed."}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Item1 description"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Item2 description"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{fmt.Sprintf("You may run '%s sl file volume-list --order 12345678' to find this file volume after it is ready.", cmd.Context.CLIName())}))
			})
		})
	})
})

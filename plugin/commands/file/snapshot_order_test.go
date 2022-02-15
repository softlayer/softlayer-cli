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
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/file"
	
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Snapshot order", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *file.SnapshotOrderCommand
		cliCommand         cli.Command
		context            plugin.PluginContext
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		context = plugin.InitPluginContext("softlayer")
		cmd = file.NewSnapshotOrderCommand(fakeUI, FakeStorageManager, context)
		cliCommand = cli.Command{
			Name:        file.FileSnapshotOrderMetaData().Name,
			Description: file.FileSnapshotOrderMetaData().Description,
			Usage:       file.FileSnapshotOrderMetaData().Usage,
			Flags:       file.FileSnapshotOrderMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Snapshot order", func() {
		Context("Snapshot order without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Snapshot order with wrong volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Volume ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Snapshot order without -s", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-s|--size] is required.")).To(BeTrue())
				Expect(err.Error()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl file volume-options' to get available options.", cmd.Context.CLIName())}))
			})
		})

		Context("Snapshot order with wrong tier", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "100", "-t", "0.3")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-t|--tier] is optional, options are: 0.25,2,4,10.")).To(BeTrue())
			})
		})

		Context("Snapshot order with -f and not continue", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "100", "-t", "0.25")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This action will incur charges on your account. Continue?"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted"}))
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
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "100", "-t", "0.25", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Order 123456 was placed."}))
			})
		})

		Context("Snapshot order with -f and continue and upgrade", func() {
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
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "100", "-t", "0.25", "-u", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Order 123456 was placed."}))
			})
		})

		Context("Snapshot order with correct parameters but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderSnapshotSpaceReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal Server Error"))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "100", "-t", "0.25", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to order snapshot space for volume 1234.Please verify your options and try again.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
	})
})

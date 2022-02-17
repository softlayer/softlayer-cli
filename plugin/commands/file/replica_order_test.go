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

var _ = Describe("Replica order", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *file.ReplicaOrderCommand
		cliCommand         cli.Command
		context            plugin.PluginContext
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		context = plugin.InitPluginContext("softlayer")
		cmd = file.NewReplicaOrderCommand(fakeUI, FakeStorageManager, context)
		cliCommand = cli.Command{
			Name:        file.FileReplicaOrderMetaData().Name,
			Description: file.FileReplicaOrderMetaData().Description,
			Usage:       file.FileReplicaOrderMetaData().Usage,
			Flags:       file.FileReplicaOrderMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Replicant order", func() {
		Context("Replicant order without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Replicant order with wrong volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Volume ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("Replicant order without -s", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-s|--snapshot-schedule] is required, options are: HOURLY, DAILY, WEEKLY.")).To(BeTrue())
			})
		})

		Context("Replicant order with wrong -s", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "-s", "yearly")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-s|--snapshot-schedule] is required, options are: HOURLY, DAILY, WEEKLY.")).To(BeTrue())
			})
		})

		Context("Replicant order without -d", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "-s", "DAILY")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-d|--datacenter] is required.")).To(BeTrue())
				Expect(err.Error()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl file volume-options' to get available options.", cmd.Context.CLIName())}))
			})
		})

		Context("Replicant order with wrong tier", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "-s", "DAILY", "-d", "dal09", "-t", "0.3")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-t|--tier] is optional, options are: 0.25,2,4,10.")).To(BeTrue())
			})
		})

		Context("Replicant order with wrong iops", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "-s", "DAILY", "-d", "dal09", "-i", "9")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: -i|--iops must be between 100 and 6000, inclusive.")).To(BeTrue())
			})
		})

		Context("Replicant order with wrong iops", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "-s", "DAILY", "-d", "dal09", "-i", "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: -i|--iops must be a multiple of 100.")).To(BeTrue())
			})
		})

		Context("Replicant order with correct volume id but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderReplicantVolumeReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "-s", "DAILY", "-d", "dal09", "-t", "4", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to order replicant for volume 123.Please verify your options and try again.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Replicant order with correct volume id", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderReplicantVolumeReturns(
					datatypes.Container_Product_Order_Receipt{
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
					},
					nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "-s", "DAILY", "-d", "dal09", "-t", "4", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Order 123456 was placed."}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Item1 description"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Item2 description"}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "-s", "DAILY", "-d", "dal09", "-i", "3000", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Order 123456 was placed."}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Item1 description"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Item2 description"}))
			})
		})
	})
})

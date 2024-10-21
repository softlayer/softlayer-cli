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

var _ = Describe("Replica order", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cliCommand         *block.ReplicaOrderCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewReplicaOrderCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
		FakeStorageManager.GetVolumeIdReturns(1234, nil)
	})

	Describe("Replicant order", func() {
		Context("Replicant order without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("Replicant order without -s", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-s|--snapshot-schedule] is required, options are: HOURLY, DAILY, WEEKLY."))
			})
		})

		Context("Replicant order with wrong -s", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "yearly")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-s|--snapshot-schedule] is required, options are: HOURLY, DAILY, WEEKLY."))
			})
		})

		Context("Replicant order without -d", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "DAILY")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-d|--datacenter] is required."))
				Expect(err.Error()).To(ContainSubstring("sl block volume-options' to get available options."))
			})
		})

		Context("Replicant order with wrong tier", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "DAILY", "-d", "dal09", "-t", "0.3")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-t|--tier] is optional, options are: 0.25,2,4,10."))
			})
		})

		Context("Replicant order with wrong iops", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "DAILY", "-d", "dal09", "-i", "9")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -i|--iops must be between 100 and 6000, inclusive."))
			})
		})

		Context("Replicant order with wrong iops", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "DAILY", "-d", "dal09", "-i", "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -i|--iops must be a multiple of 100."))
			})
		})

		Context("Replicant order with wrong os type", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "DAILY", "-d", "dal09", "-t", "4", "-o", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -o|--os-type is optional, options are: HYPER_V,LINUX,VMWARE,WINDOWS_2008,WINDOWS_GPT,WINDOWS,XEN."))
			})
		})

		Context("Replicant order with correct volume id but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderReplicantVolumeReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "DAILY", "-d", "dal09", "-t", "4", "-o", "LINUX", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to order replicant for volume 1234.Please verify your options and try again."))
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
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "DAILY", "-d", "dal09", "-t", "4", "-o", "LINUX", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Item1 description"))

			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "DAILY", "-d", "dal09", "-i", "3000", "-o", "LINUX", "-f")
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Item1 description"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Item2 description"))
			})
		})
	})
})

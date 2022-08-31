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

var _ = Describe("Replica order", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cliCommand         *file.ReplicaOrderCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "file")
		cliCommand = file.NewReplicaOrderCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("Replicant order", func() {
		Context("Replicant order without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("Replicant order with wrong volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
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

		Context("Replicant order with correct volume id but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.OrderReplicantVolumeReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "DAILY", "-d", "dal09", "-t", "4", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to order replicant for volume 123.Please verify your options and try again."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
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
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "DAILY", "-d", "dal09", "-t", "4", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "DAILY", "-d", "dal09", "-i", "3000", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
			})
		})
	})
})

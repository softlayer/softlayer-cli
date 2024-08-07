package hardware_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware billing", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cliCommand          *hardware.BillingCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewBillingCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.HardwareManager = fakeHardwareManager
	})

	Describe("hardware billing", func() {
		Context("hardware billing without id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})

		Context("hardware billing with wrong id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'. It must be a positive integer."))
			})
		})

		Context("hardware billing with server fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetHardwareReturns(datatypes.Hardware_Server{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get hardware server 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})

		Context("hardware billing", func() {
			created, _ := time.Parse(time.RFC3339, "2021-08-30T00:00:00Z")
			BeforeEach(func() {
				fakeHardwareManager.GetHardwareReturns(datatypes.Hardware_Server{
					Hardware: datatypes.Hardware{
						Id: sl.Int(1234),
						BillingItem: &datatypes.Billing_Item_Hardware{
							Billing_Item: datatypes.Billing_Item{
								Id: sl.Int(1234567),
								NextInvoiceChildren: []datatypes.Billing_Item{
									datatypes.Billing_Item{
										Description:                     sl.String("CentOS 7.x (64 bit)"),
										CategoryCode:                    sl.String("os"),
										NextInvoiceTotalRecurringAmount: sl.Float(0.00),
									},
								},
								RecurringFee:                    sl.Float(1000.00),
								NextInvoiceTotalRecurringAmount: sl.Float(1000.00),
								ProvisionTransaction: &datatypes.Provisioning_Version1_Transaction{
									CreateDate: sl.Time(created),
								},
							},
						},
					},
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1000.00"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("CentOS 7.x (64 bit)"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("os"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.00"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2021-08-30T00:00:00Z"))
			})
		})
	})
})

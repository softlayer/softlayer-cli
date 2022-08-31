package virtual_test

import (
	"errors"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"time"
)

var _ = Describe("vs billing", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeVirtualManager *testhelpers.FakeVirtualServerManager
		cmd                *virtual.BillingCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVirtualManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewBillingCommand(fakeUI, fakeVirtualManager)
		cliCommand = cli.Command{
			Name:        virtual.VSBillingMetaData().Name,
			Description: virtual.VSBillingMetaData().Description,
			Usage:       virtual.VSBillingMetaData().Usage,
			Flags:       virtual.VSBillingMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("vs billing", func() {
		Context("vs billing without id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})

		Context("hardware billing with wrong id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abcd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Virtual server ID'. It must be a positive integer."))
			})
		})

		Context("vs billing with server fails", func() {
			BeforeEach(func() {
				fakeVirtualManager.GetInstanceReturns(datatypes.Virtual_Guest{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})

		Context("hardware billing", func() {
			created, _ := time.Parse(time.RFC3339, "2021-08-30T00:00:00Z")
			BeforeEach(func() {
				fakeVirtualManager.GetInstanceReturns(datatypes.Virtual_Guest{
					Id:            sl.Int(1234),
					ProvisionDate: sl.Time(created),
					BillingItem: &datatypes.Billing_Item_Virtual_Guest{
						Billing_Item: datatypes.Billing_Item{
							Id:                              sl.Int(1234567),
							RecurringFee:                    sl.Float(1000.00),
							NextInvoiceTotalRecurringAmount: sl.Float(1000.00),
							Children: []datatypes.Billing_Item{
								datatypes.Billing_Item{
									NextInvoiceTotalRecurringAmount: sl.Float(0.00),
								},
							},
						},
					},
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1000.00"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.00"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2021-08-30T00:00:00Z"))
			})
		})
	})
})

package order_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/order"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("order quote-detail", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeOrderManager *testhelpers.FakeOrderManager
		cmd              *order.QuoteDetailCommand
		cliCommand       cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeOrderManager = new(testhelpers.FakeOrderManager)
		cmd = order.NewQuoteDetailCommand(fakeUI, fakeOrderManager)
		cliCommand = cli.Command{
			Name:        order.OrderQuoteDetailMetaData().Name,
			Description: order.OrderQuoteDetailMetaData().Description,
			Usage:       order.OrderQuoteDetailMetaData().Usage,
			Flags:       order.OrderQuoteDetailMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("order quote-detail", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCommand(cliCommand, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Quote ID'. It must be a positive integer."))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.GetQuoteReturns(datatypes.Billing_Order_Quote{}, errors.New("Failed to get Quote."))
			})
			It("Failed get Quotes", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get Quote."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerQuote := datatypes.Billing_Order_Quote{
					Id:   sl.Int(123456),
					Name: sl.String("quote1"),
					Order: &datatypes.Billing_Order{
						Items: []datatypes.Billing_Order_Item{
							datatypes.Billing_Order_Item{
								Package: &datatypes.Product_Package{
									KeyName: sl.String("CLOUD_SERVER"),
								},
								CategoryCode: sl.String("ram"),
								Description:  sl.String("1 GB"),
								Quantity:     sl.Int(2),
								RecurringFee: sl.Float(0.00),
								OneTimeFee:   sl.Float(0.00),
							},
						},
					},
				}
				fakeOrderManager.GetQuoteReturns(fakerQuote, nil)
			})
			It("Display quote", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("quote1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("CLOUD_SERVER"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ram"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.00"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.00"))
			})
		})

	})
})

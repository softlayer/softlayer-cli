package order_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/order"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("order quote-save", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeOrderManager *testhelpers.FakeOrderManager
		cmd              *order.QuoteSaveCommand
		cliCommand       cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeOrderManager = new(testhelpers.FakeOrderManager)
		cmd = order.NewQuoteSaveCommand(fakeUI, fakeOrderManager)
		cliCommand = cli.Command{
			Name:        order.OrderQuoteSaveMetaData().Name,
			Description: order.OrderQuoteSaveMetaData().Description,
			Usage:       order.OrderQuoteSaveMetaData().Usage,
			Flags:       order.OrderQuoteSaveMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("order quote-save", func() {

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
				fakeOrderManager.SaveQuoteReturns(datatypes.Billing_Order_Quote{}, errors.New("Failed to save Quote."))
			})
			It("Failed save Quote", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to save Quote."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2016-11-25T00:00:00Z")
				modified, _ := time.Parse(time.RFC3339, "2016-12-25T00:00:00Z")
				fakerQuote := datatypes.Billing_Order_Quote{
					Id:         sl.Int(123456),
					Name:       sl.String("quote1"),
					CreateDate: sl.Time(created),
					ModifyDate: sl.Time(modified),
					Status:     sl.String("SAVED"),
				}
				fakeOrderManager.SaveQuoteReturns(fakerQuote, nil)
			})
			It("Save quote", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("quote1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-11-25T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-12-25T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SAVED"))
			})
		})

	})
})

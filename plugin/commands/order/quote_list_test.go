package order_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/order"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("order quote-list", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeOrderManager *testhelpers.FakeOrderManager
		cliCommand       *order.QuoteListCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeOrderManager = new(testhelpers.FakeOrderManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = order.NewQuoteListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.OrderManager = fakeOrderManager
	})

	Describe("order quote-list", func() {

		Context("Return error", func() {
			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.GetActiveQuotesReturns([]datatypes.Billing_Order_Quote{}, errors.New("Failed to get Quotes."))
			})
			It("Failed get Quotes", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get Quotes."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2016-11-25T00:00:00Z")
				modified, _ := time.Parse(time.RFC3339, "2016-12-25T00:00:00Z")
				fakerQuotes := []datatypes.Billing_Order_Quote{
					datatypes.Billing_Order_Quote{
						Id:         sl.Int(111111),
						Name:       sl.String("quote1"),
						CreateDate: sl.Time(created),
						ModifyDate: sl.Time(modified),
						Status:     sl.String("SAVED"),
						Order: &datatypes.Billing_Order{
							Items: []datatypes.Billing_Order_Item{
								datatypes.Billing_Order_Item{
									Package: &datatypes.Product_Package{
										Id:      sl.Int(200),
										KeyName: sl.String("BARE_METAL_SERVER"),
									},
								},
							},
						},
					},
				}
				fakeOrderManager.GetActiveQuotesReturns(fakerQuotes, nil)
			})
			It("List quotes", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("quote1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-11-25T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-12-25T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SAVED"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("BARE_METAL_SERVER"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("200"))
			})
		})

	})
})

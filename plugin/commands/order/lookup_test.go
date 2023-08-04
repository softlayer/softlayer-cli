package order_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"

	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/order"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("order lookup", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *order.LookupCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		fakeOrderManager *testhelpers.FakeOrderManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeOrderManager = new(testhelpers.FakeOrderManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = order.NewLookupCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.OrderManager = fakeOrderManager
	})

	Describe("order lookup", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage : This command requires one argument"))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Order ID'. It must be a positive integer."))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.GetOrderDetailReturns(datatypes.Billing_Order{}, errors.New("Failed to get Order"))
			})
			It("Failed get order", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get Order"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				date, _ := time.Parse(time.RFC3339, "2017-11-08T00:00:00Z")
				fakerOrder := datatypes.Billing_Order{
					Id: sl.Int(123456),
					UserRecord: &datatypes.User_Customer{
						DisplayName: sl.String("Jhon"),
						UserStatus: &datatypes.User_Customer_Status{
							Name: sl.String("Active"),
						},
					},
					CreateDate:        sl.Time(date),
					ModifyDate:        sl.Time(date),
					OrderApprovalDate: sl.Time(date),
					Status:            sl.String("Approved"),
					OrderTotalAmount:  sl.Float(0.00),
					InitialInvoice: &datatypes.Billing_Invoice{
						InvoiceTotalAmount: sl.Float(0.00),
						InvoiceTopLevelItems: []datatypes.Billing_Invoice_Item{
							datatypes.Billing_Invoice_Item{
								Id: sl.Int(111111),
								Category: &datatypes.Product_Item_Category{
									Name: sl.String("Computing Instance"),
								},
								HostName:                sl.String("hostname"),
								DomainName:              sl.String("domain.com"),
								Description:             sl.String("2 x 2.0 GHz or higher Cores"),
								OneTimeAfterTaxAmount:   sl.Float(0.00),
								RecurringAfterTaxAmount: sl.Float(0.00),
								CreateDate:              sl.Time(date),
								Location: &datatypes.Location{
									Name: sl.String("dal13"),
								},
								Children: []datatypes.Billing_Invoice_Item{
									datatypes.Billing_Invoice_Item{
										Category: &datatypes.Product_Item_Category{
											Name: sl.String("RAM"),
										},
										Description: sl.String("4 GB"),
									},
								},
							},
						},
					},
					Items: []datatypes.Billing_Order_Item{
						datatypes.Billing_Order_Item{
							Description: sl.String("2 x 2.0 GHz or higher Cores"),
						},
					},
				}
				fakeOrderManager.GetOrderDetailReturns(fakerOrder, nil)
			})
			It("Display order", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--details")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Jhon (Active)"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Approved"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.00"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2 x 2.0 GHz or higher Cores"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Computing Instance"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2 x 2.0 GHz or higher Cores (hostname.domain.com)"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring(">>>"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("RAM"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("4 GB"))
			})
		})

	})
})

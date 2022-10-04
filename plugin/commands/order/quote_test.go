package order_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/order"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("order quote", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeOrderManager *testhelpers.FakeOrderManager
		fakeImageManager *testhelpers.FakeImageManager
		cliCommand       *order.QuoteCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeOrderManager = new(testhelpers.FakeOrderManager)
		fakeImageManager = new(testhelpers.FakeImageManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = order.NewQuoteCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.OrderManager = fakeOrderManager
		cliCommand.ImageManager = fakeImageManager
	})

	Describe("order quote", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--fqdn=testquote.test.com")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcde", "--fqdn=testquote.test.com")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Quote ID'. It must be a positive integer."))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--fqdn=testquote.test.com", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})

			It("Set --userdata and --userfile", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--fqdn=testquote.test.com", "--userdata=Userdata", "--userfile=tmp/userfile.txt")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '[--userdata]', '[--userfile]' are exclusive."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.GetQuoteReturns(datatypes.Billing_Order_Quote{}, errors.New("Failed to get Quote"))
			})

			It("Failed get Quote", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--fqdn=testquote.test.com")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get Quote"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakerQuote := datatypes.Billing_Order_Quote{}
				fakeOrderManager.GetQuoteReturns(fakerQuote, nil)
				fakeOrderManager.GetRecalculatedOrderContainerReturns(datatypes.Container_Product_Order{}, errors.New("Failed to get Recalculated Order Container"))
			})

			It("Failed get Recalculated Order Container", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--fqdn=testquote.test.com")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get Recalculated Order Container"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakerQuote := datatypes.Billing_Order_Quote{}
				fakerRecalculatedOrderContainer := datatypes.Container_Product_Order{}
				fakeOrderManager.GetQuoteReturns(fakerQuote, nil)
				fakeOrderManager.GetRecalculatedOrderContainerReturns(fakerRecalculatedOrderContainer, nil)
			})

			It("--fqdn option has invalid format", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--verify", "--fqdn=testquote")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is not following <hostname>.<domain.name.tld> --fqdn option format"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakerQuote := datatypes.Billing_Order_Quote{}
				fakerRecalculatedOrderContainer := datatypes.Container_Product_Order{}
				fakeOrderManager.GetQuoteReturns(fakerQuote, nil)
				fakeOrderManager.GetRecalculatedOrderContainerReturns(fakerRecalculatedOrderContainer, nil)
				fakeImageManager.GetImageReturns(datatypes.Virtual_Guest_Block_Device_Template_Group{}, errors.New("Failed to get Image"))
			})

			It("Failed get image", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--verify", "--fqdn=testquote.test.com", "--image=111111")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get Image"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakerQuote := datatypes.Billing_Order_Quote{
					Id:   sl.Int(123456),
					Name: sl.String("quote1"),
					Order: &datatypes.Billing_Order{
						Items: []datatypes.Billing_Order_Item{
							datatypes.Billing_Order_Item{
								Package: &datatypes.Product_Package{
									Id: sl.Int(1105),
								},
							},
						},
					},
				}
				fakerRecalculatedOrderContainer := datatypes.Container_Product_Order{}
				fakeOrderManager.GetQuoteReturns(fakerQuote, nil)
				fakeOrderManager.GetRecalculatedOrderContainerReturns(fakerRecalculatedOrderContainer, nil)
				fakeOrderManager.VerifyOrderReturns(datatypes.Container_Product_Order{}, errors.New("Failed to verify Quote."))
			})

			It("Failed verify quote", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--verify", "--fqdn=testquote.test.com")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to verify Quote."))
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
									Id: sl.Int(1105),
								},
							},
						},
					},
				}
				fakerRecalculatedOrderContainer := datatypes.Container_Product_Order{}
				fakerImage := datatypes.Virtual_Guest_Block_Device_Template_Group{
					GlobalIdentifier: sl.String("a8a34139-2faa-4519-aaaa-aa68a5978a47"),
				}
				fakeOrderManager.GetQuoteReturns(fakerQuote, nil)
				fakeOrderManager.GetRecalculatedOrderContainerReturns(fakerRecalculatedOrderContainer, nil)
				fakeImageManager.GetImageReturns(fakerImage, nil)
				fakeOrderManager.OrderQuoteReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Failed to order Quote."))
			})

			It("Failed order quote", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--fqdn=testquote.test.com")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to order Quote."))
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
									Id: sl.Int(1105),
								},
							},
						},
					},
				}
				fakerRecalculatedOrderContainer := datatypes.Container_Product_Order{}
				fakerVerifyOrder := datatypes.Container_Product_Order{
					UseHourlyPricing: sl.Bool(true),
					Prices: []datatypes.Product_Item_Price{
						datatypes.Product_Item_Price{
							HourlyRecurringFee: sl.Float(0.00),
							Item: &datatypes.Product_Item{
								KeyName:     sl.String("INTEL_INTEL_XEON_8260_2_4_1U"),
								Description: sl.String("Dual Intel Xeon Platinum 8260 (48 Cores, 2.4 GHz)"),
							},
						},
					},
				}
				fakeOrderManager.GetQuoteReturns(fakerQuote, nil)
				fakeOrderManager.GetRecalculatedOrderContainerReturns(fakerRecalculatedOrderContainer, nil)
				fakeOrderManager.VerifyOrderReturns(fakerVerifyOrder, nil)
			})

			It("Verify quote", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--fqdn=testquote.test.com", "--verify", "--quantity=1", "--postinstall=https://mypostinstallscript.com", "--userdata=Myuserdata")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("INTEL_INTEL_XEON_8260_2_4_1U"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Dual Intel Xeon Platinum 8260 (48 Cores, 2.4 GHz)"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.00"))
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
									Id: sl.Int(1105),
								},
							},
						},
					},
				}
				fakerRecalculatedOrderContainer := datatypes.Container_Product_Order{}
				fakerImage := datatypes.Virtual_Guest_Block_Device_Template_Group{
					GlobalIdentifier: sl.String("a8a34139-2faa-4519-aaaa-aa68a5978a47"),
				}
				created, _ := time.Parse(time.RFC3339, "2016-12-25T00:00:00Z")
				fakerOrderQuote := datatypes.Container_Product_Order_Receipt{
					OrderId:   sl.Int(333333),
					OrderDate: sl.Time(created),
					PlacedOrder: &datatypes.Billing_Order{
						Status: sl.String("PENDING_AUTO_APPROVAL"),
					},
				}
				fakeOrderManager.GetQuoteReturns(fakerQuote, nil)
				fakeOrderManager.GetRecalculatedOrderContainerReturns(fakerRecalculatedOrderContainer, nil)
				fakeImageManager.GetImageReturns(fakerImage, nil)
				fakeOrderManager.OrderQuoteReturns(fakerOrderQuote, nil)
			})

			It("Order quote", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--fqdn=testquote.test.com", "--key=111111", "--image=222222", "--complex-type=SoftLayer_Container_Product_Order_Hardware_Server")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("333333"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-12-25T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("PENDING_AUTO_APPROVAL"))
			})
		})
	})
})

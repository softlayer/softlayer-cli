package account_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Account list ItemDetail", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *account.ItemDetailCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
		fakeHandler *testhelpers.FakeTransportHandler
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = account.NewItemDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})
	AfterEach(func() {
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})

	Describe("Account item detail", func() {
		Context("Account item detail, Invalid Usage", func() {
			It("Set command without id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
			It("Set command with id like letters", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Item ID'. It must be a positive integer."))
			})
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Account item detail, correct use", func() {
			It("return account item detail", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1 x 2.0 GHz or higher Core"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Key                   Value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("createDate            2022-01-05T01:19:21Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("description           1 x 2.0 GHz or higher Core"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("FQDN                  test.test.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("hourlyRecurringFee    0.000000"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("hoursUsed             423"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Ordered By            testName (Active)"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Notes                 -"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Location              dal10"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ram                   2 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("guest_disk0           100 GB (SAN)"))
			})
			It("return account item detail in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"1 x 2.0 GHz or higher Core":`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Key":`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value":`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"createDate","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"2022-01-05T01:19:21Z""`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"description","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "1 x 2.0 GHz or higher Core""`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"FQDN","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"test.test.com""`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"hoursUsed","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"423""`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
			It("Test Fallback to GetItemDetailFromInvoiceItem", func() {
				fakeHandler.AddApiError("SoftLayer_Billing_Item", "getObject", 404, "NOT FOUND")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "999")
				Expect(err).NotTo(HaveOccurred())
				output := fakeUI.Outputs()
				apiCalls := fakeHandler.ApiCallLogs
				Expect(len(apiCalls)).To(Equal(2))
				Expect(output).To(ContainSubstring("gw.ibm.me"))
				Expect(output).To(ContainSubstring("Dual Intel Xeon Gold 5218 (32 Cores, 2.30 GHz)"))
				Expect(apiCalls[0]).To(MatchFields(IgnoreExtras, Fields{
					"Service": Equal("SoftLayer_Billing_Item"),
					"Method":  Equal("getObject"),
					"Options": PointTo(MatchFields(IgnoreExtras, Fields{"Id": PointTo(Equal(999))})),
				}))
				Expect(apiCalls[1]).To(MatchFields(IgnoreExtras, Fields{
					"Service": Equal("SoftLayer_Billing_Invoice_Item"),
					"Method":  Equal("getBillingItem"),
					"Options": PointTo(MatchFields(IgnoreExtras, Fields{"Id": PointTo(Equal(999))})),
				}))
			})
		})
	})
})

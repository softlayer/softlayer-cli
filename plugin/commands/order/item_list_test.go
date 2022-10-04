package order_test

import (
	"errors"

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

var _ = Describe("Order item-list", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *order.ItemListCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		fakeOrderManager *testhelpers.FakeOrderManager
	)
	BeforeEach(func() {
		fakeOrderManager = new(testhelpers.FakeOrderManager)
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = order.NewItemListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.OrderManager = fakeOrderManager
	})

	Describe("Order item-list", func() {
		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.ListItemsReturns([]datatypes.Product_Item{}, errors.New("This command requires one argument"))
			})
			It("Argument is not set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.ListItemsReturns([]datatypes.Product_Item{}, errors.New("Failed to list items."))
			})
			It("Package that does not exist is set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list items."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.ListItemsReturns([]datatypes.Product_Item{}, errors.New("Invalid output format, only JSON is supported now."))
			})
			It("Invalid output is set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "BARE_METAL_SERVER", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return no error", func() {
			fakeItemList := []datatypes.Product_Item{}
			BeforeEach(func() {
				fakeItemList = []datatypes.Product_Item{
					datatypes.Product_Item{
						Id: sl.Int(111111),
						ItemCategory: &datatypes.Product_Item_Category{
							CategoryCode: sl.String("bandwidth"),
						},
						KeyName:     sl.String("BANDWIDTH_0_GB_2"),
						Description: sl.String("0 GB Bandwidth Allotment"),
					},
				}
				fakeOrderManager.ListItemsReturns(fakeItemList, nil)
			})

			It("Item list is displayed", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "BARE_METAL_SERVER")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("bandwidth"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("BANDWIDTH_0_GB_2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0 GB Bandwidth Allotment"))
			})

			It("Item list is displayed in json format", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "BARE_METAL_SERVER", "--output=json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"categoryCode": "bandwidth"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"keyName": "BANDWIDTH_0_GB_2"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"description": "0 GB Bandwidth Allotment"`))
			})
		})
	})
})

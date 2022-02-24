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

var _ = Describe("Order item-list", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeOrderManager *testhelpers.FakeOrderManager
		cmd              *order.ItemListCommand
		cliCommand       cli.Command
	)
	BeforeEach(func() {
		fakeOrderManager = new(testhelpers.FakeOrderManager)
		fakeUI = terminal.NewFakeUI()
		cmd = order.NewItemListCommand(fakeUI, fakeOrderManager)
		cliCommand = cli.Command{
			Name:        order.OrderItemListMetaData().Name,
			Description: order.OrderItemListMetaData().Description,
			Usage:       order.OrderItemListMetaData().Usage,
			Flags:       order.OrderItemListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Order item-list", func() {
		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.ListItemsReturns([]datatypes.Product_Item{}, errors.New("This command requires one argument."))
			})
			It("Argument is not set", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.ListItemsReturns([]datatypes.Product_Item{}, errors.New("Failed to list items."))
			})
			It("Package that does not exist is set", func() {
				err := testhelpers.RunCommand(cliCommand, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list items."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.ListItemsReturns([]datatypes.Product_Item{}, errors.New("Invalid output format, only JSON is supported now."))
			})
			It("Invalid output is set", func() {
				err := testhelpers.RunCommand(cliCommand, "BARE_METAL_SERVER", "--output=xml")
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
				err := testhelpers.RunCommand(cliCommand, "BARE_METAL_SERVER")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("bandwidth"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("BANDWIDTH_0_GB_2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0 GB Bandwidth Allotment"))
			})

			It("Item list is displayed in json format", func() {
				err := testhelpers.RunCommand(cliCommand, "BARE_METAL_SERVER", "--output=json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"categoryCode": "bandwidth"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"keyName": "BANDWIDTH_0_GB_2"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"description": "0 GB Bandwidth Allotment"`))
			})
		})
	})
})

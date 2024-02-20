package order_test

import (
	"errors"

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

var _ = Describe("Order category-list", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *order.CategoryListCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		fakeOrderManager *testhelpers.FakeOrderManager
	)
	BeforeEach(func() {
		fakeOrderManager = new(testhelpers.FakeOrderManager)
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = order.NewCategoryListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.OrderManager = fakeOrderManager
	})

	Describe("Order category-list", func() {
		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.ListCategoriesReturns([]datatypes.Product_Package_Order_Configuration{}, errors.New("This command requires one argument"))
			})
			It("Argument is not set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.ListCategoriesReturns([]datatypes.Product_Package_Order_Configuration{}, errors.New("Failed to list categories."))
			})
			It("Package that does not exist is set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list categories."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.ListCategoriesReturns([]datatypes.Product_Package_Order_Configuration{}, errors.New("Invalid output format, only JSON is supported now."))
			})
			It("Invalid output is set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "BARE_METAL_SERVER", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return no error", func() {
			fakeCategoryList := []datatypes.Product_Package_Order_Configuration{}
			BeforeEach(func() {
				fakeCategoryList = []datatypes.Product_Package_Order_Configuration{
					datatypes.Product_Package_Order_Configuration{
						Id: sl.Int(111111),
						ItemCategory: &datatypes.Product_Item_Category{
							Name:         sl.String("Server Security"),
							CategoryCode: sl.String("trusted_platform_module"),
						},
						IsRequired: sl.Int(0),
					},
				}
				fakeOrderManager.ListCategoriesReturns(fakeCategoryList, nil)
			})

			It("Package list is displayed", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "BARE_METAL_SERVER")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Server Security"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("trusted_platform_module"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("N"))
			})

			It("Package list is displayed in json format", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "BARE_METAL_SERVER", "--output=json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"id": 111111,`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"categoryCode": "trusted_platform_module"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"name": "Server Security"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"isRequired": 0`))
			})
		})

		Context("Return no error", func() {
			fakeCategoryList := []datatypes.Product_Package_Order_Configuration{}
			BeforeEach(func() {
				fakeCategoryList = []datatypes.Product_Package_Order_Configuration{
					datatypes.Product_Package_Order_Configuration{
						Id: sl.Int(222222),
						ItemCategory: &datatypes.Product_Item_Category{
							Name:         sl.String("Server"),
							CategoryCode: sl.String("server"),
						},
						IsRequired: sl.Int(1),
					},
				}
				fakeOrderManager.ListCategoriesReturns(fakeCategoryList, nil)
			})

			It("Required package list is displayed", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "BARE_METAL_SERVER", "--required")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Server"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("server"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Y"))
			})
		})
	})
})

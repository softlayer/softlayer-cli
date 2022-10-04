package order_test

import (
	"errors"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/order"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Place", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *order.PlaceQuoteCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		OrderManager     managers.OrderManager
		fakeOrderManager *testhelpers.FakeOrderManager
	)
	BeforeEach(func() {

		filenames := []string{"getDatacenters_1"}
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession(filenames)
		OrderManager = managers.NewOrderManager(fakeSession)
		fakeOrderManager = new(testhelpers.FakeOrderManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = order.NewPlaceQuoteCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.OrderManager = OrderManager
	})
	Describe("order place-quote", func() {
		for k, _ := range order.TYPEMAP {
			Context("successfully"+k, func() {

				k := k
				It("return error", func() {
					fakeUI.Inputs("Y")
					err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type", k, "--name", "test", "--send-email")
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2523413"}))
				})

			})
		}
	})

	Describe("order place-quote", func() {
		for k, _ := range order.TYPEMAP {
			Context("successfully "+k, func() {

				k := k
				It("return in json format", func() {
					fakeUI.Inputs("Y")
					err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type", k, "--output=json")
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeUI.Outputs()).To(ContainSubstring("2523413"))
				})

			})
		}
	})

	Describe("order place-quote", func() {
		for k, _ := range order.TYPEMAP {
			Context("successfully"+k, func() {

				k := k
				It("return error", func() {
					fakeUI.Inputs("Y")
					err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type", k)
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2523413"}))
				})

			})
		}
	})

	Describe("Order place-quote", func() {
		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.PlaceQuoteReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("This command requires three arguments."))
			})
			It("Arguments is not set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires three arguments."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.PlaceQuoteReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Incorrect complex type:"))
			})
			It("Complex type is set with an invalid value", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type=invalid")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect complex type:"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.PlaceQuoteReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("failed reading file"))
			})
			It("Extras is set with an invalid file", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--extras=@invalid", "--complex-type=SoftLayer_Container_Product_Order_Virtual_Guest")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed reading file"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.PlaceQuoteReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Unable to unmarshal extras json:"))
			})
			It("Extras is set with an invalid value", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--extras=invalid", "--complex-type=SoftLayer_Container_Product_Order_Virtual_Guest")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to unmarshal extras json:"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.PlaceQuoteReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Invalid output format, only JSON is supported now."))
			})
			It("Invalid output is set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type=SoftLayer_Container_Product_Order_Virtual_Guest", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid output format, only JSON is supported now."))
			})
		})
	})
})

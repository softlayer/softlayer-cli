package order_test

import (
	"errors"

	"fmt"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/order"
)

var _ = Describe("Place", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *order.PlaceCommand
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
		cliCommand = order.NewPlaceCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.OrderManager = OrderManager
	})

	Describe("order verify", func() {
		for k, _ := range order.TYPEMAP {
			Context("successfully"+k, func() {

				k := k
				It("return no error", func() {
					err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type", k, "--billing=hourly", "--verify")
					Expect(err).NotTo(HaveOccurred())
					fmt.Println(fakeUI.Outputs())
					Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"4_PORTABLE_PUBLIC_IP_ADDRESSES"}))
				})

			})
		}

		for k, _ := range order.TYPEMAP {
			Context("successfully "+k, func() {

				k := k
				It("return in json format", func() {
					err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type", k, "--billing=monthly", "--verify", "--output=json")
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeUI.Outputs()).To(ContainSubstring("4_PORTABLE_PUBLIC_IP_ADDRESSES"))
				})

			})
		}

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.VerifyPlaceOrderReturns(datatypes.Container_Product_Order{}, errors.New("This command requires three arguments."))
			})
			It("Arguments is not set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--verify")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires three arguments."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.VerifyPlaceOrderReturns(datatypes.Container_Product_Order{}, errors.New("--billing can only be either hourly or monthly."))
			})
			It("Billing flag is set with an invalid value", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--verify", "--billing=invalid")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("--billing can only be either hourly or monthly."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.VerifyPlaceOrderReturns(datatypes.Container_Product_Order{}, errors.New("Incorrect complex type"))
			})
			It("Complex type is set with an invalid value", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--verify", "--complex-type=invalid")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect complex type"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.VerifyPlaceOrderReturns(datatypes.Container_Product_Order{}, errors.New("failed reading file"))
			})
			It("Extras is set with an invalid file", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--verify", "--extras=@invalid", "--complex-type=SoftLayer_Container_Product_Order_Virtual_Guest")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed reading file"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.VerifyPlaceOrderReturns(datatypes.Container_Product_Order{}, errors.New("Unable to unmarshal extras json:"))
			})
			It("Extras is set with an invalid value", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--verify", "--extras=invalid", "--complex-type=SoftLayer_Container_Product_Order_Virtual_Guest")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to unmarshal extras json:"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.VerifyPlaceOrderReturns(datatypes.Container_Product_Order{}, errors.New("Invalid output format, only JSON is supported now."))
			})
			It("Invalid output is set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--verify", "--complex-type=SoftLayer_Container_Product_Order_Virtual_Guest", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid output format, only JSON is supported now."))
			})
		})
	})

	Describe("order create", func() {
		for k, _ := range order.TYPEMAP {
			Context("successfully"+k, func() {

				k := k
				It("return no error", func() {
					err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type", k, "-f")
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"11493593"}))
				})

			})
		}

		for k, _ := range order.TYPEMAP {
			Context("successfully "+k, func() {

				k := k
				It("return in json format", func() {
					err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type", k, "-f", "--output=json")
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeUI.Outputs()).To(ContainSubstring("11493593"))
				})

			})
		}

		Context("Return No error", func() {
			BeforeEach(func() {
				fakeUI.Inputs("No")
			})
			It("Aborted place order", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type=SoftLayer_Container_Product_Order_Virtual_Guest")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This action will incur charges on your account. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
			})
		})
	})
})

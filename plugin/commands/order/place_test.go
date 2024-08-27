package order_test

import (
	"errors"

	"fmt"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"

	. "github.com/onsi/ginkgo/v2"
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
		fakeHandler      *testhelpers.FakeTransportHandler
	)
	BeforeEach(func() {
		filenames := []string{"getDatacenters_1"}
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession(filenames)
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)
		OrderManager = managers.NewOrderManager(fakeSession)
		fakeOrderManager = new(testhelpers.FakeOrderManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = order.NewPlaceCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.OrderManager = OrderManager
	})
    AfterEach(func() {
        // Clear API call logs and any errors that might have been set after every test
        fakeHandler.ClearApiCallLogs()
        fakeHandler.ClearErrors()
    })

	Describe("order verify", func() {
		for k, _ := range order.TYPEMAP {
			Context("successfully"+k, func() {

				k := k
				It("return no error with three arguments", func() {
					err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type", k, "--billing=hourly", "--verify")
					Expect(err).NotTo(HaveOccurred())
					fmt.Println(fakeUI.Outputs())
					Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"4_PORTABLE_PUBLIC_IP_ADDRESSES"}))
				})

				It("return no error with more of three arguments", func() {
					err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB", "CITRIX_VDC", "--complex-type", k, "--billing=hourly", "--verify")
					Expect(err).NotTo(HaveOccurred())
					fmt.Println(fakeUI.Outputs())
					Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"4_PORTABLE_PUBLIC_IP_ADDRESSES"}))
				})

			})
		}

		for k, _ := range order.TYPEMAP {
			Context("successfully "+k, func() {

				k := k
				It("return in json format with three arguments", func() {
					err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type", k, "--billing=monthly", "--verify", "--output=json")
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeUI.Outputs()).To(ContainSubstring("4_PORTABLE_PUBLIC_IP_ADDRESSES"))
				})

				It("return in json format with more of three arguments", func() {
					err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB", "CITRIX_VDC", "--complex-type", k, "--billing=monthly", "--verify", "--output=json")
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
				Expect(err.Error()).To(ContainSubstring("requires at least 3 arg(s), only received 0"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.VerifyPlaceOrderReturns(datatypes.Container_Product_Order{}, errors.New("--billing can only be either hourly or monthly."))
			})
			It("Billing flag is set with an invalid value with three arguments", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--verify", "--billing=invalid")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("--billing can only be either hourly or monthly."))
			})

			It("Billing flag is set with an invalid value with more of three arguments", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB", "CITRIX_VDC", "--verify", "--billing=invalid")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("--billing can only be either hourly or monthly."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.VerifyPlaceOrderReturns(datatypes.Container_Product_Order{}, errors.New("Incorrect complex type"))
			})
			It("Complex type is set with an invalid value with three arguments", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--verify", "--complex-type=invalid")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect complex type"))
			})

			It("Complex type is set with an invalid value with more of three arguments", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB", "CITRIX_VDC", "--verify", "--complex-type=invalid")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect complex type"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.VerifyPlaceOrderReturns(datatypes.Container_Product_Order{}, errors.New("failed reading file"))
			})
			It("Extras is set with an invalid file with three arguments", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--verify", "--extras=@invalid", "--complex-type=SoftLayer_Container_Product_Order_Virtual_Guest")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed reading file"))
			})

			It("Extras is set with an invalid file with more of three arguments", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB", "CITRIX_VDC", "--verify", "--extras=@invalid", "--complex-type=SoftLayer_Container_Product_Order_Virtual_Guest")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed reading file"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.VerifyPlaceOrderReturns(datatypes.Container_Product_Order{}, errors.New("Unable to unmarshal extras json:"))
			})
			It("Extras is set with an invalid value with arguments", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--verify", "--extras=invalid", "--complex-type=SoftLayer_Container_Product_Order_Virtual_Guest")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to unmarshal extras json:"))
			})

			It("Extras is set with an invalid value with more of three arguments", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB", "CITRIX_VDC", "--verify", "--extras=invalid", "--complex-type=SoftLayer_Container_Product_Order_Virtual_Guest")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to unmarshal extras json:"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.VerifyPlaceOrderReturns(datatypes.Container_Product_Order{}, errors.New("Invalid output format, only JSON is supported now."))
			})
			It("Invalid output is set with three arguments", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--verify", "--complex-type=SoftLayer_Container_Product_Order_Virtual_Guest", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid output format, only JSON is supported now."))
			})

			It("Invalid output is set with more of three arguments", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB", "CITRIX_VDC", "--verify", "--complex-type=SoftLayer_Container_Product_Order_Virtual_Guest", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid output format, only JSON is supported now."))
			})
		})
	})

	Describe("order create", func() {
		for k, _ := range order.TYPEMAP {
			Context("successfully"+k, func() {

				k := k
				It("return no error with three arguments", func() {
					err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type", k, "-f")
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"11493593"}))
				})

				It("return no error with more of three arguments", func() {
					err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB", "CITRIX_VDC", "--complex-type", k, "-f")
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"11493593"}))
				})

			})
		}

		for k, _ := range order.TYPEMAP {
			Context("successfully "+k, func() {

				k := k
				It("return in json format with three arguments", func() {
					err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13,EVAULT_100_GB", "CITRIX_VDC", "--complex-type", k, "-f", "--output=json")
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeUI.Outputs()).To(ContainSubstring("11493593"))
				})

				It("return in json format with more of three arguments", func() {
					err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB", "CITRIX_VDC", "--complex-type", k, "-f", "--output=json")
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeUI.Outputs()).To(ContainSubstring("11493593"))
				})

			})
		}

		Context("Return No error", func() {
			BeforeEach(func() {
				fakeUI.Inputs("No")
			})
			It("Aborted place order with three arguments", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type=SoftLayer_Container_Product_Order_Virtual_Guest")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This action will incur charges on your account. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
			})

			It("Aborted place order with more of three arguments", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "CLOUD_SERVER", "dal13", "EVAULT_100_GB", "CITRIX_VDC", "--complex-type=SoftLayer_Container_Product_Order_Virtual_Guest")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This action will incur charges on your account. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
			})
		})
	})
	Describe("softlayer-cli/issues/863", func() {
		BeforeEach(func() {
			fakeHandler.ClearApiCallLogs()
			fakeHandler.SetFileNames([]string{"getItems-835", "getDatacenters_mad02"})
		})
		It("Finds the correct price IDs", func() {
			err := testhelpers.RunCobraCommand(
				cliCommand.Command,
				"PUBLIC_CLOUD_SERVER", "MADRID02", "1_GBPS_PRIVATE_NETWORK_UPLINK", "1_IP_ADDRESS",
				"GUEST_DISK_100_GB_LOCAL", "OS_RED_HAT_ENTERPRISE_LINUX_9_X_MINIMAL_INSTALL_64_BIT",
				"MONITORING_HOST_PING", "NOTIFICATION_EMAIL_AND_TICKET", "AUTOMATED_NOTIFICATION",
				"UNLIMITED_SSL_VPN_USERS_1_PPTP_VPN_USER_PER_ACCOUNT", "REBOOT_REMOTE_CONSOLE", "BANDWIDTH_0_GB",
				"--billing=monthly",
				`--extras={"virtualGuests":[{"hostname":"testServer","domain":"ibm.com"}]}`,
				"--complex-type=SoftLayer_Container_Product_Order_Virtual_Guest",
				"--preset=BL2_8x32x100", "--verify",
			)
			Expect(err).NotTo(HaveOccurred())
			
			callLog := fakeHandler.ApiCallLogs
			Expect(len(callLog)).To(Equal(9))
			fmt.Printf(callLog[8].String())
			Expect(callLog[8].String()).To(Equal(`SoftLayer_Product_Order::verifyOrder(id=0, mask='', filter='', ` +
				`{"parameters":[{"complexType":"SoftLayer_Container_Product_Order_Virtual_Guest",` +
				`"location":"3460412","packageId":865,"presetId":785,"prices":[{"id":899},{"id":21},{"id":204637},` +
				`{"id":314158},{"id":55},{"id":57},{"id":58},{"id":420},{"id":905},{"id":22505}],"quantity":1,` +
				`"useHourlyPricing":false,"virtualGuests":[{"domain":"ibm.com","hostname":"testServer"}]}]}`,
			))
		})
	})
})

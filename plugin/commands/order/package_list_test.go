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

var _ = Describe("Order package-list", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *order.PackageListCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		fakeOrderManager *testhelpers.FakeOrderManager
	)
	BeforeEach(func() {
		fakeOrderManager = new(testhelpers.FakeOrderManager)
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = order.NewPackageListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.OrderManager = fakeOrderManager
	})

	Describe("Order package-list", func() {

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.ListPackageReturns([]datatypes.Product_Package{}, errors.New("Invalid output format, only JSON is supported now."))
			})
			It("Invalid output is set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return no error", func() {
			fakePackageList := []datatypes.Product_Package{}
			BeforeEach(func() {
				fakePackageList = []datatypes.Product_Package{
					datatypes.Product_Package{
						Id:      sl.Int(56),
						Name:    sl.String("Quad Processor Multi Core Nehalem EX"),
						KeyName: sl.String("ADDITIONAL_PRODUCTS"),
						Type: &datatypes.Product_Package_Type{
							KeyName: sl.String("BARE_METAL_CPU"),
						},
					},
				}
				fakeOrderManager.ListPackageReturns(fakePackageList, nil)
			})

			It("Package list is displayed", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("56"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Quad Processor Multi Core Nehalem EX"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ADDITIONAL_PRODUCTS"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("BARE_METAL_CPU"))
			})

			It("Package list is displayed in json format", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"id": 56`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"name": "Quad Processor Multi Core Nehalem EX"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"keyName": "BARE_METAL_CPU"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"keyName": "ADDITIONAL_PRODUCTS"`))
			})
		})
	})
})

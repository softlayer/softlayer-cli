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

var _ = Describe("Order package-locations", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *order.PackageLocationCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		fakeOrderManager *testhelpers.FakeOrderManager
	)
	BeforeEach(func() {
		fakeOrderManager = new(testhelpers.FakeOrderManager)
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = order.NewPackageLocationCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.OrderManager = fakeOrderManager
	})

	Describe("Order package-locations", func() {
		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.PackageLocationReturns([]datatypes.Location_Region{}, errors.New("This command requires one argument"))
			})
			It("Argument is not set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.PackageLocationReturns([]datatypes.Location_Region{}, errors.New("Failed to list package locations."))
			})
			It("Package that does not exist is set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "BARE_METAL_SER")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list package locations."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.PackageLocationReturns([]datatypes.Location_Region{}, errors.New("Invalid output format, only JSON is supported now."))
			})
			It("Invalid output is set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "BARE_METAL_SERVER", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return no error", func() {
			fakeLocationList := []datatypes.Location_Region{}
			BeforeEach(func() {
				fakeLocationList = []datatypes.Location_Region{
					datatypes.Location_Region{
						Locations: []datatypes.Location_Region_Location{
							datatypes.Location_Region_Location{
								Location: &datatypes.Location{
									Id:   sl.Int(265592),
									Name: sl.String("ams01"),
								},
							},
						},
						Description: sl.String("AMS01 - Amsterdam"),
						Keyname:     sl.String("AMSTERDAM"),
					},
				}
				fakeOrderManager.PackageLocationReturns(fakeLocationList, nil)
			})

			It("Location list is displayed", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "BARE_METAL_SERVER")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("265592"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ams01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("AMS01 - Amsterdam"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("AMSTERDAM"))
			})

			It("Location list is displayed in json format", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "BARE_METAL_SERVER", "--output=json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"id": 265592,`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"name": "ams01"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"description": "AMS01 - Amsterdam"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"keyname": "AMSTERDAM"`))
			})
		})
	})
})

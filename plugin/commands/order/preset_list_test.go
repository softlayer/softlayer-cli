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

var _ = Describe("Order preset-list", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *order.PresetListCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		fakeOrderManager *testhelpers.FakeOrderManager
	)
	BeforeEach(func() {
		fakeOrderManager = new(testhelpers.FakeOrderManager)
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = order.NewPresetListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.OrderManager = fakeOrderManager
	})

	Describe("Order preset-list", func() {
		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.ListPresetReturns([]datatypes.Product_Package_Preset{}, errors.New("This command requires one argument"))
			})
			It("Argument is not set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.ListPresetReturns([]datatypes.Product_Package_Preset{}, errors.New("Failed to list presets"))
			})
			It("Package that does not exist is set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "BARE_METAL_SER")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list presets"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeOrderManager.ListPresetReturns([]datatypes.Product_Package_Preset{}, errors.New("Invalid output format, only JSON is supported now."))
			})
			It("Invalid output is set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "BARE_METAL_SERVER", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return no error", func() {
			fakePresetList := []datatypes.Product_Package_Preset{}
			BeforeEach(func() {
				fakePresetList = []datatypes.Product_Package_Preset{
					datatypes.Product_Package_Preset{
						Id:          sl.Int(1278),
						Name:        sl.String("1U 4210 384GB 2x4TB RAID1"),
						KeyName:     sl.String("1U_4210S_384GB_2X4TB_RAID_1"),
						Description: sl.String("1U 4210 384GB 2x4TB RAID1"),
					},
				}
				fakeOrderManager.ListPresetReturns(fakePresetList, nil)
			})

			It("Location list is displayed", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "BARE_METAL_SERVER")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1U 4210 384GB 2x4TB RAID1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1U_4210S_384GB_2X4TB_RAID_1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1U 4210 384GB 2x4TB RAID1"))
			})

			It("Location list is displayed in json format", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "BARE_METAL_SERVER", "--output=json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"id": 1278`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"name": "1U 4210 384GB 2x4TB RAID1"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"keyName": "1U_4210S_384GB_2X4TB_RAID_1"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"description": "1U 4210 384GB 2x4TB RAID1"`))
			})
		})
	})
})

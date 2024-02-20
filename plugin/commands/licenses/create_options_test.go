package licenses_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/licenses"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("licenses create-options", func() {
	var (
		fakeUI              *terminal.FakeUI
		cliCommand          *licenses.LicensesOptionsCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
		fakeLicensesManager *testhelpers.FakeLicensesManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = licenses.NewLicensesOptionsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLicensesManager = new(testhelpers.FakeLicensesManager)
		cliCommand.LicensesManager = fakeLicensesManager

		fakeLicensesOptions := []datatypes.Product_Package{
			datatypes.Product_Package{
				Items: []datatypes.Product_Item{
					datatypes.Product_Item{
						Id:          sl.Int(123),
						Description: sl.String("item description"),
						KeyName:     sl.String("ITEM_KEY_NAME"),
						Capacity:    sl.Float(28.0),
						Prices: []datatypes.Product_Item_Price{
							datatypes.Product_Item_Price{
								RecurringFee: sl.Float(252.0),
							},
						},
					},
					datatypes.Product_Item{
						Id:          sl.Int(1234),
						Description: sl.String("item description 2"),
						KeyName:     sl.String("ITEM_KEY_NAME_2"),
						Capacity:    sl.Float(40.0),
						Prices: []datatypes.Product_Item_Price{
							datatypes.Product_Item_Price{
								RecurringFee: sl.Float(100.0),
							},
						},
					},
				},
			},
		}
		fakeLicensesManager.CreateLicensesOptionsReturns(fakeLicensesOptions, nil)

	})

	Context("licenses create options returns error", func() {
		BeforeEach(func() {
			fakeLicensesManager.CreateLicensesOptionsReturns(nil, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to licenses create options.\nInternal server error")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})

	Context("licenses create options", func() {
		It("return licenses create options", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("Id     Description          KeyName           Capacity    RecurringFee"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("123    item description     ITEM_KEY_NAME     28.000000   252.000000"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("1234   item description 2   ITEM_KEY_NAME_2   40.000000   100.000000"))
		})
	})
})

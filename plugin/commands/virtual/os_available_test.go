package virtual_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS os-available", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *virtual.OsAvailableCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		fakeOrderManager *testhelpers.FakeOrderManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeOrderManager = new(testhelpers.FakeOrderManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewOsAvailableCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.OrderManager = fakeOrderManager
	})

	Describe("VS os-available", func() {
		Context("VS list with server fails", func() {
			BeforeEach(func() {
				fakeOrderManager.ListItemsReturns([]datatypes.Product_Item{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list available OS's."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("list available OS's", func() {
			BeforeEach(func() {
				fakerItems := []datatypes.Product_Item{
					datatypes.Product_Item{
						Id:          sl.Int(3975),
						KeyName:     sl.String("OS_DEBIAN_6_X_SQUEEZE_LAMP_64_BIT"),
						Description: sl.String("Debian GNU/Linux 6.x Squeeze/Stable - LAMP Install (64 bit)"),
						Prices: []datatypes.Product_Item_Price{
							datatypes.Product_Item_Price{
								HourlyRecurringFee: sl.Float(0.0),
								RecurringFee:       sl.Float(0.2),
								SetupFee:           sl.Float(0.4),
							},
						},
					},
				}
				fakeOrderManager.ListItemsReturns(fakerItems, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("3975"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("OS_DEBIAN_6_X_SQUEEZE_LAMP_64_BIT"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Debian GNU/Linux 6.x Squeeze/Stable - LAMP Install (64 bit)"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.0"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.4"))
			})
		})
	})
})

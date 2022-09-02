package virtual_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS capacity create", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *virtual.CapacityCreateCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeVSManager *testhelpers.FakeVirtualServerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewCapacityCreateCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
	})

	Describe("VS capacity create", func() {
		Context("VS create with incorrect parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--flavor", "C1_1X1X100", "-c", "1")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown shorthand flag: 'c' in -c"))
			})
		})
		Context("VS create with no parameters", func() {
			It("return error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command)

				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted"))
			})
		})
		Context("VS create with correct parameters", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2017-11-08T00:00:00Z")
				fakeVSManager.GenerateInstanceCapacityCreationTemplateReturns(
					datatypes.Container_Product_Order_Receipt{
						OrderDate:    sl.Time(created),
						OrderId:      sl.Int(991122),
						PlacedOrder:  &datatypes.Billing_Order{Status: sl.String("OkGood")},
						OrderDetails: &datatypes.Container_Product_Order{PostTaxRecurringHourly: sl.Float(99.11)},
					},
					nil,
				)
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--flavor", "C1_1X1X100", "-i", "1", "--backendRouterId", "1234", "--name", "test")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order Date     2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Status         OkGood"))
			})
		})
	})
})

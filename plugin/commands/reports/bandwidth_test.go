package reports_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/reports"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("reports bandwidth", func() {
	var (
		fakeUI            *terminal.FakeUI
		cliCommand        *reports.BandwidthCommand
		fakeSession       *session.Session
		slCommand         *metadata.SoftlayerCommand
		fakeSearchManager *testhelpers.FakeSearchManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = reports.NewBandwidthCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeSearchManager = new(testhelpers.FakeSearchManager)
		cliCommand.SearchManager = fakeSearchManager
	})

	Describe("reports bandwidth", func() {
		Context("Return error", func() {
			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeSearchManager.AdvancedSearchReturns([]datatypes.Container_Search_Result{}, errors.New("Internal Error"))
			})
			It("Failed to get bandwidth summary", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get bandwidth summary"))
			})
		})

		Context("Advanced search to bandwidth summary", func() {
			BeforeEach(func() {
				filename := []string{"bandwidth"}
				fakeSession = testhelpers.NewFakeSoftlayerSession(filename)
				slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
				cliCommand = reports.NewBandwidthCommand(slCommand)
				cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
			})
			It("return bandwidth summary", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Device name"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("test.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SLADC307608-jt48"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Allocation"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("250.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("20.00 TB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Unlimited"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Pay-As-You-Go"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Data in"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("646.95 MB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("417.65 MB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1.04 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Not Applicable"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("tag,test,tag,test2"))
			})
			It("return bandwidth summary in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Id": "100250634",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Allocation": "250.00 GB",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Pool": "Virtual Private Rack",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Tags": ""`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})

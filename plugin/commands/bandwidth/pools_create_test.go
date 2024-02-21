package bandwidth_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/bandwidth"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("bandwidth pools-create", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *bandwidth.PoolsCreateCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = bandwidth.NewPoolsCreateCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("bandwidth pools-create", func() {
		Context("Return error", func() {
			It("Set command without Name and Region", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`required flag(s) "name", "region" not set`))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--name=NameRegion", "--region=SJC/DAL/WDC/TOR/MON", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return no error", func() {
			It("Get Bandwidth Pool with devices", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--name=NameRegion", "--region=SJC/DAL/WDC/TOR/MON")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Id"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456789"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name Pool"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("NewRegion"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SJC/DAL/WDC/TOR/MON"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Region"))
			})
			It("Get Bandwidth Pool with devices", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--name=NameRegion", "--region=SJC/DAL/WDC/TOR/MON", "--output=json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Id",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "123456789"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "NewRegion"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "SJC/DAL/WDC/TOR/MON"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})

package account_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Account Bandwidth-Pools", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *account.BandwidthPoolsCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = account.NewBandwidthPoolsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("Bandwidth-Pools Testing", func() {
		Context("Happy Path", func() {
			It("Runs without issue", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("MexRegion"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Region"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("MEX"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Allocation"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("3361 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Current Usage"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("7.70 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Cost"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("$25"))
			})

			It("Outputs JSON", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=JSON")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"ID": "265721",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "TestPool",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Projected Usage": "-",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Cost": "$55"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})

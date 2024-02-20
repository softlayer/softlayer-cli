package placementgroup_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/placementgroup"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("placementgroup credentials", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *placementgroup.PlacementGroupListCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = placementgroup.NewPlacementGroupListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("placementgroup list", func() {
		Context("Placementgroup list, correct use", func() {
			It("return placementgroup", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID       Name            Backend Router   Rule        Guests   Created"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("140665   dal05-ibmtest   -                testRule1   -        2019-06-07T19:34:55Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("71643    TestGroup       -                testRule1   -        2019-01-30T23:53:00Z"))
			})
			It("return placementgroup list in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"accountId": 99999,`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"backendRouterId": 122762,`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"rule": {`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"ruleId": 1`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
			})
		})
	})
})

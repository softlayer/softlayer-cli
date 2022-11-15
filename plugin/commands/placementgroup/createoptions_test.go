package placementgroup_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/placementgroup"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("placementgroup credentials", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *placementgroup.PlacementGroupCreateOptionsCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = placementgroup.NewPlacementGroupCreateOptionsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("placementgroup create options options", func() {
		Context("Placementgroup create options, correct use", func() {
			It("return placementgroup", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Available Router:"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Data Center   Hostname      Backend Router Id"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Dallas 1      bcr01.dal01   1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Rules:"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID   Name"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1    SPREAD"))
			})
		})
	})
})

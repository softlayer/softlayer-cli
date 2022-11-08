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
		cliCommand  *placementgroup.PlacementGroupDetailCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = placementgroup.NewPlacementGroupDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("placementgroup detail options", func() {
		Context("Return error", func() {
			It("Set command without ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
			It("Set command without ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Placement Group ID'. It must be a positive integer."))
			})
		})

		Context("Placementgroup detail, correct use", func() {
			It("return placementgroup", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name             Value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID               1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name             test-group"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Backend Router   bcr01a.mex01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Rule             SPREAD"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Created          2019-01-17T20:36:42Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Guests:          ID         FQDN                            Primary IP      Backend IP     CPU   Memory   Provisioned"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("69131875   issues10691547765077.test.com   169.57.70.180   10.131.11.14   1     1024     2019-01-17T22:47:17Z"))
			})
			It("return placementgroup detail in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345", "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"backendRouter":`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"guests":`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
			})
		})
	})
})

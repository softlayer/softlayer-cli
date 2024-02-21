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
		cliCommand  *placementgroup.PlacementGroupCreateCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = placementgroup.NewPlacementGroupCreateCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("placementgroup create", func() {
		Context("Return error", func() {
			It("Set command without name", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-n, --name' is required"))
			})

			It("Set command without backend-router-id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--name", "name")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-b, --backend-router-id' is required"))
			})

			It("Set command without backend-router-id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--name", "name", "--backend-router-id", "12345")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-r, --rule-id' is required"))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--name", "name", "--backend-router-id", "12345", "--rule-id", "54321", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format"))
			})
		})

		Context("Placementgroup create, correct use", func() {
			It("return placementgroup", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--name", "name", "--backend-router-id", "12345", "--rule-id", "54321")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Successfully created placement group: ID: 5555, Name: test01."))
			})

			It("return placementgroup in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--name", "name", "--backend-router-id", "12345", "--rule-id", "54321", "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"accountId": 123,`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"backendRouterId": 444,`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
			})
		})
	})
})

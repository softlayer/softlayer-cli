package autoscale_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/autoscale"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("autoscale detail", func() {
	var (
		fakeUI               *terminal.FakeUI
		cliCommand           *autoscale.DetailCommand
		fakeSession          *session.Session
		slCommand            *metadata.SoftlayerCommand
		fakeAutoScaleManager *testhelpers.FakeAutoScaleManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeAutoScaleManager = new(testhelpers.FakeAutoScaleManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = autoscale.NewDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.AutoScaleManager = fakeAutoScaleManager
	})

	Describe("autoscale detail", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Autoscale Group ID'. It must be a positive integer."))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeAutoScaleManager.GetScaleGroupReturns(datatypes.Scale_Group{}, errors.New("Failed to get scale group."))
			})
			It("Failed get scale group", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get scale group."))
			})
		})
	})
})

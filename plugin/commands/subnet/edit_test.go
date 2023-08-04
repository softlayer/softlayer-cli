package subnet_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/subnet"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("subnet edit", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *subnet.EditCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeNetworkManager *testhelpers.FakeNetworkManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = subnet.NewEditCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cliCommand.NetworkManager = fakeNetworkManager
	})

	Describe("subnet edit", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage : This command requires one argument"))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Subnet ID'. It must be a positive integer."))
			})

			It("Set command without option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Please pass at least one of the flags."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeNetworkManager.SetSubnetTagsReturns(false, errors.New("Failed to set tags"))
			})
			It("Failed set tags", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--tags=tag1")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to set tags"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeNetworkManager.SetSubnetNoteReturns(false, errors.New("Failed to set note"))
			})
			It("Failed set note", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--note=myNote")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to set note"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeNetworkManager.SetSubnetNoteReturns(true, nil)
				fakeNetworkManager.SetSubnetTagsReturns(true, nil)
			})

			It("Set Tags and note", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--tags=mytag1", "--note=myNote")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Set tags successfully"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Set note successfully"))
			})
		})
	})
})

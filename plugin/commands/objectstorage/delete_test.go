package objectstorage_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/objectstorage"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Object Storage delete", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *objectstorage.DeleteCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = objectstorage.NewDeleteCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("Object Storage delete", func() {
		Context("Object Storage delete, Invalid Usage", func() {
			It("Set command without id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument"))
			})
			It("Set command with id like letters", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'objectStorageID'. It must be a positive integer."))
			})
			It("Set command without confirmation", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This will cancel the object-storage:"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("and cannot be undone. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
			})
		})
		Context("Object Storage delete , correct use", func() {
			It("return objectstorage delete", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--force")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Object-storage"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("has been marked for cancellation."))
			})
			It("return objectstorage delete with inmediate option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--force", "--immediate", "--reason", "cancel test")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Object-storage"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("has been marked for immediate cancellation."))
			})
		})
	})
})

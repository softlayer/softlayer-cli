package objectstorage_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/objectstorage"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Object Storage delete Object Storages", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *objectstorage.CredentialDeleteCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = objectstorage.NewCredentialDeleteCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("Object Storage delete credential", func() {
		Context("Object Storage delete credential, Invalid Usage", func() {
			It("Set command without id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument"))
			})
			It("Set command with id like letters", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Storage ID'. It must be a positive integer."))
			})
			It("Set command with id but whitout credentialID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--credential-id' is required"))
			})
		})
		Context("Object Storage delete credential, correct use", func() {
			It("return objectstorage delete credential", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--credential-id", "654321")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Credential: 654321 was deleted."))
			})
		})
	})
})

var _ = Describe("Object Storage delete Object Storages", func() {
	var (
		fakeUI                   *terminal.FakeUI
		cliCommand               *objectstorage.CredentialDeleteCommand
		fakeSession              *session.Session
		slCommand                *metadata.SoftlayerCommand
		fakeObjectStorageManager *testhelpers.FakeObjectStorageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = objectstorage.NewCredentialDeleteCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeObjectStorageManager = new(testhelpers.FakeObjectStorageManager)
		cliCommand.ObjectStorageManager = fakeObjectStorageManager
	})

	Describe("Object Storage delete credential", func() {
		Context("Object Storage delete credential, errors", func() {
			It("return error with storageID", func() {
				fakeObjectStorageManager.DeleteCredentialReturns(errors.New("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123456'. (HTTP 404)"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--credential-id", "654321")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to find object-storage with ID: 123456."))
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123456'. (HTTP 404)"))
			})
			It("return error with credentialID", func() {
				fakeObjectStorageManager.DeleteCredentialReturns(errors.New("ObjectNotFound: Unable to find object with id of '654321'. (HTTP 404)"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--credential-id", "654321")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to find credential with ID: 654321."))
				Expect(err.Error()).To(ContainSubstring("ObjectNotFound: Unable to find object with id of '654321'. (HTTP 404)"))
			})
			It("return generic error", func() {
				fakeObjectStorageManager.DeleteCredentialReturns(errors.New("Internal Server Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--credential-id", "654321")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to delete credential: 123456."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

	})
})

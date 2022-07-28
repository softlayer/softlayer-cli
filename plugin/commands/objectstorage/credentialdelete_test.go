package objectstorage_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/objectstorage"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Object Storage delete Object Storages", func() {
	var (
		fakeUI                   *terminal.FakeUI
		cmd                      *objectstorage.CredentialDeleteCommand
		cliCommand               cli.Command
		fakeSession              *session.Session
		fakeObjectStorageManager managers.ObjectStorageManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeObjectStorageManager = managers.NewObjectStorageManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		cmd = objectstorage.NewCredentialDeleteCommand(fakeUI, fakeObjectStorageManager)
		cliCommand = cli.Command{
			Name:        objectstorage.CredentialDeleteMetaData().Name,
			Description: objectstorage.CredentialDeleteMetaData().Description,
			Usage:       objectstorage.CredentialDeleteMetaData().Usage,
			Flags:       objectstorage.CredentialDeleteMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Object Storage delete credential", func() {
		Context("Object Storage delete credential, Invalid Usage", func() {
			It("Set command without id", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument."))
			})
			It("Set command with id like letters", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Storage ID'. It must be a positive integer."))
			})
			It("Set command with id but whitout credentialID", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--credential-id' is required"))
			})
		})
		Context("Object Storage delete credential, correct use", func() {
			It("return objectstorage delete credential", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--credential-id", "654321")
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
		fakeObjectStorageManager *testhelpers.FakeObjectStorageManager
		cmd                      *objectstorage.CredentialDeleteCommand
		cliCommand               cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeObjectStorageManager = new(testhelpers.FakeObjectStorageManager)

		cmd = objectstorage.NewCredentialDeleteCommand(fakeUI, fakeObjectStorageManager)
		cliCommand = cli.Command{
			Name:        objectstorage.CredentialDeleteMetaData().Name,
			Description: objectstorage.CredentialDeleteMetaData().Description,
			Usage:       objectstorage.CredentialDeleteMetaData().Usage,
			Flags:       objectstorage.CredentialDeleteMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Object Storage delete credential", func() {
		Context("Object Storage delete credential, errors", func() {
			It("return error with storageID", func() {
				fakeObjectStorageManager.DeleteCredentialReturns(errors.New("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123456'. (HTTP 404)"))
				err := testhelpers.RunCommand(cliCommand, "123456", "--credential-id", "654321")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to find object-storage with ID: 123456."))
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123456'. (HTTP 404)"))
			})
			It("return error with credentialID", func() {
				fakeObjectStorageManager.DeleteCredentialReturns(errors.New("ObjectNotFound: Unable to find object with id of '654321'. (HTTP 404)"))
				err := testhelpers.RunCommand(cliCommand, "123456", "--credential-id", "654321")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to find credential with ID: 654321."))
				Expect(err.Error()).To(ContainSubstring("ObjectNotFound: Unable to find object with id of '654321'. (HTTP 404)"))
			})
			It("return generic error", func() {
				fakeObjectStorageManager.DeleteCredentialReturns(errors.New("Internal Server Error"))
				err := testhelpers.RunCommand(cliCommand, "123456", "--credential-id", "654321")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to delete credential: 123456."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

	})
})

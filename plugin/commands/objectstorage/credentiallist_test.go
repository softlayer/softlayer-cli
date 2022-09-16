package objectstorage_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/objectstorage"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Object Storage list Object Storages", func() {
	var (
		fakeUI                   *terminal.FakeUI
		cmd                      *objectstorage.CredentialListCommand
		cliCommand               cli.Command
		fakeSession              *session.Session
		fakeObjectStorageManager managers.ObjectStorageManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeObjectStorageManager = managers.NewObjectStorageManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		cmd = objectstorage.NewCredentialListCommand(fakeUI, fakeObjectStorageManager)
		cliCommand = cli.Command{
			Name:        objectstorage.CredentialListMetaData().Name,
			Description: objectstorage.CredentialListMetaData().Description,
			Usage:       objectstorage.CredentialListMetaData().Usage,
			Flags:       objectstorage.CredentialListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Object Storage endpoints", func() {
		Context("Object Storage endpoints, Invalid Usage", func() {
			It("Set command without id", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument"))
			})
			It("Set command with id like letters", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Storage ID'. It must be a positive integer."))
			})
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Object Storage endpoints, correct use", func() {
			It("return objectstorage endpoints", func() {
				err := testhelpers.RunCommand(cliCommand, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Id         Password                       Username                 Type Name"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("17987654   abcdefghijklmnopqrstuvwxyz     123456mnopqrstuvwxyz     S3 Compatible Signature"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("19987654   aaabcdefghijklmnopqrstuvwxyz   11123456mnopqrstuvwxyz   S3 Compatible Signature"))
				
			})
			It("return objectstorage endpoints in format json", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Id": "17987654",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Password": "abcdefghijklmnopqrstuvwxyz",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Username": "123456mnopqrstuvwxyz",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Type Name": "S3 Compatible Signature"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Id": "19987654",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Password": "aaabcdefghijklmnopqrstuvwxyz",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Username": "11123456mnopqrstuvwxyz",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Type Name": "S3 Compatible Signature"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})

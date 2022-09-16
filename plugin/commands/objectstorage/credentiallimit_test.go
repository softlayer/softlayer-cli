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

var _ = Describe("Object Storage limit Object Storages", func() {
	var (
		fakeUI                   *terminal.FakeUI
		cmd                      *objectstorage.CredentialLimitCommand
		cliCommand               cli.Command
		fakeSession              *session.Session
		fakeObjectStorageManager managers.ObjectStorageManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeObjectStorageManager = managers.NewObjectStorageManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		cmd = objectstorage.NewCredentialLimitCommand(fakeUI, fakeObjectStorageManager)
		cliCommand = cli.Command{
			Name:        objectstorage.CredentialLimitMetaData().Name,
			Description: objectstorage.CredentialLimitMetaData().Description,
			Usage:       objectstorage.CredentialLimitMetaData().Usage,
			Flags:       objectstorage.CredentialLimitMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Object Storage credential limit", func() {
		Context("Object Storage credential limit, Invalid Usage", func() {
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

		Context("Object Storage credential limit, correct use", func() {
			It("return objectstorage credential limit", func() {
				err := testhelpers.RunCommand(cliCommand, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Limit"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2"))
			})
			It("return objectstorage credential limit in format json", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Limit": "2"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})

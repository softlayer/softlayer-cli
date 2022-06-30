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
		cmd                      *objectstorage.AccountsCommand
		cliCommand               cli.Command
		fakeSession              *session.Session
		fakeObjectStorageManager managers.ObjectStorageManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeObjectStorageManager = managers.NewObjectStorageManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		cmd = objectstorage.NewAccountsCommand(fakeUI, fakeObjectStorageManager)
		cliCommand = cli.Command{
			Name:        objectstorage.AccountsMetaData().Name,
			Description: objectstorage.AccountsMetaData().Description,
			Usage:       objectstorage.AccountsMetaData().Usage,
			Flags:       objectstorage.AccountsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Object Storage accounts", func() {
		Context("Object Storage accounts, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCommand(cliCommand, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Object Storage accounts, correct use", func() {
			It("return objectstorage accounts", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Id          Name     ApiType"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456789   SLUSER   S3"))
			})
			It("return objectstorage accounts in format json", func() {
				err := testhelpers.RunCommand(cliCommand, "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Id": "123456789",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "SLUSER",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"ApiType": "S3"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})

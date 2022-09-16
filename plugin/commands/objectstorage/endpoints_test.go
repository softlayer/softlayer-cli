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
		cmd                      *objectstorage.EndpointsCommand
		cliCommand               cli.Command
		fakeSession              *session.Session
		fakeObjectStorageManager managers.ObjectStorageManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeObjectStorageManager = managers.NewObjectStorageManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		cmd = objectstorage.NewEndpointsCommand(fakeUI, fakeObjectStorageManager)
		cliCommand = cli.Command{
			Name:        objectstorage.EndpointsMetaData().Name,
			Description: objectstorage.EndpointsMetaData().Description,
			Usage:       objectstorage.EndpointsMetaData().Usage,
			Flags:       objectstorage.EndpointsMetaData().Flags,
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
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Invoice ID'. It must be a positive integer."))
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
				Expect(fakeUI.Outputs()).To(ContainSubstring("Location/Region   Url                                                        EndPoint Type   Public/Private   Legacy"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("us-geo            s3.us.cloud-object-storage.appdomain.cloud                 Cross Region    Public           False"))
			})
			It("return objectstorage endpoints in format json", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Location/Region": "us-geo",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Url": "s3.us.cloud-object-storage.appdomain.cloud",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"EndPoint Type": "Cross Region",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Public/Private": "Public",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Legacy": "False"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})

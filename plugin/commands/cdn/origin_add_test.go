package cdn_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/cdn"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Cdn origin add", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *cdn.OriginAddCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = cdn.NewOriginAddCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("Cdn origin add", func() {
		Context("Cdn origin add, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
			It("Set command without id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage : This command requires one argument"))
			})
			It("Set command without flag hostname", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456789")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`required flag(s) "origin", "path" not set`))
			})
			It("Set command without flag origin", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456789", "--origin", "123.123.123.123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`required flag(s) "path" not set`))
			})
			It("Set command without flag http or https", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456789", "--origin", "123.123.123.123", "--path", "/example/videos")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`Incorrect Usage: 'http or https' is required`))
			})
			It("Set command with flag wrong origin-type", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456789", "--origin", "123.123.123.123", "--path", "/example/videos", "--http", "80", "--origin-type", "asdfgh")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`Incorrect Usage: --origintype`))
			})
			It("Set command with flag wrong optimize", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456789", "--origin", "123.123.123.123", "--path", "/example/videos", "--http", "80", "--optimize", "notPermit")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`Incorrect Usage: --optimize`))
			})
			It("Set command with flag wrong cache-key", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456789", "--origin", "123.123.123.123", "--path", "/example/videos", "--http", "80", "--cache-key", "notPermit")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`Incorrect Usage: --cache-key`))
			})
			It("Set command with flag originType like storage and empty bucket-name", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456789", "--origin", "123.123.123.123", "--path", "/example/videos", "--http", "80", "--origin-type", "storage")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`Incorrect Usage: --bucket-name can not be empty`))
			})
		})

		Context("Cdn origin add, correct use", func() {
			It("return cdn origin add", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456789", "--origin", "123.123.123.123", "--path", "/example/videos/", "--http", "80", "--origin-type", "storage", "--bucket-name", "bucketName", "--file-extensions", "jpg")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("CDN Unique ID"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("354034879028850"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("File Extension"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("jpg,pdf,jpeg,png"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Header"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("header.test.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Path"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("/example/videos"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Cache Key Rule"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("include-all"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Performance Configuration"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("General web delivery"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Status"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("RUNNING"))
			})
			It("return cdn in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456789", "--origin", "123.123.123.123", "--path", "/example/videos/", "--http", "80", "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "CDN Unique ID",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "354034879028850"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})

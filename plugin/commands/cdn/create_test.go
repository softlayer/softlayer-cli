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

var _ = Describe("Cdn create", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *cdn.CreateCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = cdn.NewCreateCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("Cdn create", func() {
		Context("Cdn create, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
			It("Set command without flag hostname", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`required flag(s) "hostname", "origin" not set`))
			})
			It("Set command without flag origin", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "www.example.com")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`required flag(s) "origin" not set`))
			})
			It("Set command without flag http or https", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "www.example.com", "--origin", "123.45.67.8")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`Incorrect Usage: 'http or https' is required`))
			})
			It("Set command with flag wrong origin-type", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "www.example.com", "--origin", "123.45.67.8", "--http", "80", "--origin-type", "asdfgh")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`Incorrect Usage: --origintype`))
			})
			It("Set command with flag wrong ssl", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "www.example.com", "--origin", "123.45.67.8", "--http", "80", "--ssl", "asdfgh")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`Incorrect Usage: --ssl`))
			})
		})

		Context("Cdn create, correct use", func() {
			It("return cdn created", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "www.example.com", "--origin", "123.45.67.8", "--http", "80")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("CDN Unique ID"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("354034879028850"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Bucket Name"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("test-bucket-name"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hostname"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("test.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Header"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("header.test.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("IBM CNAME"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("test.cdn.appdomain.cloud"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Protocol"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("HOST_SERVER"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Certificate Type"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("WILDCARD_CERT"))
			})
			It("return cdn in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "www.example.com", "--origin", "123.45.67.8", "--http", "80", "--output", "json")
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

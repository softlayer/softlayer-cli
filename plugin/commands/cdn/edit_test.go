package cdn_test

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/cdn"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Cdn edit", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *cdn.EditCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = cdn.NewEditCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("Cdn edit", func() {
		Context("Cdn edit, Invalid Usage", func() {
			It("Set command without id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
			It("Set command with id like letters", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'cdn ID'. It must be a positive integer."))
			})
			It("Set whitout any flag", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Please pass at least one of the flags."))
			})
			It("Set invalid flag respect-headers", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--respect-headers", "2")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Option respect-headers just accept '0' or '1'"))
			})
			It("Set invalid flag cache", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--cache", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Option cache just accept: 'include-all' 'ignore-all' 'include-specified' 'ignore-specified'"))
			})
			It("Set cache whitout cache-description", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--cache", "include-specified")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: cache-description option must be used"))
			})
			It("Set cache-description whitout cache", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--header", "New header", "--cache-description", "New cache description")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: cache-description is only used with the cache option"))
			})
			It("Set invalid flag performance-configuration", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--performance-configuration", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Option performance-configuration just accept: 'General web delivery' 'Large file optimization' 'Video on demand optimization'"))
			})
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--header", "New Header", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Cdn  edit, correct use", func() {
			It("return cdn  edit", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--header", "New Header", "--http-port", "81", "--origin", "12.12.12.12", "--respect-headers", "1", "--cache", "include-specified", "--cache-description", "New cache description", "--performance-configuration", "General web delivery")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("2020-10-07T13:13:41Z"))
				Expect(results[2]).To(ContainSubstring("www.techsupport.com"))
				Expect(results[3]).To(ContainSubstring("80"))
				Expect(results[4]).To(ContainSubstring("HOST_SERVER"))
				Expect(results[5]).To(ContainSubstring("General web delivery"))
				Expect(results[6]).To(ContainSubstring("HTTP"))
			})
			It("return cdn  edit in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--header", "New Header", "--http-port", "81", "--origin", "12.12.12.12", "--respect-headers", "1", "--cache", "include-all", "--performance-configuration", "Large file optimization", "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Header",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "www.techsupport.com"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Performance Configuration",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "General web delivery"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Cname",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "cdnakaog136c8gfq12000010.cdn.appdomain.cloud"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})

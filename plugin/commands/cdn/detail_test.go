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

var _ = Describe("Cdn detail", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *cdn.DetailCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = cdn.NewDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("Cdn detail", func() {
		Context("Cdn detail, Invalid Usage", func() {
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
			It("Set invalid flag history", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--history", "100")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: history"))
			})
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Cdn  detail, correct use", func() {
			It("return cdn  detail", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("172498049151824"))
				Expect(results[2]).To(ContainSubstring("www.testgo.com"))
				Expect(results[3]).To(ContainSubstring("HTTPS"))
				Expect(results[9]).To(ContainSubstring("2.0 GB"))
				Expect(results[10]).To(ContainSubstring("3"))
				Expect(results[11]).To(ContainSubstring("1.7 %"))

			})
			It("return cdn  detail in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Unique id",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "172498049151824"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Total bandwidth",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "2.0 GB"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Hit Radio",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "1.7 %"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})

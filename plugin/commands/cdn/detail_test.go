package cdn_test

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/cdn"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Cdn list Detail", func() {
	var (
		fakeUI         *terminal.FakeUI
		cmd            *cdn.DetailCommand
		cliCommand     cli.Command
		fakeSession    *session.Session
		fakeCdnManager managers.CdnManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeCdnManager = managers.NewCdnManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		cmd = cdn.NewDetailCommand(fakeUI, fakeCdnManager)
		cliCommand = cli.Command{
			Name:        cdn.DetailMetaData().Name,
			Description: cdn.DetailMetaData().Description,
			Usage:       cdn.DetailMetaData().Usage,
			Flags:       cdn.DetailMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Cdn detail", func() {
		Context("Cdn detail, Invalid Usage", func() {
			It("Set command without id", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
			It("Set command with id like letters", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'cdn ID'. It must be a positive integer."))
			})
			It("Set invalid flag history", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "--history", "100")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: history"))
			})
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Cdn  detail, correct use", func() {
			It("return cdn  detail", func() {
				err := testhelpers.RunCommand(cliCommand, "123")
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
				err := testhelpers.RunCommand(cliCommand, "123", "--output", "json")
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

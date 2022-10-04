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

var _ = Describe("Cdn list Cdn", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *cdn.ListCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = cdn.NewListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("Cdn list", func() {
		Context("Cdn list, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Cdn list, correct use", func() {
			It("return cdn list", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Unique Id      Domain          Origin         Vendor   Cname                       Status"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("321654987123   test.com        10.32.12.125   akamai   test.cdn.appdomain.cloud    CNAME_CONFIGURATION"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("654321789321   www.test2.com   10.32.12.125   akamai   test2.cdn.appdomain.cloud   ERROR"))
			})
			It("return cdn cdn in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Unique Id": "321654987123",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Domain": "www.test2.com",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Origin": "10.32.12.125",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Status": "CNAME_CONFIGURATION"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})

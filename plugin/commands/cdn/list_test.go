package cdn_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/cdn"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Cdn list Cdn", func() {
	var (
		fakeUI         *terminal.FakeUI
		cmd            *cdn.ListCommand
		cliCommand     cli.Command
		fakeSession    *session.Session
		fakeCdnManager managers.CdnManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeCdnManager = managers.NewCdnManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		cmd = cdn.NewListCommand(fakeUI, fakeCdnManager)
		cliCommand = cli.Command{
			Name:        cdn.ListMetaData().Name,
			Description: cdn.ListMetaData().Description,
			Usage:       cdn.ListMetaData().Usage,
			Flags:       cdn.ListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Cdn list", func() {
		Context("Cdn list, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCommand(cliCommand, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Cdn list, correct use", func() {
			It("return cdn list", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Unique Id      Domain          Origin         Vendor   Cname                       Status"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("321654987123   test.com        10.32.12.125   akamai   test.cdn.appdomain.cloud    CNAME_CONFIGURATION"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("654321789321   www.test2.com   10.32.12.125   akamai   test2.cdn.appdomain.cloud   ERROR"))
			})
			It("return cdn cdn in format json", func() {
				err := testhelpers.RunCommand(cliCommand, "--output", "json")
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

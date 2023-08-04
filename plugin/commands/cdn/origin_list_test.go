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

var _ = Describe("Cdn origin list Cdn", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *cdn.OriginListCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = cdn.NewOriginListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("Cdn origin list", func() {
		Context("Cdn origin list, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456789", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})

			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage : This command requires one argument."))
			})
		})

		Context("Cdn origin list, correct use", func() {
			It("return cdn list", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456789")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Path"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Origin"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("/example1/*"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123.123.123.123"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("RUNNING"))
			})

			It("return origin list in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456789", "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Path": "/example1/*",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Origin": "123.123.123.123",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Http Port": "80",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Status": "RUNNING"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})

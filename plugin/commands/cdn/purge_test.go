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

var _ = Describe("Cdn purge", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *cdn.PurgeCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = cdn.NewPurgeCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("Cdn purge", func() {
		Context("Cdn purge, Invalid Usage", func() {
			It("Set command without id and path", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires two arguments."))
			})
		})

		Context("Cdn purge, correct use", func() {
			It("return cdn purge", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456789", "/example/")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Date"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Path"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Saved"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Status"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2023-04-28 12:06:27"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("/example/"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("UNSAVED"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SUCCESS"))
			})
			It("return cdn purge in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456789", "/example/", "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Date": "2023-04-28 12:06:27",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Path": "/example/",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Saved": "UNSAVED",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Status": "SUCCESS"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})

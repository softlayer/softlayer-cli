package account_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Account list EventDetail", func() {
	var (
		fakeUI				*terminal.FakeUI
		cliCommand			*account.EventDetailCommand
		fakeSession   		*session.Session
		slCommand			*metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand  = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = account.NewEventDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("Account events", func() {
		Context("Account events, Invalid Usage", func() {
			It("Set command without id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
			It("Set command with id like letters", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Event ID'. It must be a positive integer."))
			})
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Account events, correct use", func() {
			It("return account events", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("ACTION REQUIRED - Windows"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Id       Status      Type           Start   End"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("340846   Published   ANNOUNCEMENT   -       -"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Id          Hostname                     Label"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("112238162   resource.softlayer.test      Capacity - Windows"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("121450334   resource2.softlayer2.test2   Capacity - Windows"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("======= Update #1 on 2022-03-23T00:50:57Z ======="))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Updated"))
			})
			It("return account events in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`ACTION REQUIRED - Windows`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Id": "340846",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Status": "Published",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Type": "ANNOUNCEMENT",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Updates": "======= Update #1 on 2022-03-23T00:50:57Z ======="`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Updates": "Updated message"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})

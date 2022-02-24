package ticket_test

import (
	"errors"
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ticket"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("ticket create", func() {
	var (
		fakeUI            *terminal.FakeUI
		fakeTicketManager *testhelpers.FakeTicketManager
		cmd               *ticket.CreateStandardTicketCommand
		cliCommand        cli.Command
		ticket_id         int
	)
	ticket_id = 12345
	BeforeEach(func() {
		fakeTicketManager = new(testhelpers.FakeTicketManager)
		fakeUI = terminal.NewFakeUI()
		cmd = ticket.NewCreateStandardTicketCommand(fakeUI, fakeTicketManager)
		cliCommand = cli.Command{
			Name:        ticket.TicketCreateMetaData().Name,
			Description: ticket.TicketCreateMetaData().Description,
			Usage:       ticket.TicketCreateMetaData().Usage,
			Flags:       ticket.TicketCreateMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Ticket create", func() {
		Context("ticket create", func() {
			It("Success: inline body", func() {
				fakeTicketManager.CreateStandardTicketReturns(&ticket_id, nil)
				err := testhelpers.RunCommand(cliCommand, "--subject-id=1522", `--title="Test Title"`, `--body="Test Contents"`)
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Ticket ID: 12345"}))
			})

			It("Success: body from TicketManager.GetText()", func() {
				fakeTicketManager.GetTextReturns("Body goes here", nil)
				fakeTicketManager.CreateStandardTicketReturns(&ticket_id, nil)
				err := testhelpers.RunCommand(cliCommand, "--subject-id=1522", `--title="Test Title"`)
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Ticket ID: 12345"}))
				Expect(err).ToNot(HaveOccurred())
			})

			It("Success with all options", func() {
				fakeTicketManager.CreateStandardTicketReturns(&ticket_id, nil)
				err := testhelpers.RunCommand(cliCommand, "--subject-id=1522", `--title="Test Title"`, "--attachment=111111", `--rootpwd="thisisapassword"`)
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Ticket ID: 12345"}))
				Expect(err).ToNot(HaveOccurred())
			})

			It("Failure: No Subjet ID", func() {
				err := testhelpers.RunCommand(cliCommand, `--title="Test Title"`, "--attachment=111111", `--rootpwd="thisisapassword"`)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstrings([]string{"requires the --subject-id"}))
			})

			It("Failure: No Title", func() {
				err := testhelpers.RunCommand(cliCommand, "--subject-id=1522", "--attachment=111111", `--rootpwd="thisisapassword"`)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstrings([]string{"requires the --title"}))
			})

			It("Failure: No body", func() {
				fakeTicketManager.GetTextReturns("nil", errors.New("Body Error"))
				err := testhelpers.RunCommand(cliCommand, "--subject-id=1522", `--title="Test Title"`, "--attachment=111111", `--rootpwd="thisisapassword"`)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstrings([]string{"Body Error"}))
			})

			It("Failure: No root password", func() {
				err := testhelpers.RunCommand(cliCommand, "--subject-id=1522", `--title="Test Title"`, "--attachment=111111")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstrings([]string{"Root password must be provided"}))
			})

			It("Failure: API failre", func() {
				fakeTicketManager.CreateStandardTicketReturns(nil, errors.New("API Error"))
				err := testhelpers.RunCommand(cliCommand, "--subject-id=1522", `--title="Test Title"`, "--attachment=111111", `--rootpwd="thisisapassword"`)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstrings([]string{"API Error"}))

			})
		})
	})
})

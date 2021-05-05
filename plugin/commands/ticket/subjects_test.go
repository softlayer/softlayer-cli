package ticket_test

import (
	"errors"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ticket"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("ticket subjects", func() {
	var (
		fakeUI            *terminal.FakeUI
		fakeTicketManager *testhelpers.FakeTicketManager
		cmd               *ticket.SubjectsTicketCommand
		cliCommand        cli.Command
	)
	fakeTicketManager = new(testhelpers.FakeTicketManager)
	session := testhelpers.NewFakeSoftlayerSession(nil)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		cmd = ticket.NewSubjectsTicketCommand(fakeUI, fakeTicketManager)
		cliCommand = cli.Command{
			Name:        metadata.TicketSubjectsMetaData().Name,
			Description: metadata.TicketSubjectsMetaData().Description,
			Usage:       metadata.TicketSubjectsMetaData().Usage,
			Flags:       metadata.TicketSubjectsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Ticket subjects", func() {
		Context("ticket subjects", func() {
			It("Normal command call", func() {
				var returnData []datatypes.Ticket_Subject
				// This just loads data from the fixtures JSON file.
				_ = session.DoRequest("SoftLayer_Ticket_Subjet", "getAllObjects", nil, nil, &returnData)
				fakeTicketManager.GetSubjectsReturns(&returnData, nil)
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).ToNot(HaveOccurred())
			})

			It("return fail 1", func() {
				err := testhelpers.RunCommand(cliCommand, "argument")
				Expect(err).To(HaveOccurred())
			})

			It("API Failure", func() {
				fakeTicketManager.GetSubjectsReturns(nil, errors.New("Internal Server Error"))
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

package ticket_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ticket"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ticket Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test ticket commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		ticketCommands := ticket.SetupCobraCommands(slMeta)
		Expect(ticketCommands.Name()).To(Equal("ticket"))
	})
	Context("Ticket Namespace", func() {
		It("Ticket Name Space", func() {
			Expect(ticket.TicketNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(ticket.TicketNamespace().Name).To(ContainSubstring("ticket"))
			Expect(ticket.TicketNamespace().Description).To(ContainSubstring("Classic infrastructure Manage Tickets"))
		})
	})
})

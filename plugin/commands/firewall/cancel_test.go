package firewall_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/firewall"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("firewall cancel", func() {
	var (
		fakeUI              *terminal.FakeUI
		cliCommand          *firewall.CancelCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
		fakeFirewallManager managers.FirewallManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeFirewallManager = managers.NewFirewallManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = firewall.NewCancelCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.FirewallManager = fakeFirewallManager
	})

	Describe("firewall cancel", func() {

		Context("Firewall cancel, Invalid Usage", func() {
			It("Set without ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument"))
			})
			It("Set invalid ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--force")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid ID 123456: ID should be of the form xxx:yyy, xxx is the type of the firewall, yyy is the positive integer ID."))
			})
			It("Set invalid flag force", func() {
				fakeUI.Inputs("abc")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "vs:123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
			})
			It("Set No in flag force", func() {
				fakeUI.Inputs("no")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "vs:123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted"))
			})
		})

		Context("Return no error", func() {
			It("firewall vs canceled", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "vs:123", "--force")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("vs:123 is being cancelled!"))
			})
			It("firewall vlan canceled", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "vlan:123", "--force")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("vlan:123 is being cancelled!"))
			})
		})
	})
})

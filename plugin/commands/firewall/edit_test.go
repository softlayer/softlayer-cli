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

var _ = Describe("firewall edit", func() {
	var (
		fakeUI              *terminal.FakeUI
		cliCommand          *firewall.EditCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
		fakeFirewallManager managers.FirewallManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeFirewallManager = managers.NewFirewallManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = firewall.NewEditCommand(slCommand)
		cliCommand.FirewallManager = fakeFirewallManager
	})

	Describe("firewall edit", func() {
		Context("Firewall edit, Invalid Usage", func() {
			It("Set without ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument"))
			})
			It("Set invalid ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid ID 123456: ID should be of the form xxx:yyy, xxx is the type of the firewall, yyy is the positive integer ID."))
			})
			It("Set multivlan", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "multiVlan:123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("All multi vlan rules must be managed through the FortiGate dashboard using the provided credentials."))
			})
		})
	})
})

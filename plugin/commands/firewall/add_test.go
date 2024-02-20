package firewall_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/firewall"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("firewall add", func() {
	var (
		fakeUI              *terminal.FakeUI
		cliCommand          *firewall.AddCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
		fakeFirewallManager managers.FirewallManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeFirewallManager = managers.NewFirewallManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = firewall.NewAddCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.FirewallManager = fakeFirewallManager
	})

	Describe("firewall add", func() {

		Context("Firewall add, Invalid Usage", func() {
			It("Set without ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument"))
			})
			It("Set without type", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--type' is required"))
			})
			It("Set command with id like letters", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc", "--type", "vlan")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Target ID'. It must be a positive integer."))
			})
			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--type", "vlan", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
			It("Set invalid flag force", func() {
				fakeUI.Inputs("abc")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--type", "vlan")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
			})
			It("Set No in flag force", func() {
				fakeUI.Inputs("no")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--type", "vlan")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted"))
			})
		})

		Context("Return no error", func() {
			It("Firewall type vlan created", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--type", "vlan", "--force")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Product: CDN 25 GB Storage"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 11493593 was placed to create a firewall."))
			})
			It("Firewall type vs created", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--type", "vs", "--force")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Product: CDN 25 GB Storage"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 11493593 was placed to create a firewall."))
			})
			It("Firewall type hardware created", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--type", "hardware", "--force")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Product: CDN 25 GB Storage"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 11493593 was placed to create a firewall."))
			})
		})
	})
})

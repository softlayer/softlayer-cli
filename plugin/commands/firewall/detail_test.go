package firewall_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/firewall"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("firewall detail", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeFirewallManager *testhelpers.FakeFirewallManager
		cmd                 *firewall.DetailCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeFirewallManager = new(testhelpers.FakeFirewallManager)
		cmd = firewall.NewDetailCommand(fakeUI, fakeFirewallManager)
		cliCommand = cli.Command{
			Name:        firewall.FirewallDetailMetaData().Name,
			Description: firewall.FirewallDetailMetaData().Description,
			Usage:       firewall.FirewallDetailMetaData().Usage,
			Flags:       firewall.FirewallDetailMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("firewall detail", func() {

		Context("Return error", func() {

			It("Set without ID", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument."))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCommand(cliCommand, "vs:123456", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeFirewallManager.ParseFirewallIDReturns("", 0, errors.New("Failed to parse firewall ID"))
			})

			It("Set invalid ID", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to parse firewall ID"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeFirewallManager.ParseFirewallIDReturns("vlan", 123456, nil)
				fakeFirewallManager.GetDedicatedFirewallRulesReturns([]datatypes.Network_Vlan_Firewall_Rule{}, errors.New("Failed to get dedicated firewall rules."))
			})
			It("Failed get vlan firewall", func() {
				err := testhelpers.RunCommand(cliCommand, "vlan:123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get dedicated firewall rules."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeFirewallManager.ParseFirewallIDReturns("vs", 123456, nil)
				fakeFirewallManager.GetStandardFirewallRulesReturns([]datatypes.Network_Component_Firewall_Rule{}, errors.New("Failed to get standard firewall rules."))
			})
			It("Failed get standard firewall", func() {
				err := testhelpers.RunCommand(cliCommand, "vs:123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get standard firewall rules."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeFirewallManager.ParseFirewallIDReturns("vlan", 123456, nil)
				fakerRules := []datatypes.Network_Vlan_Firewall_Rule{
					datatypes.Network_Vlan_Firewall_Rule{
						OrderValue:                sl.Int(1),
						Action:                    sl.String("permit"),
						Protocol:                  sl.String("tcp"),
						SourceIpAddress:           sl.String("0.0.0.0"),
						SourceIpSubnetMask:        sl.String("0.0.0.0"),
						DestinationIpAddress:      sl.String("any on server"),
						DestinationPortRangeStart: sl.Int(85),
						DestinationPortRangeEnd:   sl.Int(85),
						DestinationIpSubnetMask:   sl.String("255.255.255.255"),
					},
				}
				fakeFirewallManager.GetDedicatedFirewallRulesReturns(fakerRules, nil)
			})

			It("get dedicated firewalls rules", func() {
				err := testhelpers.RunCommand(cliCommand, "vlan:123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("permit"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("tcp"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.0.0.0"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.0.0.0"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("any on server:85-85"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("255.255.255.255"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeFirewallManager.ParseFirewallIDReturns("vs", 123456, nil)
				fakerRules := []datatypes.Network_Component_Firewall_Rule{
					datatypes.Network_Component_Firewall_Rule{
						OrderValue:                sl.Int(1),
						Action:                    sl.String("permit"),
						Protocol:                  sl.String("tcp"),
						SourceIpAddress:           sl.String("0.0.0.0"),
						SourceIpSubnetMask:        sl.String("0.0.0.0"),
						DestinationIpAddress:      sl.String("any on server"),
						DestinationPortRangeStart: sl.Int(85),
						DestinationPortRangeEnd:   sl.Int(85),
						DestinationIpSubnetMask:   sl.String("255.255.255.255"),
					},
				}
				fakeFirewallManager.GetStandardFirewallRulesReturns(fakerRules, nil)
			})

			It("get stadard firewalls rules", func() {
				err := testhelpers.RunCommand(cliCommand, "vs:123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("permit"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("tcp"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.0.0.0"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.0.0.0"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("any on server:85-85"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("255.255.255.255"))
			})
		})
	})
})

package firewall_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/firewall"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("firewall detail", func() {
	var (
		fakeUI              *terminal.FakeUI
		cliCommand          *firewall.DetailCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
		FakeFirewallManager *testhelpers.FakeFirewallManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = firewall.NewDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		FakeFirewallManager = new(testhelpers.FakeFirewallManager)
		cliCommand.FirewallManager = FakeFirewallManager
	})

	Describe("firewall detail", func() {

		Context("Return error", func() {

			It("Set without ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument"))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "vs:123456", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				FakeFirewallManager.ParseFirewallIDReturns("", 0, errors.New("Failed to parse firewall ID"))
			})

			It("Set invalid ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to parse firewall ID"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				FakeFirewallManager.ParseFirewallIDReturns("vlan", 123456, nil)
				FakeFirewallManager.GetDedicatedFirewallRulesReturns([]datatypes.Network_Vlan_Firewall_Rule{}, errors.New("Failed to get dedicated firewall rules.\n"))
			})
			It("Failed get vlan firewall", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "vlan:123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get dedicated firewall rules."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				FakeFirewallManager.ParseFirewallIDReturns("vs", 123456, nil)
				FakeFirewallManager.GetStandardFirewallRulesReturns([]datatypes.Network_Component_Firewall_Rule{}, errors.New("Failed to get standard firewall rules."))
			})
			It("Failed get standard firewall", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "vs:123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get standard firewall rules."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				FakeFirewallManager.ParseFirewallIDReturns("multiVlan", 123456, nil)
				FakeFirewallManager.GetMultiVlanFirewallReturns(datatypes.Network_Vlan_Firewall{}, errors.New("Failed to get multi vlan firewall."))
			})
			It("Failed get standard firewall", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "multiVlan:123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get multi vlan firewall."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				FakeFirewallManager.ParseFirewallIDReturns("vlan", 123456, nil)
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
				FakeFirewallManager.GetDedicatedFirewallRulesReturns(fakerRules, nil)
			})

			It("get dedicated firewalls rules", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "vlan:123456")
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
				FakeFirewallManager.ParseFirewallIDReturns("vs", 123456, nil)
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
				FakeFirewallManager.GetStandardFirewallRulesReturns(fakerRules, nil)
			})

			It("get standard firewalls rules", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "vs:123456")
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
				FakeFirewallManager.ParseFirewallIDReturns("multiVlan", 123456, nil)
				fakerMultiVlan := datatypes.Network_Vlan_Firewall{
					NetworkGateway: &datatypes.Network_Gateway{
						Name: sl.String("firewall1"),
						PublicIpAddress: &datatypes.Network_Subnet_IpAddress{
							IpAddress: sl.String("1.1.1.1"),
						},
						PrivateIpAddress: &datatypes.Network_Subnet_IpAddress{
							IpAddress: sl.String("192.168.1.2"),
						},
						PublicIpv6Address: &datatypes.Network_Subnet_IpAddress{
							IpAddress: sl.String("2607:f0d0:1704:0020:0000:0000:0000:0002"),
						},
						PublicVlan: &datatypes.Network_Vlan{
							VlanNumber: sl.Int(1111),
						},
						PrivateVlan: &datatypes.Network_Vlan{
							VlanNumber: sl.Int(2222),
						},
					},
					Datacenter: &datatypes.Location{
						LongName: sl.String("Dallas 13"),
					},
					FirewallType: sl.String("fortigate-security-appliance-10gb"),
					ManagementCredentials: &datatypes.Software_Component_Password{
						Username: sl.String("myUsername"),
						Password: sl.String("test1234."),
					},
					Rules: []datatypes.Network_Vlan_Firewall_Rule{
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
					},
				}
				FakeFirewallManager.GetMultiVlanFirewallReturns(fakerMultiVlan, nil)
			})

			It("get multi vlan with credentials", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "multiVlan:123456", "--credentials")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("firewall1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1.1.1.1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("192.168.1.2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2607:f0d0:1704:0020:0000:0000:0000:0002"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2222"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Dallas 13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("fortigate-security-appliance-10gb"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("myUsername"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("test1234."))
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

package firewall_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/firewall"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("firewall list", func() {
	var (
		fakeUI              *terminal.FakeUI
		cliCommand          *firewall.ListCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
		FakeFirewallManager *testhelpers.FakeFirewallManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = firewall.NewListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		FakeFirewallManager = new(testhelpers.FakeFirewallManager)
		cliCommand.FirewallManager = FakeFirewallManager
	})

	Describe("firewall list", func() {

		Context("Return error", func() {

			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				FakeFirewallManager.GetFirewallsReturns([]datatypes.Network_Vlan{}, errors.New("Failed to get firewalls on your account"))
			})
			It("Failed get firewalls", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get firewalls on your account"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerVlans := []datatypes.Network_Vlan{
					datatypes.Network_Vlan{
						Id:                           sl.Int(222222),
						HighAvailabilityFirewallFlag: sl.Bool(true),
						NetworkVlanFirewall: &datatypes.Network_Vlan_Firewall{
							Id: sl.Int(111111),
						},
						DedicatedFirewallFlag: sl.Int(1),
					},
					datatypes.Network_Vlan{
						FirewallGuestNetworkComponents: []datatypes.Network_Component_Firewall{
							datatypes.Network_Component_Firewall{
								Id:     sl.Int(333333),
								Status: sl.String("allow_edit"),
								GuestNetworkComponent: &datatypes.Virtual_Guest_Network_Component{
									GuestId: sl.Int(444444),
								},
							},
						},
					},
					datatypes.Network_Vlan{
						FirewallNetworkComponents: []datatypes.Network_Component_Firewall{
							datatypes.Network_Component_Firewall{
								Id:     sl.Int(555555),
								Status: sl.String("allow_edit"),
								NetworkComponent: &datatypes.Network_Component{
									DownlinkComponent: &datatypes.Network_Component{
										HardwareId: sl.Int(666666),
									},
								},
							},
						},
					},
				}
				fakerMultiVlans := []datatypes.Network_Gateway{
					datatypes.Network_Gateway{
						InsideVlans: []datatypes.Network_Gateway_Vlan{
							datatypes.Network_Gateway_Vlan{
								NetworkVlanId: sl.Int(888888),
							},
						},
						NetworkFirewall: &datatypes.Network_Vlan_Firewall{
							Id:           sl.Int(777777),
							FirewallType: sl.String("fortigate-security-appliance-10gb"),
							Datacenter: &datatypes.Location{
								Name: sl.String("dal13"),
							},
						},
						Name: sl.String("testfirewall"),
						Members: []datatypes.Network_Gateway_Member{
							datatypes.Network_Gateway_Member{
								Hardware: &datatypes.Hardware{
									Hostname: sl.String("dft03.pod03.dal13"),
								},
							},
						},
						PublicIpAddress: &datatypes.Network_Subnet_IpAddress{
							IpAddress: sl.String("65.65.65.65"),
						},
						PrivateIpAddress: &datatypes.Network_Subnet_IpAddress{
							IpAddress: sl.String("10.2.2.2"),
						},
						Status: &datatypes.Network_Gateway_Status{
							KeyName: sl.String("ACTIVE"),
						},
					},
				}
				FakeFirewallManager.GetFirewallsReturns(fakerVlans, nil)
				FakeFirewallManager.GetMultiVlanFirewallsReturns(fakerMultiVlans, nil)
			})

			It("get firewalls", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("vlan:111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("VLAN - dedicated"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("HA"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("222222"))

				Expect(fakeUI.Outputs()).To(ContainSubstring("vs:333333"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Virtual Server - standard"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("-"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("444444"))

				Expect(fakeUI.Outputs()).To(ContainSubstring("server:555555"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hardware Server - standard"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("-"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("666666"))

				Expect(fakeUI.Outputs()).To(ContainSubstring("multiVlan:777777"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("testfirewall"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("fortigate-security-appliance-10gb"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dft03.pod03.dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("65.65.65.65"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("10.2.2.2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1 VLANs"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ACTIVE"))
			})
		})
	})
})

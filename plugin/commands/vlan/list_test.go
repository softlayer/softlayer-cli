package vlan_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/vlan"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VLAN List", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *vlan.ListCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeNetworkManager *testhelpers.FakeNetworkManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = vlan.NewListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cliCommand.NetworkManager = fakeNetworkManager
	})

	Describe("VLAN list", func() {
		Context("VLAN list but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListVlansReturns([]datatypes.Network_Vlan{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list VLANs on your account."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
			It("return error", func() {
				fakeNetworkManager.ListVlansReturns([]datatypes.Network_Vlan{}, nil)
				fakeNetworkManager.GetPodsReturns([]datatypes.Network_Pod{}, errors.New("Internal Server Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get Pods."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("VLAN list with wrong --sortby", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListVlansReturns([]datatypes.Network_Vlan{}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "abcd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --sortby abcd is not supported."))
			})
		})

		Context("VLAN list", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListVlansReturns([]datatypes.Network_Vlan{}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{""}))
			})
		})

		Context("VLAN list", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListVlansReturns([]datatypes.Network_Vlan{
					datatypes.Network_Vlan{
						Id:                 sl.Int(123456),
						VlanNumber:         sl.Int(100),
						FullyQualifiedName: sl.String("dal05.fcr01.784"),
						Name:               sl.String("Bill"),
						FirewallInterfaces: []datatypes.Network_Firewall_Module_Context_Interface{
							datatypes.Network_Firewall_Module_Context_Interface{
								Id: sl.Int(1),
							},
							datatypes.Network_Firewall_Module_Context_Interface{
								Id: sl.Int(2),
							},
						},
						PrimaryRouter: &datatypes.Hardware_Router{
							Hardware_Switch: datatypes.Hardware_Switch{
								Hardware: datatypes.Hardware{
									Id: sl.Int(987654),
									Datacenter: &datatypes.Location{
										Name: sl.String("dal05"),
									},
								},
							},
						},
						BillingItem: &datatypes.Billing_Item{
							Id: sl.Int(456321),
						},
						AttachedNetworkGateway: &datatypes.Network_Gateway{
							Name: sl.String("support"),
						},
						NetworkSpace:               sl.String("Public"),
						HardwareCount:              sl.Uint(uint(30)),
						VirtualGuestCount:          sl.Uint(uint(40)),
						TotalPrimaryIpAddressCount: sl.Uint(uint(50)),
						TagReferences: []datatypes.Tag_Reference{
							datatypes.Tag_Reference{
								Tag: &datatypes.Tag{
									Name: sl.String("Tag"),
								},
							},
						},
					},
					datatypes.Network_Vlan{
						Id:                 sl.Int(123458),
						VlanNumber:         sl.Int(73),
						FullyQualifiedName: sl.String("dal06.fcr01.797"),
						Name:               sl.String("Anne"),
						FirewallInterfaces: []datatypes.Network_Firewall_Module_Context_Interface{},
						PrimaryRouter: &datatypes.Hardware_Router{
							Hardware_Switch: datatypes.Hardware_Switch{
								Hardware: datatypes.Hardware{
									Id: sl.Int(87654),
									Datacenter: &datatypes.Location{
										Name: sl.String("dal06"),
									},
								},
							},
						},
						NetworkVlanFirewall: &datatypes.Network_Vlan_Firewall{
							FullyQualifiedDomainName: sl.String("firewall"),
						},
						NetworkSpace:               sl.String("Private"),
						HardwareCount:              sl.Uint(uint(5)),
						VirtualGuestCount:          sl.Uint(uint(6)),
						TotalPrimaryIpAddressCount: sl.Uint(uint(7)),
						TagReferences: []datatypes.Tag_Reference{
							datatypes.Tag_Reference{
								Tag: &datatypes.Tag{
									Name: sl.String("Tag"),
								},
							},
						},
					},
				}, nil)
				fakeNetworkManager.GetPodsReturns([]datatypes.Network_Pod{
					datatypes.Network_Pod{
						BackendRouterId:  sl.Int(987654),
						FrontendRouterId: sl.Int(123456),
						Capabilities:     []string{"CLOSURE_ANNOUNCED"},
						Name:             sl.String("dal05.pod01"),
					},
					datatypes.Network_Pod{
						BackendRouterId:  sl.Int(213456),
						FrontendRouterId: sl.Int(87654),
						Capabilities:     []string{},
						Name:             sl.String("dal06.pod02"),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("123456"))
				Expect(results[2]).To(ContainSubstring("123458"))
			})

			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "number")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("73"))
				Expect(results[2]).To(ContainSubstring("100"))
			})

			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "name")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("Anne"))
				Expect(results[2]).To(ContainSubstring("Bill"))
			})

			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "firewall")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("No"))
				Expect(results[2]).To(ContainSubstring("Yes"))
			})

			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "datacenter")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("dal05"))
				Expect(results[2]).To(ContainSubstring("dal06"))
			})

			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "hardware")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("5"))
				Expect(results[2]).To(ContainSubstring("30"))
			})

			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "virtual_servers")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("6"))
				Expect(results[2]).To(ContainSubstring("40"))
			})

			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "public_ips")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("7"))
				Expect(results[2]).To(ContainSubstring("50"))
			})

			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})

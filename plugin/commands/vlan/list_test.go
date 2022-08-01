package vlan_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/vlan"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VLAN List", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *vlan.ListCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = vlan.NewListCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        vlan.VlanListMetaData().Name,
			Description: vlan.VlanListMetaData().Description,
			Usage:       vlan.VlanListMetaData().Usage,
			Flags:       vlan.VlanListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VLAN list", func() {
		Context("VLAN list but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListVlansReturns([]datatypes.Network_Vlan{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to list VLANs on your account.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
			It("return error", func() {
				fakeNetworkManager.ListVlansReturns([]datatypes.Network_Vlan{}, nil)
				fakeNetworkManager.GetPodsReturns([]datatypes.Network_Pod{}, errors.New("Internal Server Error"))
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get Pods.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("VLAN list with wrong --sortby", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListVlansReturns([]datatypes.Network_Vlan{}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "abcd")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: --sortby abcd is not supported.")).To(BeTrue())
			})
		})

		Context("VLAN list", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListVlansReturns([]datatypes.Network_Vlan{}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
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
								datatypes.Hardware{
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
				err := testhelpers.RunCommand(cliCommand, "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "123456")).To(BeTrue())
				Expect(strings.Contains(results[2], "123458")).To(BeTrue())
			})

			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "number")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "73")).To(BeTrue())
				Expect(strings.Contains(results[2], "100")).To(BeTrue())
			})

			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "name")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "Anne")).To(BeTrue())
				Expect(strings.Contains(results[2], "Bill")).To(BeTrue())
			})

			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "firewall")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "No")).To(BeTrue())
				Expect(strings.Contains(results[2], "Yes")).To(BeTrue())
			})

			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "datacenter")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "dal05")).To(BeTrue())
				Expect(strings.Contains(results[2], "dal06")).To(BeTrue())
			})

			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "hardware")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "5")).To(BeTrue())
				Expect(strings.Contains(results[2], "30")).To(BeTrue())
			})

			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "virtual_servers")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "6")).To(BeTrue())
				Expect(strings.Contains(results[2], "40")).To(BeTrue())
			})

			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "public_ips")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "7")).To(BeTrue())
				Expect(strings.Contains(results[2], "50")).To(BeTrue())
			})

			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})

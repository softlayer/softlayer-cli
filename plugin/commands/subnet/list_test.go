package subnet_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/subnet"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Subnet list", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *subnet.ListCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeNetworkManager *testhelpers.FakeNetworkManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = subnet.NewListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cliCommand.NetworkManager = fakeNetworkManager
	})

	Describe("Subnet list", func() {
		Context("Subnet list with wrong --sortby", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --sortby abc is not supported."))
			})
		})

		Context("Subnet list but server fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListSubnetsReturns([]datatypes.Network_Subnet{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list subnets on your account."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Subnet list with different --sortby columns", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListSubnetsReturns([]datatypes.Network_Subnet{
					datatypes.Network_Subnet{
						Id:                sl.Int(321),
						NetworkIdentifier: sl.String("9.9.9.9"),
						SubnetType:        sl.String("SECONDARY"),
						NetworkVlan: &datatypes.Network_Vlan{
							Id:           sl.Int(7654),
							NetworkSpace: sl.String("PUBLIC"),
						},
						AddressSpace: sl.String("PUBLIC"),
						Datacenter: &datatypes.Location_Datacenter{
							Location: datatypes.Location{
								Name: sl.String("tok02"),
							},
						},
					},
					datatypes.Network_Subnet{
						Id:                sl.Int(123),
						NetworkIdentifier: sl.String("9.0.9.9"),
						SubnetType:        sl.String("PRIMARY"),
						NetworkVlan: &datatypes.Network_Vlan{
							Id:           sl.Int(4567),
							NetworkSpace: sl.String("PRIVATE"),
						},
						AddressSpace: sl.String("PRIVATE"),
						Datacenter: &datatypes.Location_Datacenter{
							Location: datatypes.Location{
								Name: sl.String("dal10"),
							},
						},
						IpAddresses: []datatypes.Network_Subnet_IpAddress{
							datatypes.Network_Subnet_IpAddress{
								Id: sl.Int(345),
							},
							datatypes.Network_Subnet_IpAddress{
								Id: sl.Int(456),
							},
						},
						Hardware: []datatypes.Hardware{
							datatypes.Hardware{
								Hostname:                sl.String("hw1"),
								Domain:                  sl.String("wilma.com"),
								PrimaryIpAddress:        sl.String("9.9.9.3"),
								PrimaryBackendIpAddress: sl.String("1.2.1.2"),
							},
						},
						VirtualGuests: []datatypes.Virtual_Guest{
							datatypes.Virtual_Guest{
								Hostname:                sl.String("vs1"),
								Domain:                  sl.String("wilma.com"),
								PrimaryIpAddress:        sl.String("9.9.9.2"),
								PrimaryBackendIpAddress: sl.String("1.2.1.1"),
							},
						},
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("123"))
				Expect(results[2]).To(ContainSubstring("321"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "identifier")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("9.0.9.9"))
				Expect(results[2]).To(ContainSubstring("9.9.9.9"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "type")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("Primary"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "network_space")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("Private"))
				Expect(results[2]).To(ContainSubstring("Public"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "IPs")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("0"))
				Expect(results[2]).To(ContainSubstring("2"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "hardware")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("0"))
				Expect(results[2]).To(ContainSubstring("1"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "vs")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("0"))
				Expect(results[2]).To(ContainSubstring("1"))
			})
		})
	})
})

package subnet_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/subnet"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Subnet list", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *subnet.ListCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = subnet.NewListCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        metadata.SubnetListMetaData().Name,
			Description: metadata.SubnetListMetaData().Description,
			Usage:       metadata.SubnetListMetaData().Usage,
			Flags:       metadata.SubnetListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Subnet list", func() {
		Context("Subnet list with wrong --sortby", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: --sortby abc is not supported.")).To(BeTrue())
			})
		})

		Context("Subnet list but server fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListSubnetsReturns([]datatypes.Network_Subnet{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to list subnets on your account.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
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
				err := testhelpers.RunCommand(cliCommand, "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "123")).To(BeTrue())
				Expect(strings.Contains(results[2], "321")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "identifier")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "9.0.9.9")).To(BeTrue())
				Expect(strings.Contains(results[2], "9.9.9.9")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "type")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "PRIMARY")).To(BeTrue())
				Expect(strings.Contains(results[2], "SECONDARY")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "network_space")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "PRIVATE")).To(BeTrue())
				Expect(strings.Contains(results[2], "PUBLIC")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "datacenter")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "4567")).To(BeTrue())
				Expect(strings.Contains(results[2], "7654")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "IPs")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "0")).To(BeTrue())
				Expect(strings.Contains(results[2], "2")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "hardware")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "0")).To(BeTrue())
				Expect(strings.Contains(results[2], "1")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "vs")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "0")).To(BeTrue())
				Expect(strings.Contains(results[2], "1")).To(BeTrue())
			})
		})
	})
})

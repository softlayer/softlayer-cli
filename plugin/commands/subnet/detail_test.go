package subnet_test

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
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/subnet"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Subnet detail", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *subnet.DetailCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = subnet.NewDetailCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        subnet.SubnetDetailMetaData().Name,
			Description: subnet.SubnetDetailMetaData().Description,
			Usage:       subnet.SubnetDetailMetaData().Usage,
			Flags:       subnet.SubnetDetailMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Subnet detail", func() {
		Context("Subnet detail without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Subnet detail with wrong subnet id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Subnet ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Subnet detail with correct subnet id but server fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetSubnetReturns(datatypes.Network_Subnet{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get subnet: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Subnet detail with correct subnet id", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetSubnetReturns(datatypes.Network_Subnet{
					Id:                sl.Int(1234),
					NetworkIdentifier: sl.String("9.9.9.9"),
					Cidr:              sl.Int(10),
					SubnetType:        sl.String("PRIMARY"),
					NetworkVlan: &datatypes.Network_Vlan{
						NetworkSpace: sl.String("PUBLIC"),
					},
					Gateway:          sl.String("9.9.9.1"),
					BroadcastAddress: sl.String("9.9.9.0"),
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
					VirtualGuests: []datatypes.Virtual_Guest{
						datatypes.Virtual_Guest{
							Hostname:                sl.String("vs1"),
							Domain:                  sl.String("wilma.com"),
							PrimaryIpAddress:        sl.String("9.9.9.2"),
							PrimaryBackendIpAddress: sl.String("1.2.1.1"),
						},
					},
					Hardware: []datatypes.Hardware{
						datatypes.Hardware{
							Hostname:                sl.String("hw1"),
							Domain:                  sl.String("wilma.org"),
							PrimaryIpAddress:        sl.String("9.9.9.3"),
							PrimaryBackendIpAddress: sl.String("1.2.1.2"),
						},
					},
				}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.9/10"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"PRIMARY"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"PUBLIC"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.1"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.0"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"dal10"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"vs1"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"wilma.com"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.2"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1.2.1.1"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"hw1"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"wilma.org"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.3"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1.2.1.2"}))
			})

			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--no-vs")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.9/10"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"PRIMARY"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"PUBLIC"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.1"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.0"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"dal10"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"vs1"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"wilma.com"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"9.9.9.2"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"1.2.1.1"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"hw1"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"wilma.org"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.3"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1.2.1.2"}))
			})

			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--no-hardware")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.9/10"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"PRIMARY"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"PUBLIC"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.1"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.0"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"dal10"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"vs1"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"wilma.com"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.2"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1.2.1.1"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"hw1"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"wilma.org"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"9.9.9.3"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"1.2.1.2"}))
			})
		})
	})
})

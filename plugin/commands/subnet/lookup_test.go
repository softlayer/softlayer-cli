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

var _ = Describe("Subnet lookup", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *subnet.LookupCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = subnet.NewLookupCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        subnet.SubnetLookupMetaData().Name,
			Description: subnet.SubnetLookupMetaData().Description,
			Usage:       subnet.SubnetLookupMetaData().Usage,
			Flags:       subnet.SubnetLookupMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Subnet lookup", func() {
		Context("Subnet lookup without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Subnet detail with correct ipaddress but server fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.IPLookupReturns(datatypes.Network_Subnet_IpAddress{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1.2.3.4")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to lookup IP address: 1.2.3.4.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Subnet detail with correct ipaddress but not found", func() {
			BeforeEach(func() {
				fakeNetworkManager.IPLookupReturns(datatypes.Network_Subnet_IpAddress{}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1.2.3.4")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"IP address 1.2.3.4 is not found."}))
			})
		})

		Context("Subnet detail with correct ipaddress", func() {
			BeforeEach(func() {
				fakeNetworkManager.IPLookupReturns(datatypes.Network_Subnet_IpAddress{
					Id:        sl.Int(1234),
					IpAddress: sl.String("9.9.9.8"),
					Subnet: &datatypes.Network_Subnet{
						Id:                sl.Int(4567),
						NetworkIdentifier: sl.String("9.9.9.9"),
						Cidr:              sl.Int(10),
						SubnetType:        sl.String("PRIMARY"),
						NetworkVlan: &datatypes.Network_Vlan{
							NetworkSpace: sl.String("PUBLIC"),
						},
						Netmask: sl.String("9.9.9.0"),
						Gateway: sl.String("9.9.9.1"),
						Datacenter: &datatypes.Location_Datacenter{
							Location: datatypes.Location{
								Name: sl.String("dal10"),
							},
						},
					},
					VirtualGuest: &datatypes.Virtual_Guest{
						Id:                       sl.Int(765432),
						Hostname:                 sl.String("vs1"),
						Domain:                   sl.String("wilma.com"),
						PrimaryIpAddress:         sl.String("9.9.9.2"),
						FullyQualifiedDomainName: sl.String("vs1.wilma.com"),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1.2.3.4")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.8"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"4567"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.9/10"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.0"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.1"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"PRIMARY"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"765432"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"vs1.wilma.com"}))
			})
		})
	})
})

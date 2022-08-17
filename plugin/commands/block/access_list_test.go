package block_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

var _ = Describe("Access List Tests", func() {
    var (
        fakeUI              *terminal.FakeUI
        cliCommand          *block.AccessListCommand
        fakeSession         *session.Session
        slCommand           *metadata.SoftlayerCommand
        FakeStorageManager *testhelpers.FakeStorageManager
    )
    BeforeEach(func() {
        fakeUI = terminal.NewFakeUI()
        fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		FakeStorageManager = new(testhelpers.FakeStorageManager)
        slCommand  = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
        cliCommand = block.NewAccessListCommand(slCommand)
        cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
        cliCommand.StorageManager = FakeStorageManager
    })

	Describe("Access List", func() {
		Context("Access list without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("Access Authorize with wrong volume id", func() {
			It("error resolving volume ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
			})
		})

		Context("Access Authorize with correct volume id", func() {
			BeforeEach(func() {
				FakeStorageManager.GetVolumeAccessListReturns(datatypes.Network_Storage{
					AllowedVirtualGuests: []datatypes.Virtual_Guest{
						datatypes.Virtual_Guest{
							Id:                      sl.Int(12345678),
							Hostname:                sl.String("wilma.org"),
							PrimaryBackendIpAddress: sl.String("1.2.3.4"),
							AllowedHost: &datatypes.Network_Storage_Allowed_Host{
								Name: sl.String("vs-abc"),
								Credential: &datatypes.Network_Storage_Credential{
									Username: sl.String("vs-bcd"),
									Password: sl.String("xxxxxxxx"),
								},
							},
						},
					},
					AllowedHardware: []datatypes.Hardware{
						datatypes.Hardware{
							Id:                      sl.Int(87654321),
							Hostname:                sl.String("wilma.com"),
							PrimaryBackendIpAddress: sl.String("4.3.2.1"),
							AllowedHost: &datatypes.Network_Storage_Allowed_Host{
								Name: sl.String("hw-abc"),
								Credential: &datatypes.Network_Storage_Credential{
									Username: sl.String("hw-bcd"),
									Password: sl.String("yyyyyyyy"),
								},
							},
						},
					},
					AllowedSubnets: []datatypes.Network_Subnet{
						datatypes.Network_Subnet{
							Id:                sl.Int(12387654),
							NetworkIdentifier: sl.String("9.9.9.9"),
							EndPointIpAddress: &datatypes.Network_Subnet_IpAddress{
								IpAddress: sl.String("9.9.9.9"),
							},
							AllowedHost: &datatypes.Network_Storage_Allowed_Host{
								Name: sl.String("sn-abc"),
								Credential: &datatypes.Network_Storage_Credential{
									Username: sl.String("sn-bcd"),
									Password: sl.String("zzzzzzzz"),
								},
							},
						},
					},
					AllowedIpAddresses: []datatypes.Network_Subnet_IpAddress{
						datatypes.Network_Subnet_IpAddress{
							Id:        sl.Int(87612345),
							IpAddress: sl.String("8.8.8.8"),
							AllowedHost: &datatypes.Network_Storage_Allowed_Host{
								Name: sl.String("ip-abc"),
								Credential: &datatypes.Network_Storage_Credential{
									Username: sl.String("ip-bcd"),
									Password: sl.String("vvvvvvvv"),
								},
							},
						},
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"12345678"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"wilma.org"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1.2.3.4"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"vs-abc"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"vs-bcd"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"xxxxxxxx"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"87654321"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"wilma.org"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"4.3.2.1"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"hw-abc"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"hw-bcd"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"yyyyyyyy"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"12387654"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.9"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"sn-abc"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"sn-bcd"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"zzzzzzzz"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"87612345"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"8.8.8.8"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"ip-abc"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"ip-bcd"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"vvvvvvvv"}))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--column", "id")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"12345678"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"wilma.org"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"1.2.3.4"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"vs-abc"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"vs-bcd"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"xxxxxxxx"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"87654321"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"wilma.org"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"4.3.2.1"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"hw-abc"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"hw-bcd"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"yyyyyyyy"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"12387654"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"9.9.9.9"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"sn-abc"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"sn-bcd"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"zzzzzzzz"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"87612345"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"8.8.8.8"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"ip-abc"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"ip-bcd"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"vvvvvvvv"}))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--column", "id")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"12345678"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"wilma.org"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"1.2.3.4"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"vs-abc"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"vs-bcd"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"xxxxxxxx"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"87654321"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"wilma.org"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"4.3.2.1"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"hw-abc"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"hw-bcd"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"yyyyyyyy"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"12387654"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"9.9.9.9"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"sn-abc"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"sn-bcd"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"zzzzzzzz"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"87612345"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"8.8.8.8"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"ip-abc"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"ip-bcd"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"vvvvvvvv"}))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--column", "id", "--column", "username")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"12345678"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"wilma.org"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"1.2.3.4"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"vs-abc"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"vs-bcd"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"xxxxxxxx"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"87654321"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"wilma.org"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"4.3.2.1"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"hw-abc"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"hw-bcd"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"yyyyyyyy"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"12387654"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"9.9.9.9"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"sn-abc"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"sn-bcd"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"zzzzzzzz"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"87612345"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"8.8.8.8"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"ip-abc"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"ip-bcd"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"vvvvvvvv"}))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--column", "id", "--column", "username")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"12345678"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"wilma.org"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"1.2.3.4"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"vs-abc"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"vs-bcd"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"xxxxxxxx"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"87654321"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"wilma.org"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"4.3.2.1"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"hw-abc"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"hw-bcd"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"yyyyyyyy"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"12387654"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"9.9.9.9"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"sn-abc"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"sn-bcd"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"zzzzzzzz"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"87612345"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"8.8.8.8"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"ip-abc"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"ip-bcd"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"vvvvvvvv"}))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(MatchRegexp(`.*12345678.*\n.*12387654.*\n.*87612345.*\n.*87654321.*`))
			})
		})

		Context("Access Authorize with correct volume id but server fails", func() {
			BeforeEach(func() {
				FakeStorageManager.GetVolumeAccessListReturns(datatypes.Network_Storage{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"OK"}))
				Expect(strings.Contains(err.Error(), "Failed to get access list for volume 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
	})
})

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
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Access Authorize", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *block.AccessAuthorizeCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = block.NewAccessAuthorizeCommand(fakeUI, FakeStorageManager, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        block.BlockAccessAuthorizeMetaData().Name,
			Description: block.BlockAccessAuthorizeMetaData().Description,
			Usage:       block.BlockAccessAuthorizeMetaData().Usage,
			Flags:       block.BlockAccessAuthorizeMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Access Authorize", func() {
		Context("Access Authorize without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)

				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Access Authorize with wrong volume id", func() {
			It("error resolving volume ID", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
			})
		})

		Context("Access Authorize with correct volume id and virtual server id", func() {
			BeforeEach(func() {
				FakeStorageManager.AuthorizeHostToVolumeReturns([]datatypes.Network_Storage_Allowed_Host{}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--virtual-id", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The virtual server 5678 was authorized to access 1234."}))
			})
		})

		Context("Access Authorize with correct volume id and hardware server id", func() {
			BeforeEach(func() {
				FakeStorageManager.AuthorizeHostToVolumeReturns([]datatypes.Network_Storage_Allowed_Host{}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--hardware-id", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The hardware server 5678 was authorized to access 1234."}))
			})
		})

		Context("Access Authorize with correct volume id and ip address id", func() {
			BeforeEach(func() {
				FakeStorageManager.AuthorizeHostToVolumeReturns([]datatypes.Network_Storage_Allowed_Host{}, nil)
			})
			It("Success with multipl IP Ids", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--ip-address-id", "5678", "--ip-address-id", "9999")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The IP address 5678 was authorized to access 1234."}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The IP address 9999 was authorized to access 1234."}))
				volId, _, _, ipArg, _ := FakeStorageManager.AuthorizeHostToVolumeArgsForCall(0)
				Expect(ipArg).To(Equal([]int{5678, 9999}))
				Expect(volId).To(Equal(1234))
			})
			It("Success with single IP Ids", func() {
				// Testing this because when splitting out sl into its own module, intSlices seem to be duplicating first value
				err := testhelpers.RunCommand(cliCommand, "1234", "--ip-address-id", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The IP address 5678 was authorized to access 1234."}))
				volId, _, _, ipArg, _ := FakeStorageManager.AuthorizeHostToVolumeArgsForCall(0)
				Expect(ipArg).To(Equal([]int{5678}))
				Expect(volId).To(Equal(1234))
			})
		})

		Context("Access Authorize with correct volume id and ip address", func() {
			BeforeEach(func() {
				FakeStorageManager.AuthorizeHostToVolumeReturns([]datatypes.Network_Storage_Allowed_Host{}, nil)
				fakeNetworkManager.IPLookupReturns(datatypes.Network_Subnet_IpAddress{Id: sl.Int(5678)}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--ip-address", "1.2.3.4")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The IP address 5678 was authorized to access 1234."}))
			})
		})

		Context("Access Authorize with correct volume id and wrong ip address", func() {
			BeforeEach(func() {
				FakeStorageManager.AuthorizeHostToVolumeReturns([]datatypes.Network_Storage_Allowed_Host{}, nil)
				fakeNetworkManager.IPLookupReturns(datatypes.Network_Subnet_IpAddress{}, errors.New("Not Found"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--ip-address", "1.2.3.4")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"OK"}))
				Expect(strings.Contains(err.Error(), "IP address 1.2.3.4 is not found on your account.Please confirm IP and try again.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Not Found")).To(BeTrue())
			})
		})

		Context("Access Authorize with correct volume id but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.AuthorizeHostToVolumeReturns([]datatypes.Network_Storage_Allowed_Host{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--virtual-id", "5678")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"OK"}))
				Expect(strings.Contains(err.Error(), "Failed to authorize host to volume")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
	})
})

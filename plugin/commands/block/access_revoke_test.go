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
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Access Revoke", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *block.AccessRevokeCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = block.NewAccessRevokeCommand(fakeUI, FakeStorageManager, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        metadata.BlockAccessRevokeMetaData().Name,
			Description: metadata.BlockAccessRevokeMetaData().Description,
			Usage:       metadata.BlockAccessRevokeMetaData().Usage,
			Flags:       metadata.BlockAccessRevokeMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Access Revoke", func() {
		Context("Access revoke without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Access revoke with wrong volume id", func() {
			It("error resolving volume ID", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Volume ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Access revoke with correct volume id and virtual server id", func() {
			BeforeEach(func() {
				FakeStorageManager.DeauthorizeHostToVolumeReturns([]datatypes.Network_Storage_Allowed_Host{}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--virtual-id", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Access to 1234 was revoked for virtual server 5678"}))
			})
		})

		Context("Access revoke with correct volume id and hardware server id", func() {
			BeforeEach(func() {
				FakeStorageManager.DeauthorizeHostToVolumeReturns([]datatypes.Network_Storage_Allowed_Host{}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--hardware-id", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Access to 1234 was revoked for hardware server 5678."}))
			})
		})

		Context("Access revoke with correct volume id and ip address id", func() {
			BeforeEach(func() {
				FakeStorageManager.DeauthorizeHostToVolumeReturns([]datatypes.Network_Storage_Allowed_Host{}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--ip-address-id", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Access to 1234 was revoked for IP address 5678."}))
			})
		})

		Context("Access revoke with correct volume id and ip address", func() {
			BeforeEach(func() {
				FakeStorageManager.DeauthorizeHostToVolumeReturns([]datatypes.Network_Storage_Allowed_Host{}, nil)
				fakeNetworkManager.IPLookupReturns(datatypes.Network_Subnet_IpAddress{Id: sl.Int(5678)}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--ip-address", "1.2.3.4")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Access to 1234 was revoked for IP address 5678."}))
			})
		})

		Context("Access revoke with correct volume id and wrong ip address", func() {
			BeforeEach(func() {
				FakeStorageManager.DeauthorizeHostToVolumeReturns([]datatypes.Network_Storage_Allowed_Host{}, nil)
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
				FakeStorageManager.DeauthorizeHostToVolumeReturns([]datatypes.Network_Storage_Allowed_Host{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--virtual-id", "5678")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"OK"}))
				Expect(strings.Contains(err.Error(), "Failed to revoke access to volume 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
	})
})
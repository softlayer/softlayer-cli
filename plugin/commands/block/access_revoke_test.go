package block_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Access Revoke", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *block.AccessRevokeCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
		FakeStorageManager *testhelpers.FakeStorageManager
		fakeNetworkManager *testhelpers.FakeNetworkManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewAccessRevokeCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
		cliCommand.NetworkManager = fakeNetworkManager
		FakeStorageManager.GetVolumeIdReturns(1234, nil)
	})

	Describe("Access Revoke", func() {
		Context("Syntax Errors", func() {
			It("Require One Argument", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})

		Context("Successful Revokations", func() {
			BeforeEach(func() {
				FakeStorageManager.DeauthorizeHostToVolumeReturns([]datatypes.Network_Storage_Allowed_Host{}, nil)
			})
			It("Virtual Server", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--virtual-id", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Access to 1234 was revoked for Virtual Server 5678"))
			})
			It("Hardware Server", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--hardware-id", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Access to 1234 was revoked for Hardware Server 5678."))
			})
			It("Single IP Address ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--ip-address-id", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Access to 1234 was revoked for IP Address 5678."))
			})
			It("Single IP address", func() {
				fakeNetworkManager.IPLookupReturns(datatypes.Network_Subnet_IpAddress{Id: sl.Int(5678)}, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--ip-address", "1.2.3.4")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Access to 1234 was revoked for IP Address 5678."))
			})
		})

		Context("Error Handling", func() {
			BeforeEach(func() {
				FakeStorageManager.DeauthorizeHostToVolumeReturns([]datatypes.Network_Storage_Allowed_Host{}, nil)				
			})
			It("IP Not Found", func() {
				fakeNetworkManager.IPLookupReturns(datatypes.Network_Subnet_IpAddress{}, errors.New("Not Found"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--ip-address", "1.2.3.4")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("IP address 1.2.3.4 is not found on your account."))
				Expect(err.Error()).To(ContainSubstring("Not Found"))
			})
			It("API error", func() {
				FakeStorageManager.DeauthorizeHostToVolumeReturns(
					[]datatypes.Network_Storage_Allowed_Host{}, errors.New("Internal Server Error"),
				)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--virtual-id", "5678")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to revoke access to volume 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
			It("Subnet not removed because isci isolation", func() {
				FakeStorageManager.RemoveSubnetsFromAclReturns([]int{}, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--subnet-id", "5678")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to remove subnet id: 5678"))
			})
			It("Subnet not removed because wrong subnet returned", func() {
				FakeStorageManager.RemoveSubnetsFromAclReturns([]int{999}, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--subnet-id", "5678")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to remove subnet id: 5678"))
			})
			It("Subnet not removed because API error", func() {
				FakeStorageManager.RemoveSubnetsFromAclReturns([]int{}, errors.New("API ERROR"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--subnet-id", "5678")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("API ERROR"))
			})
		})
	})
})

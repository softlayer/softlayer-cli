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

var _ = Describe("Access Authorize", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *block.AccessAuthorizeCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
		FakeStorageManager *testhelpers.FakeStorageManager
		fakeNetworkManager *testhelpers.FakeNetworkManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession(nil)
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewAccessAuthorizeCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
		cliCommand.NetworkManager = fakeNetworkManager
		FakeStorageManager.GetVolumeIdReturns(1234, nil)
	})

	Describe("Access Authorize", func() {
		Context("Syntax Errors", func() {
			It("Require One Argument", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
			It("Valid allowed host id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "test", "--subnet-id=1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Allowed Host IDENTIFIER'"))
			})
		})
		Context("Successful Authorizations", func() {
			BeforeEach(func() {
				FakeStorageManager.AuthorizeHostToVolumeReturns([]datatypes.Network_Storage_Allowed_Host{}, nil)
			})
			It("Virtual Server", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--virtual-id", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The Virtual Server 5678 was authorized to access 1234."))
			})
			It("Hardware Server", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--hardware-id", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The Hardware Server 5678 was authorized to access 1234."))
			})
			It("Multiple IP Address IDs", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--ip-address-id", "5678", "--ip-address-id", "9999")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The IP Address 5678 was authorized to access 1234."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("The IP Address 9999 was authorized to access 1234."))
				volId, _, _, ipArg, _ := FakeStorageManager.AuthorizeHostToVolumeArgsForCall(0)
				Expect(ipArg).To(Equal([]int{5678, 9999}))
				Expect(volId).To(Equal(1234))
			})
			It("Single IP Address ID", func() {
				// Testing this because when splitting out sl into its own module, intSlices seem to be duplicating first value
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--ip-address-id", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The IP Address 5678 was authorized to access 1234."))
				volId, _, _, ipArg, _ := FakeStorageManager.AuthorizeHostToVolumeArgsForCall(0)
				Expect(ipArg).To(Equal([]int{5678}))
				Expect(volId).To(Equal(1234))
			})
			It("IP Address", func() {
				fakeNetworkManager.IPLookupReturns(datatypes.Network_Subnet_IpAddress{Id: sl.Int(5678)}, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--ip-address", "1.2.3.4")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The IP Address 5678 was authorized to access 1234."))
			})
			It("Subnet", func() {
				FakeStorageManager.AssignSubnetsToAclReturns([]int{5678}, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--subnet-id", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The Subnet 5678 was authorized to access 1234."))
			})
		})

		Context("Error Handling", func() {
			It("IP Address not found", func() {
				FakeStorageManager.AuthorizeHostToVolumeReturns([]datatypes.Network_Storage_Allowed_Host{}, nil)
				fakeNetworkManager.IPLookupReturns(datatypes.Network_Subnet_IpAddress{}, errors.New("Not Found"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--ip-address", "1.2.3.4")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("OK"))
				Expect(err.Error()).To(ContainSubstring("IP address 1.2.3.4 is not found on your account."))
				Expect(err.Error()).To(ContainSubstring("Not Found"))
			})
			It("Other API Error", func() {
				FakeStorageManager.AuthorizeHostToVolumeReturns(
					[]datatypes.Network_Storage_Allowed_Host{}, errors.New("Internal Server Error"),
				)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--virtual-id", "5678")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("OK"))
				Expect(err.Error()).To(ContainSubstring("Failed to authorize host to volume"))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
			It("Subnet not added because isci isolation", func() {
				FakeStorageManager.AssignSubnetsToAclReturns([]int{}, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--subnet-id", "5678")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Make sure ISCSI Isolation is enabled for this account"))
			})
			It("Subnet not added because wrong subnet returned", func() {
				FakeStorageManager.AssignSubnetsToAclReturns([]int{999}, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--subnet-id", "5678")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Make sure ISCSI Isolation is enabled for this account"))
			})
			It("Subnet not added because API error", func() {
				FakeStorageManager.AssignSubnetsToAclReturns([]int{}, errors.New("API ERROR"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--subnet-id", "5678")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to assign subnet id: 5678 to allowed host id: 1234"))
				Expect(err.Error()).To(ContainSubstring("API ERROR"))
			})
		})
	})
})

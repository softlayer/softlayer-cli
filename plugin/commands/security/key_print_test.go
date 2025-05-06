package security_test

import (
	"errors"
	"os"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/security"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Key print", func() {
	var (
		fakeUI              *terminal.FakeUI
		cliCommand          *security.KeyPrintCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
		fakeSecurityManager *testhelpers.FakeSecurityManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = security.NewKeyPrintCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeSecurityManager = new(testhelpers.FakeSecurityManager)
		cliCommand.SecurityManager = fakeSecurityManager
	})

	Describe("Key print", func() {
		Context("Key print without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("Key print with wrong key ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'SSH Key ID'. It must be a positive integer."))
			})
		})
		Context("Key print but server API call fails", func() {
			BeforeEach(func() {
				fakeSecurityManager.GetSSHKeyReturns(datatypes.Security_Ssh_Key{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get SSH Key 1234"))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})
		Context("Key print", func() {
			BeforeEach(func() {
				fakeSecurityManager.GetSSHKeyReturns(datatypes.Security_Ssh_Key{
					Id:    sl.Int(1234),
					Label: sl.String("label"),
					Notes: sl.String("notes"),
					Key:   sl.String("ssh-rsa djghtbtmfhgentongwfrdnglkhsdye"),
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"label"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"notes"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"ssh-rsa djghtbtmfhgentongwfrdnglkhsdye"}))
			})
			It("return no error", func() {
				if os.Getenv("OS") == "Windows_NT" {
					Skip("Test doesn't work in windows.")
				}
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f", "/tmp/key")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"label"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"notes"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"ssh-rsa djghtbtmfhgentongwfrdnglkhsdye"}))
			})
		})
	})
})

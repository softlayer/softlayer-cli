package security_test

import (
	"errors"
	"strings"

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

var _ = Describe("Key List", func() {
	var (
		fakeUI              *terminal.FakeUI
		cliCommand          *security.KeyListCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
		fakeSecurityManager *testhelpers.FakeSecurityManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = security.NewKeyListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeSecurityManager = new(testhelpers.FakeSecurityManager)
		cliCommand.SecurityManager = fakeSecurityManager
	})

	Describe("Key list", func() {
		Context("Key list but server API call fails", func() {
			BeforeEach(func() {
				fakeSecurityManager.ListSSHKeysReturns([]datatypes.Security_Ssh_Key{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list SSH keys on your account."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Key list with wrong sortby", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --sortby abc is not supported."))
			})
		})

		Context("Key list with different sortby", func() {
			BeforeEach(func() {
				fakeSecurityManager.ListSSHKeysReturns([]datatypes.Security_Ssh_Key{
					datatypes.Security_Ssh_Key{
						Id:          sl.Int(123),
						Label:       sl.String("mon"),
						Fingerprint: sl.String("92:e1:82:20:c4:f4:c3:1c:ca:57:ce:5f:10:5b:93:31"),
						Notes:       sl.String("Docker"),
					},
					datatypes.Security_Ssh_Key{
						Id:          sl.Int(110),
						Label:       sl.String("nom"),
						Fingerprint: sl.String("25:2a:e1:97:57:e4:d2:7c:ed:e0:57:85:eb:e9:c2:a8"),
						Notes:       sl.String("Armer"),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("110"))
				Expect(results[2]).To(ContainSubstring("123"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "label")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("mon"))
				Expect(results[2]).To(ContainSubstring("nom"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "fingerprint")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("25:2a:e1:97:57:e4:d2:7c:ed:e0:57:85:eb:e9:c2:a8"))
				Expect(results[2]).To(ContainSubstring("92:e1:82:20:c4:f4:c3:1c:ca:57:ce:5f:10:5b:93:31"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "note")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("Armer"))
				Expect(results[2]).To(ContainSubstring("Docker"))
			})
		})
	})
})

package security_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/security"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Key List", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeSecurityManager *testhelpers.FakeSecurityManager
		cmd                 *security.KeyListCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSecurityManager = new(testhelpers.FakeSecurityManager)
		cmd = security.NewKeyListCommand(fakeUI, fakeSecurityManager)
		cliCommand = cli.Command{
			Name:        metadata.SecuritySSHKeyListMetaData().Name,
			Description: metadata.SecuritySSHKeyListMetaData().Description,
			Usage:       metadata.SecuritySSHKeyListMetaData().Usage,
			Flags:       metadata.SecuritySSHKeyListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Key list", func() {
		Context("Key list but server API call fails", func() {
			BeforeEach(func() {
				fakeSecurityManager.ListSSHKeysReturns([]datatypes.Security_Ssh_Key{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to list SSH keys on your account.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Key list with wrong sortby", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: --sortby abc is not supported.")).To(BeTrue())
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
				err := testhelpers.RunCommand(cliCommand, "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "110")).To(BeTrue())
				Expect(strings.Contains(results[2], "123")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "label")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "mon")).To(BeTrue())
				Expect(strings.Contains(results[2], "nom")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "fingerprint")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "25:2a:e1:97:57:e4:d2:7c:ed:e0:57:85:eb:e9:c2:a8")).To(BeTrue())
				Expect(strings.Contains(results[2], "92:e1:82:20:c4:f4:c3:1c:ca:57:ce:5f:10:5b:93:31")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "note")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "Armer")).To(BeTrue())
				Expect(strings.Contains(results[2], "Docker")).To(BeTrue())
			})
		})
	})
})

package security_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"reflect"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/security"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Security Suite")
}

// These are all the commands in security.go
var availableCommands = []string{
	"ssl-add",
	"ssl-edit",
	"ssl-list",
	"ssl-remove",
	"security-sshkey-list",
	"sshkey-add",
	"sshkey-print",
	"security-cert-download",
	"security-cert-remove",
	"ssl-download",
	"security-sshkey-add",
	"security-sshkey-print",
	"sshkey-edit",
	"sshkey-remove",
	"sshkey-list",
	"security-cert-add",
	"security-cert-edit",
	"security-cert-list",
	"security-sshkey-edit",
	"security-sshkey-remove",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test security.GetCommandActionBindings()", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	commands := security.GetCommandActionBindings(fakeUI, fakeSession)
	Context("Test Actions", func() {
		for _, cmdName := range availableCommands {
			//necessary to ensure the correct value is passed to the closure
			cmdName := cmdName
			It("ibmcloud sl "+cmdName, func() {
				command, exists := commands[cmdName]
				Expect(exists).To(BeTrue(), cmdName+" not found")
				// Checks to make sure we actually have a function here.
				// Test the actual function works in the specific commands test file.
				Expect(reflect.ValueOf(command).Kind().String()).To(Equal("func"))
				context := testhelpers.GetCliContext(cmdName)
				err := command(context)
				// some commands work without arguments
				if err == nil {
					Expect(err).NotTo(HaveOccurred())
				} else {
					Expect(err).To(HaveOccurred())
				}
			})
		}
	})
	Context("New commands testable", func() {
		for cmdName, _ := range commands {
			//necessary to ensure the correct value is passed to the closure
			cmdName := cmdName
			It("availableCommands["+cmdName+"]", func() {
				found := false
				for _, value := range availableCommands {
					if value == cmdName {
						found = true
						break
					}
				}
				Expect(found).To(BeTrue(), cmdName+" needs to be added to availableCommands[] in securty_test.go")
			})
		}
	})
})

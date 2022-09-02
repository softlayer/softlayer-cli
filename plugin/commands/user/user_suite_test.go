package user_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestUser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the SetupCobraCommands
var _ = Describe("Test user commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)

	Context("New commands testable", func() {
		autoscaleCommands := user.SetupCobraCommands(slMeta)
		Expect(autoscaleCommands.Name()).To(Equal("user"))
	})

	Context("User Namespace", func() {
		It("User Name Space", func() {
			Expect(user.UserNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(user.UserNamespace().Name).To(ContainSubstring("user"))
			Expect(user.UserNamespace().Description).To(ContainSubstring("Classic infrastructure Manage Users"))
		})
	})
})

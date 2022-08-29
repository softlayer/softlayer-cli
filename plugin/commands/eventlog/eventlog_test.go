package eventlog_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/eventlog"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Event Log Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the SetupCobraCommands
var _ = Describe("Test eventlog commands", func() {

	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)

	Context("New commands testable", func() {
		eventlogCommands := eventlog.SetupCobraCommands(slMeta)
		Expect(eventlogCommands.Name()).To(Equal("event-log"))
	})

	Context("Event Log Namespace", func() {
		It("Event Log Name Space", func() {
			Expect(eventlog.EventLogNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(eventlog.EventLogNamespace().Name).To(ContainSubstring("event-log"))
			Expect(eventlog.EventLogNamespace().Description).To(ContainSubstring("Classic infrastructure Event Log Group"))
		})
	})
})

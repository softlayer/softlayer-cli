package meta_test

import (
	"fmt"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/meta"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Metadata Suite")
}

var _ = Describe("Metadata list Metadata", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *meta.MetaCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
		// fakeManager   *testhelpers.FakeMetadataManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		// fakeManager = new(testhelpers.FakeMetadataManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = meta.NewMetaCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		// cliCommand.Manager = fakeManager
	})

	Describe("Metadata command", func() {
		Context("Metadata options, Invalid Usage", func() {
			It("Set command without option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
			It("Set unavailable option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid argument \"abc\""))
			})
		})

		Context("Metadata options, correct use", func() {
			It("return table with all datas from network", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "network")
				Expect(err).NotTo(HaveOccurred())
				fmt.Println("responseee:", fakeUI.Outputs())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name            Value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Mac addresses   \"00:a1:b2:c3:d4:e5\""))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Router          fcr02.dal06"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Vlans           12345"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Vlan ids        1234567"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Mac addresses   \"11:a1:b2:c3:d4:e5\""))
			})
			It("return backend id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "backend_ip")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("10.222.111.199"))
			})
			It("return backend ip", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "backend_ip")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("10.222.111.199"))
			})
			It("return backend mac", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "backend_mac")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"11:a1:b2:c3:d4:e5"`))
			})
			It("return name datacenter", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "datacenter")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal06"))
			})
			It("return id datacenter", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "datacenter_id")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("7654321"))
			})
			It("return fully qualified domain name", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "fqdn")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("hostname.test.com"))
			})
			It("return frontend mac", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "frontend_mac")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"00:a1:b2:c3:d4:e5"`))
			})
			It("return id machine", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "id")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("987654321"))
			})
			It("return ip machine", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "ip")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("169.11.22.199"))
			})
			It("return provision state", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "provision_state")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("COMPLETE"))
			})
			It("return tags", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "tags")
				Expect(err).NotTo(HaveOccurred())
				fmt.Println("tagss:", fakeUI.Outputs())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"testTags"`))
			})
			It("return user data", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "user_data")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("userData"))
			})
		})
	})
})

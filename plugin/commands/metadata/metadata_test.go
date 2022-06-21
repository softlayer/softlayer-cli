package metadata_test

import (
	"reflect"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Metadata Suite")
}

var availableCommands = []string{
	"sl-metadata",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test metadata.GetCommandActionBindings()", func() {
	var (
		context plugin.PluginContext
	)
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	context = plugin.InitPluginContext("softlayer")
	commands := metadata.GetCommandActionBindings(context, fakeUI, fakeSession)

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
				Expect(found).To(BeTrue(), cmdName+" needs to be added to availableCommands[] in metadata.go")
			})
		}
	})
})

var _ = Describe("Metadata list Metadata", func() {
	var (
		fakeUI              *terminal.FakeUI
		cmd                 *metadata.MetadataCommand
		cliCommand          cli.Command
		fakeSession         *session.Session
		fakeMetadataManager managers.MetadataManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeMetadataManager = managers.NewMetadataManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		cmd = metadata.NewMetadataCommand(fakeUI, fakeMetadataManager)
		cliCommand = cli.Command{
			Name:        metadata.MetadataMetadata().Name,
			Description: metadata.MetadataMetadata().Description,
			Usage:       metadata.MetadataMetadata().Usage,
			Flags:       metadata.MetadataMetadata().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Metadata command", func() {
		Context("Metadata options, Invalid Usage", func() {
			It("Set command without option", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
			It("Set unavailable option", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This option is not available."))
			})
		})

		Context("Metadata options, correct use", func() {
			It("return table with all datas from network", func() {
				err := testhelpers.RunCommand(cliCommand, "network")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name            Value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Mac addresses   00:a1:b2:c3:d4:e5"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Router          fcr02.dal06"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Vlans           12345"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Vlan ids        1234567"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Mac addresses   11:a1:b2:c3:d4:e5"))
			})
			It("return backend id", func() {
				err := testhelpers.RunCommand(cliCommand, "backend_ip")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("10.222.111.199"))
			})
			It("return backend ip", func() {
				err := testhelpers.RunCommand(cliCommand, "backend_ip")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("10.222.111.199"))
			})
			It("return backend mac", func() {
				err := testhelpers.RunCommand(cliCommand, "backend_mac")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("11:a1:b2:c3:d4:e5"))
			})
			It("return name datacenter", func() {
				err := testhelpers.RunCommand(cliCommand, "datacenter")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal06"))
			})
			It("return id datacenter", func() {
				err := testhelpers.RunCommand(cliCommand, "datacenter_id")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("7654321"))
			})
			It("return fully qualified domain name", func() {
				err := testhelpers.RunCommand(cliCommand, "fqdn")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("hostname.test.com"))
			})
			It("return frontend mac", func() {
				err := testhelpers.RunCommand(cliCommand, "frontend_mac")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("00:a1:b2:c3:d4:e5"))
			})
			It("return id machine", func() {
				err := testhelpers.RunCommand(cliCommand, "id")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("987654321"))
			})
			It("return ip machine", func() {
				err := testhelpers.RunCommand(cliCommand, "ip")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("169.11.22.199"))
			})
			It("return provision state", func() {
				err := testhelpers.RunCommand(cliCommand, "provision_state")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("COMPLETE"))
			})
			It("return tags", func() {
				err := testhelpers.RunCommand(cliCommand, "tags")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("testTags"))
			})
			It("return user data", func() {
				err := testhelpers.RunCommand(cliCommand, "user_data")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("userData"))
			})
		})
	})
})

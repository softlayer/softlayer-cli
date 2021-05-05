package virtual_test

import (
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS edit", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.EditCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewEditCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        metadata.VSEditMataData().Name,
			Description: metadata.VSEditMataData().Description,
			Usage:       metadata.VSEditMataData().Usage,
			Flags:       metadata.VSEditMataData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS edit", func() {
		Context("VS edit without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("VS edit with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Virtual server ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("VS edit with both -u and -f", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-u", "mydata", "-F", "/tmp/datafile")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: '[-u|--userdata]', '[-F|--userfile]' are exclusive.")).To(BeTrue())
			})
		})
		Context("VS edit with wrong public speed", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--public-speed", "9")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: Public network interface speed must be in: 0, 10, 100, 1000, 10000 (Mbps).")).To(BeTrue())
			})
		})
		Context("VS edit with wrong private speed", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--private-speed", "9")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: Private network interface speed must be in: 0, 10, 100, 1000, 10000 (Mbps).")).To(BeTrue())
			})
		})

		Context("VS edit with hostname fails", func() {
			BeforeEach(func() {
				fakeVSManager.EditInstanceReturns([]bool{false}, []string{"Failed to update the hostname/domain/note of virtual server instance: 1234."})
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-H", "vs-abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to update the hostname/domain/note of virtual server instance: 1234.")).To(BeTrue())
			})
		})
		Context("VS edit with hostname succeed", func() {
			BeforeEach(func() {
				fakeVSManager.EditInstanceReturns([]bool{true}, []string{"The hostname of virtual server instance: 1234 was updated."})
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-H", "vs-abc")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The hostname of virtual server instance: 1234 was updated."}))
			})
		})

		Context("VS edit with domain fails", func() {
			BeforeEach(func() {
				fakeVSManager.EditInstanceReturns([]bool{false}, []string{"Failed to update the hostname/domain/note of virtual server instance: 1234."})
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-D", "wilma.com")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to update the hostname/domain/note of virtual server instance: 1234.")).To(BeTrue())
			})
		})
		Context("VS edit with domain succeed", func() {
			BeforeEach(func() {
				fakeVSManager.EditInstanceReturns([]bool{true}, []string{"The domain of virtual server instance: 1234 was updated."})
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-D", "wilma.com")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The domain of virtual server instance: 1234 was updated."}))
			})
		})
		Context("VS edit with userdata fails", func() {
			BeforeEach(func() {
				fakeVSManager.EditInstanceReturns([]bool{false}, []string{"Failed to update the user data of virtual server instance: 1234."})
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-u", "mydata")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to update the user data of virtual server instance: 1234.")).To(BeTrue())
			})
		})
		Context("VS edit with user data succeed", func() {
			BeforeEach(func() {
				fakeVSManager.EditInstanceReturns([]bool{true}, []string{"The user data of virtual server instance: 1234 was updated."})
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-u", "mydata")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The user data of virtual server instance: 1234 was updated."}))
			})
		})
		Context("VS edit with tags fails", func() {
			BeforeEach(func() {
				fakeVSManager.EditInstanceReturns([]bool{false}, []string{"Failed to update the tags of virtual server instance: 1234."})
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--tag", "mytags")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to update the tags of virtual server instance: 1234.")).To(BeTrue())
			})
		})
		Context("VS edit with tag succeed", func() {
			BeforeEach(func() {
				fakeVSManager.EditInstanceReturns([]bool{true}, []string{"The tags of virtual server instance: 1234 was updated."})
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--tag", "mytags")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The tags of virtual server instance: 1234 was updated."}))
			})
		})
		Context("VS edit with public-speed fails", func() {
			BeforeEach(func() {
				fakeVSManager.EditInstanceReturns([]bool{false}, []string{"Failed to update the public network speed of virtual server instance: 1234."})
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--public-speed", "1000")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to update the public network speed of virtual server instance: 1234.")).To(BeTrue())
			})
		})
		Context("VS edit with public-speed succeed", func() {
			BeforeEach(func() {
				fakeVSManager.EditInstanceReturns([]bool{true}, []string{"The public network speed of virtual server instance: 1234 was updated."})
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--public-speed", "1000")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The public network speed of virtual server instance: 1234 was updated."}))
			})
		})
		Context("VS edit with private-speed fails", func() {
			BeforeEach(func() {
				fakeVSManager.EditInstanceReturns([]bool{false}, []string{"Failed to update the private network speed of virtual server instance: 1234."})
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--private-speed", "1000")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to update the private network speed of virtual server instance: 1234.")).To(BeTrue())
			})
		})
		Context("VS edit with private-speed succeed", func() {
			BeforeEach(func() {
				fakeVSManager.EditInstanceReturns([]bool{true}, []string{"The private network speed of virtual server instance: 1234 was updated."})
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--private-speed", "1000")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The private network speed of virtual server instance: 1234 was updated."}))
			})
		})
	})
})

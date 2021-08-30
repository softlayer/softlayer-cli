package virtual_test

import (
	"errors"
	"fmt"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS reload", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.ReloadCommand
		cliCommand    cli.Command
		context       plugin.PluginContext
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		context = plugin.InitPluginContext("softlayer")
		cmd = virtual.NewReloadCommand(fakeUI, fakeVSManager, context)
		cliCommand = cli.Command{
			Name:        metadata.VSReloadMetaData().Name,
			Description: metadata.VSReloadMetaData().Description,
			Usage:       metadata.VSReloadMetaData().Usage,
			Flags:       metadata.VSReloadMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS reload", func() {
		Context("VS reload without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("VS reload with wrong vs ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Virtual server ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("VS reload with correct vs ID but not continue", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This will reload operating system of virtual server instance: 1234 and cannot be undone. Continue?"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted."}))
			})
		})

		Context("VS reload with correct vs ID but server fails", func() {
			BeforeEach(func() {
				fakeVSManager.ReloadInstanceReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to reload virtual server instance: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("VS reload with correct vs ID ", func() {
			BeforeEach(func() {
				fakeVSManager.ReloadInstanceReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{fmt.Sprintf("System reloading for virtual server instance: 1234 is in progress. Run '%s sl vs ready 1234' to check whether it is ready later on.", cmd.Context.CLIName())}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f", "-i", "http://abc/script.sh")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{fmt.Sprintf("System reloading for virtual server instance: 1234 is in progress. Run '%s sl vs ready 1234' to check whether it is ready later on.", cmd.Context.CLIName())}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f", "--image", "456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{fmt.Sprintf("System reloading for virtual server instance: 1234 is in progress. Run '%s sl vs ready 1234' to check whether it is ready later on.", cmd.Context.CLIName())}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f", "-k", "456", "-k", "678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{fmt.Sprintf("System reloading for virtual server instance: 1234 is in progress. Run '%s sl vs ready 1234' to check whether it is ready later on.", cmd.Context.CLIName())}))
			})
		})
	})
})

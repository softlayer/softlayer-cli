package virtual_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS resume", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.ResumeCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewResumeCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        virtual.VSResumeMetaData().Name,
			Description: virtual.VSResumeMetaData().Description,
			Usage:       virtual.VSResumeMetaData().Usage,
			Flags:       virtual.VSResumeMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS resume", func() {
		Context("VS resume without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("VS resume with wrong vs ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Virtual server ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("VS resume with correct vs ID but not continue", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This will resume virtual server instance: 1234. Continue?"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted."}))
			})
		})

		Context("VS resume with correct vs ID but server fails", func() {
			BeforeEach(func() {
				fakeVSManager.ResumeInstanceReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to resume virtual server instance: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("VS resume with correct vs ID ", func() {
			BeforeEach(func() {
				fakeVSManager.ResumeInstanceReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Virtual server instance: 1234 was resumed."}))
			})
		})
	})
})

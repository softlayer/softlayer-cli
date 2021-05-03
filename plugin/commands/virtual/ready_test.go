package virtual_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS ready", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.ReadyCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewReadyCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        metadata.VSReadyMataData().Name,
			Description: metadata.VSReadyMataData().Description,
			Usage:       metadata.VSReadyMataData().Usage,
			Flags:       metadata.VSReadyMataData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS ready", func() {
		Context("VS ready without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("VS ready with wrong vs ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Virtual server ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("VS ready with correct vs ID but server fails", func() {
			BeforeEach(func() {
				fakeVSManager.InstanceIsReadyReturns(false, "", errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to check virtual server instance 1234 is ready.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("VS ready with correct vs ID ", func() {
			BeforeEach(func() {
				fakeVSManager.InstanceIsReadyReturns(true, "", nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Virtual server instance: 1234 is ready."}))
			})
		})

		Context("VS ready with correct vs ID ", func() {
			BeforeEach(func() {
				fakeVSManager.InstanceIsReadyReturns(false, "Virtual guest instance 1234 is paused.", nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Not ready: Virtual guest instance 1234 is paused."}))
			})
		})
	})
})

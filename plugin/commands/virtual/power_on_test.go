package virtual_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS poweron", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.PowerOnCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewPowerOnCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        virtual.VSPowerOnMetaData().Name,
			Description: virtual.VSPowerOnMetaData().Description,
			Usage:       virtual.VSPowerOnMetaData().Usage,
			Flags:       virtual.VSPowerOnMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS poweron", func() {
		Context("VS poweron without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("VS poweron with wrong vs ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Virtual server ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("VS poweron with correct vs ID but not continue", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This will power on virtual server instance: 1234. Continue?"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted."}))
			})
		})

		Context("VS poweron with correct vs ID but server fails", func() {
			It("return error", func() {
				fakeVSManager.PowerOnInstanceReturns(errors.New("Internal Server Error"))
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to power on virtual server instance: 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
			It("return error2", func() {
				fakeVSManager.PowerOnInstanceReturns(errors.New("{\"error\":\"Internal Error\",\"code\":\"SoftLayer_Exception_Public\"}"))
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to power on virtual server instance: 1234."))
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_Public"))
			})
			It("return error3", func() {
				returnErr := sl.Error{
					StatusCode: 200,
					Exception:  "SoftLayer_Exception_Public",
					Message:    "Internal Error",
				}
				fakeVSManager.PowerOnInstanceReturns(returnErr)
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to power on virtual server instance: 1234."))
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_Public"))
			})
		})

		Context("VS poweron with correct vs ID ", func() {
			BeforeEach(func() {
				fakeVSManager.PowerOnInstanceReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Virtual server instance: 1234 was power on."}))
			})
		})
	})
})

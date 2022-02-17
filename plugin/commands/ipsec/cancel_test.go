package ipsec_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ipsec"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("IPSec cancel", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeIPSecManager *testhelpers.FakeIPSECManager
		cmd              *ipsec.CancelCommand
		cliCommand       cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeIPSecManager = new(testhelpers.FakeIPSECManager)
		cmd = ipsec.NewCancelCommand(fakeUI, fakeIPSecManager)
		cliCommand = cli.Command{
			Name:        ipsec.IpsecCancelMetaData().Name,
			Description: ipsec.IpsecCancelMetaData().Description,
			Usage:       ipsec.IpsecCancelMetaData().Usage,
			Flags:       ipsec.IpsecCancelMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Context("cancel without contextID", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
		})
	})
	Context("cancel with wrong context id", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "abc")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Invalid input for 'Context ID'. It must be a positive integer.")).To(BeTrue())
		})
	})
	Context("cancel without confirmation", func() {
		It("return aborted", func() {
			fakeUI.Inputs("No")
			err := testhelpers.RunCommand(cliCommand, "1234")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This will cancel the IPSec: 1234 and cannot be undone. Continue?"}))
		})
	})
	Context("cancel with server fails", func() {
		BeforeEach(func() {
			fakeIPSecManager.CancelTunnelContextReturns(errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-f")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to cancel IPSec 1234.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("cancel with correct context id", func() {
		BeforeEach(func() {
			fakeIPSecManager.ApplyConfigurationReturns(nil)
		})
		It("return no error", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-f")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"IPSec 1234 is cancelled."}))
		})
	})
})

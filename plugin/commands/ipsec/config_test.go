package ipsec_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/ipsec"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("IPSec config", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeIPSecManager *testhelpers.FakeIPSECManager
		cmd              *ipsec.ConfigCommand
		cliCommand       cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeIPSecManager = new(testhelpers.FakeIPSECManager)
		cmd = ipsec.NewConfigCommand(fakeUI, fakeIPSecManager)
		cliCommand = cli.Command{
			Name:        metadata.IpsecConfigMetaData().Name,
			Description: metadata.IpsecConfigMetaData().Description,
			Usage:       metadata.IpsecConfigMetaData().Usage,
			Flags:       metadata.IpsecConfigMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Context("config without contextID", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
		})
	})
	Context("config with wrong context id", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "abc")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Invalid input for 'Context ID'. It must be a positive integer.")).To(BeTrue())
		})
	})
	Context("config with server fails", func() {
		BeforeEach(func() {
			fakeIPSecManager.ApplyConfigurationReturns(errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "1234")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to enqueue configuration request for IPSec 1234.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("config with correct context id", func() {
		BeforeEach(func() {
			fakeIPSecManager.ApplyConfigurationReturns(nil)
		})
		It("return no error", func() {
			err := testhelpers.RunCommand(cliCommand, "1234")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Configuration request received for IPSec 1234."}))
		})
	})
})

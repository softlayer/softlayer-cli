package ipsec_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ipsec"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("IPSec update translation", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeIPSecManager *testhelpers.FakeIPSECManager
		cmd              *ipsec.UpdateTranslationCommand
		cliCommand       cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeIPSecManager = new(testhelpers.FakeIPSECManager)
		cmd = ipsec.NewUpdateTranslationCommand(fakeUI, fakeIPSecManager)
		cliCommand = cli.Command{
			Name:        metadata.IpsecTransUpdataMetaData().Name,
			Description: metadata.IpsecTransUpdataMetaData().Description,
			Usage:       metadata.IpsecTransUpdataMetaData().Usage,
			Flags:       metadata.IpsecTransUpdataMetaData().Flags,
			Action:      cmd.Run,
		}
	})
	Context("update translation without context id", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires two arguments.")).To(BeTrue())
		})
	})
	Context("update translation without translation id", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "123")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires two arguments.")).To(BeTrue())
		})
	})
	Context("update translation with wrong context id", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "abc", "456")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Invalid input for 'Context ID'. It must be a positive integer.")).To(BeTrue())
		})
	})
	Context("update translation with wrong translation id", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "123", "abc")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Invalid input for 'Translation ID'. It must be a positive integer.")).To(BeTrue())
		})
	})
	Context("update translation with fail to update translation", func() {
		BeforeEach(func() {
			fakeIPSecManager.GetTranslationReturns(datatypes.Network_Tunnel_Module_Context_Address_Translation{Id: sl.Int(456)}, nil)
			fakeIPSecManager.UpdateTranslationReturns(datatypes.Network_Tunnel_Module_Context_Address_Translation{}, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "123", "456")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to update translation with ID 456 in IPSec 123.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("update translation ", func() {
		BeforeEach(func() {
			fakeIPSecManager.GetTranslationReturns(datatypes.Network_Tunnel_Module_Context_Address_Translation{Id: sl.Int(456)}, nil)
			fakeIPSecManager.UpdateTranslationReturns(datatypes.Network_Tunnel_Module_Context_Address_Translation{Id: sl.Int(456)}, nil)
		})
		It("succeed", func() {
			err := testhelpers.RunCommand(cliCommand, "123", "456", "-s", "1.2.3.4")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated translation with ID 456 in IPSec 123."}))
		})
		It("succeed", func() {
			err := testhelpers.RunCommand(cliCommand, "123", "456", "-r", "1.2.3.4")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated translation with ID 456 in IPSec 123."}))
		})
		It("succeed", func() {
			err := testhelpers.RunCommand(cliCommand, "123", "456", "-n", "test")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated translation with ID 456 in IPSec 123."}))
		})
		It("succeed", func() {
			err := testhelpers.RunCommand(cliCommand, "123", "456", "-s", "1.2.3.4", "-n", "test")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated translation with ID 456 in IPSec 123."}))
		})
		It("succeed", func() {
			err := testhelpers.RunCommand(cliCommand, "123", "456", "-r", "1.2.3.4", "-n", "test")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated translation with ID 456 in IPSec 123."}))
		})
		It("succeed", func() {
			err := testhelpers.RunCommand(cliCommand, "123", "456", "-s", "1.2.3.4", "-r", "5.6.7.8", "-n", "test")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated translation with ID 456 in IPSec 123."}))
		})
	})
})

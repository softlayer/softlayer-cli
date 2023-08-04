package ipsec_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ipsec"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("IPSec add translation", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeIPSecManager *testhelpers.FakeIPSECManager
		cliCommand       *ipsec.AddTranslationCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeIPSecManager = new(testhelpers.FakeIPSECManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = ipsec.NewAddTranslationCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.IPSECManager = fakeIPSecManager
	})
	Context("add translation without context id", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage : This command requires one argument")).To(BeTrue())
		})
	})
	Context("add translation with wrong context id", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Invalid input for 'Context ID'. It must be a positive integer.")).To(BeTrue())
		})
	})
	Context("add translation without static IP", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: '-s|--static-ip' is required")).To(BeTrue())
		})
	})
	Context("add translation without remote IP", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "1.2.3.4")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: '-r|--remote-ip' is required")).To(BeTrue())
		})
	})
	Context("add translation with get context fails", func() {
		BeforeEach(func() {
			fakeIPSecManager.GetTunnelContextReturns(datatypes.Network_Tunnel_Module_Context{}, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "1.2.3.4", "-r", "5.6.7.8")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to get IPSec with ID 123.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("add translation with create translation fails", func() {
		BeforeEach(func() {
			fakeIPSecManager.GetTunnelContextReturns(datatypes.Network_Tunnel_Module_Context{Id: sl.Int(123)}, nil)
			fakeIPSecManager.CreateTranslationReturns(datatypes.Network_Tunnel_Module_Context_Address_Translation{}, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "1.2.3.4", "-r", "5.6.7.8")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to create translation for IPSec with ID 123.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("add translation ", func() {
		BeforeEach(func() {
			fakeIPSecManager.GetTunnelContextReturns(datatypes.Network_Tunnel_Module_Context{Id: sl.Int(123)}, nil)
			fakeIPSecManager.CreateTranslationReturns(datatypes.Network_Tunnel_Module_Context_Address_Translation{Id: sl.Int(567)}, nil)
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "1.2.3.4", "-r", "5.6.7.8")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Created translation from 1.2.3.4 to 5.6.7.8 #567."}))
		})
	})
})

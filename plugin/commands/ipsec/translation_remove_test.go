package ipsec_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ipsec"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("IPSec remove translation", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeIPSecManager *testhelpers.FakeIPSECManager
		cliCommand       *ipsec.RemoveTranslationCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeIPSecManager = new(testhelpers.FakeIPSECManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = ipsec.NewRemoveTranslationCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.IPSECManager = fakeIPSecManager
	})
	Context("remove translation without context id", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires two arguments.")).To(BeTrue())
		})
	})
	Context("remove translation without translation id", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires two arguments.")).To(BeTrue())
		})
	})
	Context("remove translation with wrong context id", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "abc", "456")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Invalid input for 'Context ID'. It must be a positive integer.")).To(BeTrue())
		})
	})
	Context("remove translation with wrong translation id", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "abc")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Invalid input for 'Translation ID'. It must be a positive integer.")).To(BeTrue())
		})
	})
	Context("remove translation with fail to get translation", func() {
		BeforeEach(func() {
			fakeIPSecManager.GetTranslationReturns(datatypes.Network_Tunnel_Module_Context_Address_Translation{}, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "456")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to get translation with ID 456 from IPSec 123.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("remove translation with fail to remove translation", func() {
		BeforeEach(func() {
			fakeIPSecManager.GetTranslationReturns(datatypes.Network_Tunnel_Module_Context_Address_Translation{Id: sl.Int(456)}, nil)
			fakeIPSecManager.RemoveTranslationReturns(errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "456")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to remove translation with ID 456 from IPSec 123.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("remove translation ", func() {
		BeforeEach(func() {
			fakeIPSecManager.GetTranslationReturns(datatypes.Network_Tunnel_Module_Context_Address_Translation{Id: sl.Int(456)}, nil)
			fakeIPSecManager.RemoveTranslationReturns(nil)
		})
		It("succeed", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "456")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Removed translation with ID 456 from IPSec 123."}))
		})
	})
})

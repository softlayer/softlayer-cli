package ipsec_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ipsec"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("IPSec remove subnet", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeIPSecManager *testhelpers.FakeIPSECManager
		cliCommand       *ipsec.RemoveSubnetCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeIPSecManager = new(testhelpers.FakeIPSECManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = ipsec.NewRemoveSubnetCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.IPSECManager = fakeIPSecManager
	})
	Context("remove subnet without context id", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires three arguments.")).To(BeTrue())
		})
	})
	Context("remove subnet with wrong context id", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "abc", "bcd", "efg")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Invalid input for 'Context ID'. It must be a positive integer.")).To(BeTrue())
		})
	})
	Context("remove subnet with wrong subnet id", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "bcd", "efg")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Invalid input for 'Subnet ID'. It must be a positive integer.")).To(BeTrue())
		})
	})
	Context("remove subnet with wrong subnet type", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "456", "efg")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: The subnet type has to be either internal, or remote or service.")).To(BeTrue())
		})
	})
	Context("remove internal subnet with server fail", func() {
		BeforeEach(func() {
			fakeIPSecManager.RemoveInternalSubnetReturns(errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "456", "internal")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to remove internal subnet #456 from IPSec 123.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("remove remote subnet with server fail", func() {
		BeforeEach(func() {
			fakeIPSecManager.RemoveRemoteSubnetReturns(errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "456", "remote")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to remove remote subnet #456 from IPSec 123.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("remove service subnet with server fail", func() {
		BeforeEach(func() {
			fakeIPSecManager.RemoveServiceSubnetReturns(errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "456", "service")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to remove service subnet #456 from IPSec 123.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("remove internal subnet", func() {
		BeforeEach(func() {
			fakeIPSecManager.RemoveInternalSubnetReturns(nil)
		})
		It("succeed", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "456", "internal")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Removed internal subnet #456 from IPSec 123."}))
		})
	})
	Context("remove remote subnet", func() {
		BeforeEach(func() {
			fakeIPSecManager.RemoveRemoteSubnetReturns(nil)
		})
		It("succeed", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "456", "remote")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Removed remote subnet #456 from IPSec 123."}))
		})
	})
	Context("remove service subnet", func() {
		BeforeEach(func() {
			fakeIPSecManager.RemoveServiceSubnetReturns(nil)
		})
		It("succeed", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "456", "service")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Removed service subnet #456 from IPSec 123."}))
		})
	})
})

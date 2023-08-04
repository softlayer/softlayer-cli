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
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ipsec"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("IPSec add subnet", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeIPSecManager *testhelpers.FakeIPSECManager
		cliCommand       *ipsec.AddSubnetCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeIPSecManager = new(testhelpers.FakeIPSECManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = ipsec.NewAddSubnetCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.IPSECManager = fakeIPSecManager
	})
	Context("Add subnet without context id", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage : This command requires one argument")).To(BeTrue())
		})
	})
	Context("Add subnet with wrong context id", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Invalid input for 'Context ID'. It must be a positive integer.")).To(BeTrue())
		})
	})
	Context("Add subnet with wrong subnet type", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "456", "-t", "abc")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: The subnet type has to be either internal, or remote or service.")).To(BeTrue())
		})
	})
	Context("Add subnet without subnetId or subnet identifier", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-t", "remote")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: Either -s|--subnet-id or -n|--network must be provided.")).To(BeTrue())
		})
	})
	Context("Add subnet with subnet identifier but wrong subnet type", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-n", "1.1.2.3", "-t", "internal")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: Unable to create internal subnet.")).To(BeTrue())
		})
	})
	//TODO create remote subnet cases
	Context("Add internal subnet with get context fail", func() {
		BeforeEach(func() {
			fakeIPSecManager.GetTunnelContextReturns(datatypes.Network_Tunnel_Module_Context{}, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "456", "-t", "internal")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to get IPSec with ID 123.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("Add internal subnet with add subnet fail", func() {
		BeforeEach(func() {
			fakeIPSecManager.AddInternalSubnetReturns(errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "456", "-t", "internal")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to add internal subnet #456 to IPSec 123.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("Add remote subnet with add subnet fail", func() {
		BeforeEach(func() {
			fakeIPSecManager.AddRemoteSubnetReturns(errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "456", "-t", "remote")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to add remote subnet #456 to IPSec 123.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("Add service subnet with add subnet fail", func() {
		BeforeEach(func() {
			fakeIPSecManager.AddServiceSubnetReturns(errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "456", "-t", "service")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to add service subnet #456 to IPSec 123.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})

	Context("Add internal subnet", func() {
		BeforeEach(func() {
			fakeIPSecManager.AddInternalSubnetReturns(nil)
		})
		It("succeed", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "456", "-t", "internal")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Added internal subnet #456 to IPSec 123."}))
		})
	})
	Context("Add remote subnet", func() {
		BeforeEach(func() {
			fakeIPSecManager.AddRemoteSubnetReturns(nil)
		})
		It("Add", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "456", "-t", "remote")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Added remote subnet #456 to IPSec 123."}))
		})
	})
	Context("Add service subnet", func() {
		BeforeEach(func() {
			fakeIPSecManager.AddServiceSubnetReturns(nil)
		})
		It("succeed", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "-s", "456", "-t", "service")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Added service subnet #456 to IPSec 123."}))
		})
	})
})

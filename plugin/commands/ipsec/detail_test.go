package ipsec_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ipsec"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("IPSec detail", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeIPSecManager *testhelpers.FakeIPSECManager
		cliCommand       *ipsec.DetailCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeIPSecManager = new(testhelpers.FakeIPSECManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = ipsec.NewDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.IPSECManager = fakeIPSecManager
	})
	Context("detail without contextID", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument")).To(BeTrue())
		})
	})
	Context("detail with wrong context id", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Invalid input for 'Context ID'. It must be a positive integer.")).To(BeTrue())
		})
	})
	Context("detail with server fails", func() {
		BeforeEach(func() {
			fakeIPSecManager.GetTunnelContextReturns(datatypes.Network_Tunnel_Module_Context{}, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to get IPSec with ID 1234.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
})

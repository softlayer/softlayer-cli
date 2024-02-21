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

var _ = Describe("IPSec order", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeIPSecManager *testhelpers.FakeIPSECManager
		cliCommand       *ipsec.OrderCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeIPSecManager = new(testhelpers.FakeIPSECManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = ipsec.NewOrderCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.IPSECManager = fakeIPSecManager
	})
	Context("order without -d", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: '-d|--datacenter' is required")).To(BeTrue())
		})
	})
	Context("order with server fails", func() {
		BeforeEach(func() {
			fakeIPSecManager.OrderTunnelContextReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-d", "dal09")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to order IPSec.Please try again later.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("order", func() {
		BeforeEach(func() {
			fakeIPSecManager.OrderTunnelContextReturns(datatypes.Container_Product_Order_Receipt{
				OrderId: sl.Int(12345),
			}, nil)
		})
		It("succeed", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-d", "dal09")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Order 12345 was placed."}))
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"You may run 'ibmcloud sl ipsec list --order 12345' to find this IPSec VPN after it is ready."}))
		})
	})
})

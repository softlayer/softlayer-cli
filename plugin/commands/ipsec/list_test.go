package ipsec_test

import (
	"errors"
	"strings"
	"time"

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

var _ = Describe("IPSec list", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeIPSecManager *testhelpers.FakeIPSECManager
		cliCommand       *ipsec.ListCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeIPSecManager = new(testhelpers.FakeIPSECManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = ipsec.NewListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.IPSECManager = fakeIPSecManager
	})
	Context("list with server fails", func() {
		BeforeEach(func() {
			fakeIPSecManager.GetTunnelContextsReturns(nil, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to get IPSec on your account.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("list", func() {
		BeforeEach(func() {
			fakeIPSecManager.GetTunnelContextsReturns(nil, nil)
		})
		It("return no ipsec", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"No IPSec was found."}))
		})
	})
	Context("list", func() {
		BeforeEach(func() {
			created := time.Now()
			fakeIPSecManager.GetTunnelContextsReturns([]datatypes.Network_Tunnel_Module_Context{
				datatypes.Network_Tunnel_Module_Context{
					Id:                    sl.Int(123),
					Name:                  sl.String("abc"),
					FriendlyName:          sl.String("ABC"),
					InternalPeerIpAddress: sl.String("1.1.1.2"),
					CustomerPeerIpAddress: sl.String("2.2.2.3"),
					CreateDate:            sl.Time(created),
				},
			}, nil)
		})
		It("return ipseclist", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"123"}))
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"abc"}))
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"ABC"}))
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1.1.1.2"}))
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2.2.2.3"}))
		})
	})
})

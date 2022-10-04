package vlan_test

import (
	"errors"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/vlan"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VLAN options", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *vlan.OptionsCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeNetworkManager *testhelpers.FakeNetworkManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = vlan.NewOptionsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cliCommand.NetworkManager = fakeNetworkManager
	})

	Describe("VLAN options", func() {
		Context("VLAN options but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListDatacentersReturns(nil, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list datacenters."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("VLAN options", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListDatacentersReturns(map[int]string{1: "dal07"}, nil)
				fakeNetworkManager.ListRoutersReturns([]string{"bcr01a.dal07", "fcr01a.dal07"}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"public,private"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"bcr01a.dal07,fcr01a.dal07"}))
			})
		})
	})
})

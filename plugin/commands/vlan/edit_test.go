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

var _ = Describe("VLAN Edit", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *vlan.EditCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeNetworkManager *testhelpers.FakeNetworkManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = vlan.NewEditCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cliCommand.NetworkManager = fakeNetworkManager
	})

	Describe("VLAN edit", func() {
		Context("VLAN edit without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("VLAN edit with wrong vlan id", func() {
			It("error resolving vlan ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'VLAN ID'. It must be a positive integer."))
			})
		})

		Context("VLAN edit without -n", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-n|--name' is required"))
			})
		})

		Context("VLAN edit with correct vlan id but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.EditVlanReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-n", "myvlan")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to edit VLAN: 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("VLAN edit with correct vlan id", func() {
			BeforeEach(func() {
				fakeNetworkManager.EditVlanReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-n", "myvlan")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"VLAN 1234 was updated."}))
			})
		})
	})
})

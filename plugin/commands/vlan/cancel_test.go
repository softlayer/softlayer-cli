package vlan_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/vlan"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VLAN Cancel", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *vlan.CancelCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeNetworkManager *testhelpers.FakeNetworkManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = vlan.NewCancelCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cliCommand.NetworkManager = fakeNetworkManager
	})

	Describe("VLAN cancel", func() {
		Context("VLAN cancel without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("VLAN cancel with wrong vlan id", func() {
			It("error resolving vlan ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'VLAN ID'. It must be a positive integer."))
			})
		})

		Context("VLAN cancel with correct vlan id", func() {
			BeforeEach(func() {
				fakeNetworkManager.CancelVLANReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("VLAN 1234 was cancelled."))
			})
		})

		Context("VLAN cancel with correct vlan id but vlan is not found", func() {
			BeforeEach(func() {
				fakeNetworkManager.CancelVLANReturns(errors.New("SoftLayer_Exception_ObjectNotFound"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to find VLAN with ID 1234."))
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ObjectNotFound"))
			})
		})

		Context("VLAN cancel with correct vlan id but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.CancelVLANReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to cancel VLAN 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Unable to cancel due to reasons.", func() {
			BeforeEach(func() {
				fakeError := []string{"BAD"}
				fakeNetworkManager.GetCancelFailureReasonsReturns(fakeError)
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to cancel VLAN 1234."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("BAD"))
			})
		})
	})
})

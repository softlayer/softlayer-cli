package vlan_test

import (
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/vlan"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VLAN Detail", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *vlan.DetailCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeHandler        *testhelpers.FakeTransportHandler
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession(nil)
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = vlan.NewDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})
	AfterEach(func() {
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})

	Describe("VLAN detail", func() {
		Context("VLAN detail without ID", func() {
			It("Error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstrings())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("VLAN detail with wrong vlan id", func() {
			It("Error resolving vlan ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'VLAN ID'. It must be a positive integer."))
			})
		})

		Context("VLAN detail with correct vlan id but server API call fails", func() {
			BeforeEach(func() {
				fakeHandler.AddApiError("SoftLayer_Network_Vlan", "getObject", 500, "Internal Server Error")
			})
			It("Error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get VLAN: 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("VLAN Happy Path", func() {
			It("Success", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1262125"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("169.55.16.48"))
			})
		})
		Context("VLAN Happy Path", func() {
			It("Success", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--output=JSON")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`domain": "stage1.ng.bluemix.net"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"subnetType": "ADDITIONAL_PRIMARY"`))
			})
		})
		Context("VLAN Happy Path: trunk details", func() {
			It("Success", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("testibm"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SECONDARY_ON_VLAN"))
			})
		})
	})
})

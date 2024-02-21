package hardware_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VLAN-ADD Tests", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *hardware.VlanAddCommand
		fakeSession *session.Session
		fakeHandler *testhelpers.FakeTransportHandler
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession(nil)
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewVlanAddCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})
	AfterEach(func() {
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})

	Describe("VLAN-ADD", func() {
		Context("Argument Tests", func() {
			It("Missing HardwareID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: requires at least 2 arg(s), only received 0"))
			})
			It("Bad HardwareID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "asdf", "12345")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'."))
			})
			It("Bad VlanId", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345", "zzzz")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'VLAN ID'."))
			})
		})
		Context("API Errors", func() {
			It("Bad Hardware", func() {
				fakeHandler.AddApiError("SoftLayer_Hardware_Server", "getObject", 500, "Internal Server Error")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1000", "9999")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get hardware server: 1000."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
			It("Bad Vlan", func() {
				fakeHandler.AddApiError("SoftLayer_Network_Vlan", "getObject", 500, "Internal Server Error")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1000", "9999")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get VLAN: 9999."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
			It("Trunk Fail", func() {
				fakeHandler.AddApiError("SoftLayer_Network_Component", "addNetworkVlanTrunks", 500, "Trunk Error")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1000", "9999")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Trunk Error: Trunk Error (HTTP 500)"))
			})
			It("Trunk Fail Private", func() {
				fakeHandler.AddApiError("SoftLayer_Network_Component", "addNetworkVlanTrunks", 500, "Trunk Error")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1000", "9990")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Trunk Error: Trunk Error (HTTP 500)"))
			})
			AfterEach(func() {
				fakeHandler.ClearErrors()
			})
		})
		Context("Happy Path", func() {
			It("1 Pub 1 Pri", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1000", "9999", "9990")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("9990    445566   fmirPublic"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("10011   1163     fmirPublic"))
				apiCalls := fakeHandler.ApiCallLogs

				Expect(len(apiCalls)).To(Equal(5))
				// Trying out https://pkg.go.dev/github.com/onsi/gomega/gstruct for matching API calls
				Expect(apiCalls[0]).To(MatchFields(IgnoreExtras, Fields{
					"Service": Equal("SoftLayer_Hardware_Server"),
					"Method":  Equal("getObject"),
					"Options": PointTo(MatchFields(IgnoreExtras, Fields{"Id": PointTo(Equal(1000))})),
				}))
				Expect(apiCalls[1]).To(MatchFields(IgnoreExtras, Fields{
					"Service": Equal("SoftLayer_Network_Vlan"),
					"Method":  Equal("getObject"),
					"Options": PointTo(MatchFields(IgnoreExtras, Fields{"Id": PointTo(Equal(9999))})),
				}))
				Expect(apiCalls[2]).To(MatchFields(IgnoreExtras, Fields{
					"Service": Equal("SoftLayer_Network_Vlan"),
					"Method":  Equal("getObject"),
					"Options": PointTo(MatchFields(IgnoreExtras, Fields{"Id": PointTo(Equal(9990))})),
				}))
				Expect(apiCalls[3]).To(MatchFields(IgnoreExtras, Fields{
					"Service": Equal("SoftLayer_Network_Component"),
					"Method":  Equal("addNetworkVlanTrunks"),
					"Options": PointTo(MatchFields(IgnoreExtras, Fields{"Id": PointTo(Equal(10011))})),
				}))
				Expect(apiCalls[4]).To(MatchFields(IgnoreExtras, Fields{
					"Service": Equal("SoftLayer_Network_Component"),
					"Method":  Equal("addNetworkVlanTrunks"),
					"Options": PointTo(MatchFields(IgnoreExtras, Fields{"Id": PointTo(Equal(90099))})),
				}))
			})
		})
	})
})

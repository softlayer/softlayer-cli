package hardware_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VLAN-REMOVE Tests", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *hardware.VlanTrunkableCommand
		fakeSession *session.Session
		fakeHandler *testhelpers.FakeTransportHandler
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession(nil)
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewVlanTrunkableCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})
	AfterEach(func() {
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})

	Describe("VLAN-TRUNKABLE", func() {
		Context("Argument Tests", func() {
			It("Missing HardwareID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
			It("Bad HardwareID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "asdf")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'."))
			})
		})
		Context("API Errors", func() {
			It("Bad Hardware", func() {
				fakeHandler.AddApiError("SoftLayer_Hardware_Server", "getObject", 500, "Internal Server Error")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1001")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get hardware server: 1001."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})

			AfterEach(func() {
				fakeHandler.ClearErrors()
			})
		})
		Context("Happy Path", func() {
			It("Happy Path", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1001")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1404269   dal10.bcr01.1632       ibm iSCSI DAL10   PRIVATE"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2282899   dal10.fcr01.1163       fmirPublic        PUBLIC"))
				apiCalls := fakeHandler.ApiCallLogs

				Expect(len(apiCalls)).To(Equal(1))
				// Trying out https://pkg.go.dev/github.com/onsi/gomega/gstruct for matching API calls
				Expect(apiCalls[0]).To(MatchFields(IgnoreExtras, Fields{
					"Service": Equal("SoftLayer_Hardware_Server"),
					"Method":  Equal("getObject"),
					"Options": PointTo(MatchFields(IgnoreExtras, Fields{"Id": PointTo(Equal(1001))})),
				}))
			})
		})
	})
})

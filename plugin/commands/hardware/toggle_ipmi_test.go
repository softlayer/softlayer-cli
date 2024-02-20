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

var _ = Describe("TOGGLE-IPMI Tests", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *hardware.ToggleIPMICommand
		fakeSession *session.Session
		fakeHandler *testhelpers.FakeTransportHandler
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession(nil)
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewToggleIPMICommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})
	AfterEach(func() {
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})

	Describe("TOGGLE-IPMI", func() {
		Context("Argument Tests", func() {
			It("Missing HardwareID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
			It("Bad HardwareID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "asdf", "--enable")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'."))
			})
			It("Both enable and disble", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345", "--enable", "--disable")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--enable', '--disable' are exclusive."))
			})
			It("Niether enable nor disble", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Either '--enable' or '--disable' is required."))
			})
		})
		Context("API Errors", func() {
			It("Bad Hardware", func() {
				fakeHandler.AddApiError("SoftLayer_Hardware_Server", "toggleManagementInterface", 500, "Internal Server Error")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1000", "--enable")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to toggle IPMI interface of hardware server '1000'."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
			AfterEach(func() {
				fakeHandler.ClearErrors()
			})
		})
		Context("Happy Path", func() {
			It("Success Toggle", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1000", "--enable")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Successfully send request to toggle IPMI interface of hardware server '1000'."))
				apiCalls := fakeHandler.ApiCallLogs

				Expect(len(apiCalls)).To(Equal(1))
				// Trying out https://pkg.go.dev/github.com/onsi/gomega/gstruct for matching API calls
				Expect(apiCalls[0]).To(MatchFields(IgnoreExtras, Fields{
					"Service": Equal("SoftLayer_Hardware_Server"),
					"Method":  Equal("toggleManagementInterface"),
					"Options": PointTo(MatchFields(IgnoreExtras, Fields{"Id": PointTo(Equal(1000))})),
					"Args":    MatchAllElementsWithIndex(IndexIdentity, Elements{"0": PointTo(Equal(true))}),
				}))
			})
			It("Success UnToggle", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1000", "--disable")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Successfully send request to toggle IPMI interface of hardware server '1000'."))
				apiCalls := fakeHandler.ApiCallLogs

				Expect(len(apiCalls)).To(Equal(1))
				// Trying out https://pkg.go.dev/github.com/onsi/gomega/gstruct for matching API calls
				Expect(apiCalls[0]).To(MatchFields(IgnoreExtras, Fields{
					"Service": Equal("SoftLayer_Hardware_Server"),
					"Method":  Equal("toggleManagementInterface"),
					"Options": PointTo(MatchFields(IgnoreExtras, Fields{"Id": PointTo(Equal(1000))})),
					"Args":    MatchAllElementsWithIndex(IndexIdentity, Elements{"0": PointTo(Equal(false))}),
				}))
			})
		})
	})
})
